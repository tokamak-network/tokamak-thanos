package batcher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	_ "net/http/pprof"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"

	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/ethereum-optimism/optimism/op-batcher/metrics"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum-optimism/optimism/op-service/dial"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
)

var (
	ErrBatcherNotRunning = errors.New("batcher is not running")
	emptyTxData          = txData{
		frames: []frameData{
			{
				data: []byte{},
			},
		},
	}
	SetMaxDASizeMethod = "miner_setMaxDASize"
)

type txRef struct {
	id       txID
	isCancel bool
	isBlob   bool
}

func (r txRef) String() string {
	return r.string(func(id txID) string { return id.String() })
}

func (r txRef) TerminalString() string {
	return r.string(func(id txID) string { return id.TerminalString() })
}

func (r txRef) string(txIDStringer func(txID) string) string {
	if r.isCancel {
		if r.isBlob {
			return "blob-cancellation"
		} else {
			return "calldata-cancellation"
		}
	}
	return txIDStringer(r.id)
}

type L1Client interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
}

type L2Client interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

type RollupClient interface {
	SyncStatus(ctx context.Context) (*eth.SyncStatus, error)
}

// DriverSetup is the collection of input/output interfaces and configuration that the driver operates on.
type DriverSetup struct {
	Log               log.Logger
	Metr              metrics.Metricer
	RollupConfig      *rollup.Config
	Config            BatcherConfig
	Txmgr             txmgr.TxManager
	L1Client          L1Client
	EndpointProvider  dial.L2EndpointProvider
	ChannelConfig     ChannelConfigProvider
	AltDA             *altda.DAClient
	ChannelOutFactory ChannelOutFactory
	ActiveSeqChanged  chan struct{} // optional
}

// BatchSubmitter encapsulates a service responsible for submitting L2 tx
// batches to L1 for availability.
type BatchSubmitter struct {
	DriverSetup

	wg                               *sync.WaitGroup
	shutdownCtx, killCtx             context.Context
	cancelShutdownCtx, cancelKillCtx context.CancelFunc

	mutex   sync.Mutex
	running bool

	txpoolMutex       sync.Mutex // guards txpoolState and txpoolBlockedBlob
	txpoolState       TxPoolState
	txpoolBlockedBlob bool

	channelMgrMutex sync.Mutex // guards channelMgr and prevCurrentL1
	channelMgr      *channelManager
	prevCurrentL1   eth.L1BlockRef // cached CurrentL1 from the last syncStatus
}

// NewBatchSubmitter initializes the BatchSubmitter driver from a preconfigured DriverSetup
func NewBatchSubmitter(setup DriverSetup) *BatchSubmitter {
	state := NewChannelManager(setup.Log, setup.Metr, setup.ChannelConfig, setup.RollupConfig)
	if setup.ChannelOutFactory != nil {
		state.SetChannelOutFactory(setup.ChannelOutFactory)
	}

	return &BatchSubmitter{
		DriverSetup: setup,
		channelMgr:  state,
	}
}

func (l *BatchSubmitter) StartBatchSubmitting() error {
	l.Log.Info("Starting Batch Submitter")

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.running {
		return errors.New("batcher is already running")
	}
	l.running = true

	l.shutdownCtx, l.cancelShutdownCtx = context.WithCancel(context.Background())
	l.killCtx, l.cancelKillCtx = context.WithCancel(context.Background())
	l.clearState(l.shutdownCtx)
	l.wg = &sync.WaitGroup{}

	if err := l.waitForL2Genesis(); err != nil {
		return fmt.Errorf("error waiting for L2 genesis: %w", err)
	}

	if l.Config.WaitNodeSync {
		err := l.waitNodeSync()
		if err != nil {
			return fmt.Errorf("error waiting for node sync: %w", err)
		}
	}

	receiptsCh := make(chan txmgr.TxReceipt[txRef])

	l.txpoolState = TxpoolGood // no need to lock mutex as no other routines yet exist

	// Channels used to signal between the loops
	pendingBytesUpdated := make(chan int64, 1)
	publishSignal := make(chan struct{})

	// DA throttling loop should always be started except for testing (indicated by ThrottleThreshold == 0)
	if l.Config.ThrottleThreshold > 0 {
		l.wg.Add(1)
		go l.throttlingLoop(l.wg, pendingBytesUpdated) // ranges over pendingBytesUpdated channel
	} else {
		l.Log.Warn("Throttling loop is DISABLED due to 0 throttle-interval. This should not be disabled in prod.")
	}

	l.wg.Add(3)
	go l.receiptsLoop(l.wg, receiptsCh)                                            // ranges over receiptsCh channel
	go l.publishingLoop(l.killCtx, l.wg, receiptsCh, publishSignal)                // ranges over publishSignal, spawns routines which send on receiptsCh. Closes receiptsCh when done.
	go l.blockLoadingLoop(l.shutdownCtx, l.wg, pendingBytesUpdated, publishSignal) // sends on pendingBytesUpdated (if throttling enabled), and publishSignal. Closes them both when done

	l.Log.Info("Batch Submitter started")
	return nil
}

// waitForL2Genesis waits for the L2 genesis time to be reached.
func (l *BatchSubmitter) waitForL2Genesis() error {
	genesisTime := time.Unix(int64(l.RollupConfig.Genesis.L2Time), 0)
	now := time.Now()
	if now.After(genesisTime) {
		return nil
	}

	l.Log.Info("Waiting for L2 genesis", "genesisTime", genesisTime)

	// Create a ticker that fires every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	genesisTrigger := time.After(time.Until(genesisTime))

	for {
		select {
		case <-ticker.C:
			remaining := time.Until(genesisTime)
			l.Log.Info("Waiting for L2 genesis", "remainingTime", remaining.Round(time.Second))
		case <-genesisTrigger:
			l.Log.Info("L2 genesis time reached")
			return nil
		case <-l.shutdownCtx.Done():
			return errors.New("batcher stopped")
		}
	}
}

func (l *BatchSubmitter) StopBatchSubmittingIfRunning(ctx context.Context) error {
	err := l.StopBatchSubmitting(ctx)
	if errors.Is(err, ErrBatcherNotRunning) {
		return nil
	}
	return err
}

// StopBatchSubmitting stops the batch-submitter loop, and force-kills if the provided ctx is done.
func (l *BatchSubmitter) StopBatchSubmitting(ctx context.Context) error {
	l.Log.Info("Stopping Batch Submitter")

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if !l.running {
		return ErrBatcherNotRunning
	}
	l.running = false

	// go routine will call cancelKill() if the passed in ctx is ever Done
	cancelKill := l.cancelKillCtx
	wrapped, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-wrapped.Done()
		cancelKill()
	}()

	l.cancelShutdownCtx()
	l.wg.Wait()
	l.cancelKillCtx()

	l.Log.Info("Batch Submitter stopped")
	return nil
}

// loadBlocksIntoState loads the blocks between start and end (inclusive).
// If there is a reorg, it will return an error.
func (l *BatchSubmitter) loadBlocksIntoState(ctx context.Context, start, end uint64) error {
	if end < start {
		return fmt.Errorf("start number is > end number %d,%d", start, end)
	}

	// we don't want to print it in the 1-block case as `loadBlockIntoState` already does
	if end > start {
		l.Log.Info("Loading range of multiple blocks into state", "start", start, "end", end)
	}

	var latestBlock *types.Block
	// Add all blocks to "state"
	for i := start; i <= end; i++ {
		block, err := l.loadBlockIntoState(ctx, i)
		if errors.Is(err, ErrReorg) {
			l.Log.Warn("Found L2 reorg", "block_number", i)
			return err
		} else if err != nil {
			l.Log.Warn("Failed to load block into state", "err", err)
			return err
		}
		latestBlock = block
	}

	l2ref, err := derive.L2BlockToBlockRef(l.RollupConfig, latestBlock)
	if err != nil {
		l.Log.Warn("Invalid L2 block loaded into state", "err", err)
		return err
	}

	l.Metr.RecordL2BlocksLoaded(l2ref)
	return nil
}

// loadBlockIntoState fetches & stores a single block into `state`. It returns the block it loaded.
func (l *BatchSubmitter) loadBlockIntoState(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	l2Client, err := l.EndpointProvider.EthClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting L2 client: %w", err)
	}

	cCtx, cancel := context.WithTimeout(ctx, l.Config.NetworkTimeout)
	defer cancel()

	block, err := l2Client.BlockByNumber(cCtx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return nil, fmt.Errorf("getting L2 block: %w", err)
	}

	l.channelMgrMutex.Lock()
	defer l.channelMgrMutex.Unlock()
	if err := l.channelMgr.AddL2Block(block); err != nil {
		return nil, fmt.Errorf("adding L2 block to state: %w", err)
	}

	l.Log.Info("Added L2 block to local state", "block", eth.ToBlockID(block), "tx_count", len(block.Transactions()), "time", block.Time())
	return block, nil
}

func (l *BatchSubmitter) getSyncStatus(ctx context.Context) (*eth.SyncStatus, error) {
	rollupClient, err := l.EndpointProvider.RollupClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting rollup client: %w", err)
	}

	var (
		syncStatus *eth.SyncStatus
		backoff    = time.Second
		maxBackoff = 30 * time.Second
	)
	timer := time.NewTimer(backoff)
	defer timer.Stop()

	for {
		cCtx, cancel := context.WithTimeout(ctx, l.Config.NetworkTimeout)
		syncStatus, err = rollupClient.SyncStatus(cCtx)
		cancel()

		// Ensure that we have the sync status
		if err != nil {
			return nil, fmt.Errorf("failed to get sync status: %w", err)
		}

		// If we have a head, break out of the loop
		if syncStatus.HeadL1 != (eth.L1BlockRef{}) {
			break
		}

		// Empty sync status, implement backoff
		l.Log.Info("Received empty sync status, backing off", "backoff", backoff)
		select {
		case <-timer.C:
			backoff *= 2
			backoff = min(backoff, maxBackoff)
			// Reset timer to tick of the new backoff time again
			timer.Reset(backoff)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return syncStatus, nil
}

// The following things occur:
// New L2 block (reorg or not)
// L1 transaction is confirmed
//
// What the batcher does:
// Ensure that channels are created & submitted as frames for an L2 range
//
// Error conditions:
// Submitted batch, but it is not valid
// Missed L2 block somehow.

type TxPoolState int

const (
	// Txpool states.  Possible state transitions:
	//   TxpoolGood -> TxpoolBlocked:
	//     happens when ErrAlreadyReserved is ever returned by the TxMgr.
	//   TxpoolBlocked -> TxpoolCancelPending:
	//     happens once the send loop detects the txpool is blocked, and results in attempting to
	//     send a cancellation transaction.
	//   TxpoolCancelPending -> TxpoolGood:
	//     happens once the cancel transaction completes, whether successfully or in error.
	TxpoolGood TxPoolState = iota
	TxpoolBlocked
	TxpoolCancelPending
)

// sendToThrottlingLoop sends the current pending bytes to the throttling loop.
// It is not blocking, no signal will be sent if the channel is full.
func (l *BatchSubmitter) sendToThrottlingLoop(pendingBytesUpdated chan int64) {
	if l.Config.ThrottleThreshold == 0 {
		return
	}

	l.channelMgrMutex.Lock()
	pendingBytes := l.channelMgr.PendingDABytes()
	l.channelMgrMutex.Unlock()

	// notify the throttling loop it may be time to initiate throttling without blocking
	select {
	case pendingBytesUpdated <- pendingBytes:
	default:
	}
}

// trySignal trys to send an emptry struct  on the provided channel.
// It is not blocking, no signal will be sent if the channel is full.
func trySignal(c chan struct{}) {
	select {
	case c <- struct{}{}:
	default:
	}
}

// setTxPoolState locks the mutex, sets the parameters to the supplied ones, and release the mutex.
func (l *BatchSubmitter) setTxPoolState(txPoolState TxPoolState, txPoolBlockedBlob bool) {
	l.txpoolMutex.Lock()
	l.txpoolState = txPoolState
	l.txpoolBlockedBlob = txPoolBlockedBlob
	l.txpoolMutex.Unlock()
}

// syncAndPrune computes actions to take based on the current sync status, prunes the channel manager state
// and returns blocks to load.
func (l *BatchSubmitter) syncAndPrune(syncStatus *eth.SyncStatus) *inclusiveBlockRange {
	l.channelMgrMutex.Lock()
	defer l.channelMgrMutex.Unlock()

	// Decide appropriate actions
	syncActions, outOfSync := computeSyncActions(*syncStatus, l.prevCurrentL1, l.channelMgr.blocks, l.channelMgr.channelQueue, l.Log, l.Config.PreferLocalSafeL2)

	if outOfSync {
		// If the sequencer is out of sync
		// do nothing and wait to see if it has
		// got in sync on the next tick.
		l.Log.Warn("Sequencer is out of sync, retrying next tick.")
		return syncActions.blocksToLoad
	}

	l.prevCurrentL1 = syncStatus.CurrentL1

	// Manage existing state / garbage collection
	if syncActions.clearState != nil {
		l.channelMgr.Clear(*syncActions.clearState)
	} else {
		l.channelMgr.PruneSafeBlocks(syncActions.blocksToPrune)
		l.channelMgr.PruneChannels(syncActions.channelsToPrune)
	}
	return syncActions.blocksToLoad
}

// publishingLoop:
// -  waits for a signal that blocks have been loaded
// -  drives the creation of channels and frames
// -  sends transactions to the DA layer
func (l *BatchSubmitter) publishingLoop(ctx context.Context, wg *sync.WaitGroup, receiptsCh chan txmgr.TxReceipt[txRef], publishSignal chan struct{}) {
	defer close(receiptsCh)
	defer wg.Done()

	daGroup := &errgroup.Group{}
	// errgroup with limit of 0 means no goroutine is able to run concurrently,
	// so we only set the limit if it is greater than 0.
	if l.Config.MaxConcurrentDARequests > 0 {
		daGroup.SetLimit(int(l.Config.MaxConcurrentDARequests))
	}
	txQueue := txmgr.NewQueue[txRef](ctx, l.Txmgr, l.Config.MaxPendingTransactions)

	for range publishSignal {
		if !l.checkTxpool(txQueue, receiptsCh) {
			continue
		}
		l.publishStateToL1(ctx, txQueue, receiptsCh, daGroup)
	}

	// We _must_ wait for all senders on receiptsCh to finish before we can close it.
	if err := txQueue.Wait(); err != nil {
		if !errors.Is(err, context.Canceled) {
			l.Log.Error("error waiting for transactions to complete", "err", err)
		}
	}
	l.Log.Info("publishingLoop returning")
}

// blockLoadingLoop
// -  polls the sequencer,
// -  prunes the channel manager state (i.e. safe blocks)
// -  loads unsafe blocks from the sequencer
func (l *BatchSubmitter) blockLoadingLoop(ctx context.Context, wg *sync.WaitGroup, pendingBytesUpdated chan int64, publishSignal chan struct{}) {
	ticker := time.NewTicker(l.Config.PollInterval)
	defer ticker.Stop()
	defer close(pendingBytesUpdated)
	defer close(publishSignal)
	defer wg.Done()
	for {
		select {
		case <-ticker.C:
			syncStatus, err := l.getSyncStatus(ctx)
			if err != nil {
				l.Log.Warn("could not get sync status", "err", err)
				continue
			}

			blocksToLoad := l.syncAndPrune(syncStatus)

			if blocksToLoad != nil {
				// Get fresh unsafe blocks
				if err := l.loadBlocksIntoState(ctx, blocksToLoad.start, blocksToLoad.end); errors.Is(err, ErrReorg) {
					l.Log.Warn("error loading blocks, clearing state and waiting for node sync", "err", err)
					l.waitNodeSyncAndClearState()
				} else {
					l.sendToThrottlingLoop(pendingBytesUpdated) // we have increased the pending data. Signal the throttling loop to check if it should throttle.
				}
			}
			trySignal(publishSignal) // always signal the write loop to ensure we periodically publish even if we aren't loading blocks
		case <-ctx.Done():
			l.Log.Info("blockLoadingLoop returning")
			return
		}
	}
}

// receiptsLoop handles transaction receipts from the DA layer
func (l *BatchSubmitter) receiptsLoop(wg *sync.WaitGroup, receiptsCh chan txmgr.TxReceipt[txRef]) {
	defer wg.Done()
	l.Log.Info("Starting receipts processing loop")
	for r := range receiptsCh {

		if errors.Is(r.Err, txpool.ErrAlreadyReserved) && l.txpoolState == TxpoolGood {
			l.setTxPoolState(TxpoolBlocked, r.ID.isBlob)
			l.Log.Warn("incompatible tx in txpool", "id", r.ID, "is_blob", r.ID.isBlob)
		} else if r.ID.isCancel && l.txpoolState == TxpoolCancelPending {
			// Set state to TxpoolGood even if the cancellation transaction ended in error
			// since the stuck transaction could have cleared while we were waiting.
			l.setTxPoolState(TxpoolGood, l.txpoolBlockedBlob)
			l.Log.Info("txpool may no longer be blocked", "err", r.Err)
		}
		l.Log.Info("Handling receipt", "id", r.ID)
		l.handleReceipt(r)
	}
	l.Log.Info("receiptsLoop returning")
}

// throttlingLoop monitors the backlog in bytes we need to make available, and appropriately enables or disables
// throttling of incoming data prevent the backlog from growing too large. By looping & calling the miner API setter
// continuously, we ensure the engine currently in use is always going to be reset to the proper throttling settings
// even in the event of sequencer failover.
func (l *BatchSubmitter) throttlingLoop(wg *sync.WaitGroup, pendingBytesUpdated chan int64) {
	defer wg.Done()
	l.Log.Info("Starting DA throttling loop")

	retryInterval := 10 * time.Second
	retryTimer := time.NewTimer(retryInterval)
	retryTimer.Stop()

	updateParams := func(pendingBytes int64) {
		retryTimer.Stop()
		ctx, cancel := context.WithTimeout(l.shutdownCtx, l.Config.NetworkTimeout)
		defer cancel()
		cl, err := l.EndpointProvider.EthClient(ctx)
		if err != nil {
			l.Log.Error("Can't reach sequencer execution RPC", "err", err)
			return
		}

		maxTxSize := uint64(0)
		maxBlockSize := l.Config.ThrottleAlwaysBlockSize
		if pendingBytes > int64(l.Config.ThrottleThreshold) {
			l.Log.Warn("Pending bytes over limit, throttling DA", "bytes", pendingBytes, "limit", l.Config.ThrottleThreshold)
			maxTxSize = l.Config.ThrottleTxSize
			if maxBlockSize == 0 || (l.Config.ThrottleBlockSize != 0 && l.Config.ThrottleBlockSize < maxBlockSize) {
				maxBlockSize = l.Config.ThrottleBlockSize
			}
		}
		var (
			success bool
			rpcErr  rpc.Error
		)
		err = cl.Client().CallContext(
			ctx, &success, SetMaxDASizeMethod, hexutil.Uint64(maxTxSize), hexutil.Uint64(maxBlockSize),
		)
		if errors.Is(ctx.Err(), context.Canceled) {
			// If the context was cancelled, our work is done and we expect an error here:
			// So log it quietly and exit.
			l.Log.Debug("DA throttling context cancelled")
			return
		}
		if errors.As(err, &rpcErr) && eth.ErrorCode(rpcErr.ErrorCode()).IsGenericRPCError() {
			l.Log.Error("SetMaxDASize rpc unavailable or broken, shutting down. Either enable it or disable throttling.", "err", err)
			// We'd probably hit this error right after startup, so a short shutdown duration should suffice.
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			// Call StopBatchSubmitting in another goroutine to avoid deadlock.
			go func() {
				// Always returns nil. An error is only returned to expose this function as an RPC.
				_ = l.StopBatchSubmitting(ctx)
			}()
			return
		} else if err != nil {
			l.Log.Error("SetMaxDASize rpc failed, retrying.", "err", err)
			retryTimer.Reset(retryInterval)
			return
		}
		if !success {
			l.Log.Error("Result of SetMaxDASize was false, retrying.")
			retryTimer.Reset(retryInterval)
		}
	}

	cachedPendingBytes := int64(0)

	for {
		select {
		case pendingBytes, ok := <-pendingBytesUpdated:
			if !ok {
				// If the channel was closed, this is our signal to exit
				l.Log.Info("throttlingLoop returning")
				return
			}
			updateParams(pendingBytes)
			cachedPendingBytes = pendingBytes
		case <-l.ActiveSeqChanged:
			updateParams(cachedPendingBytes)
		case <-retryTimer.C:
			updateParams(cachedPendingBytes)
		}
	}
}

func (l *BatchSubmitter) waitNodeSyncAndClearState() {
	// Wait for any in flight transactions
	// to be ingested by the node before
	// we start loading blocks again.
	err := l.waitNodeSync()
	if err != nil {
		l.Log.Warn("error waiting for node sync", "err", err)
	}
	l.clearState(l.shutdownCtx)
}

// waitNodeSync Check to see if there was a batcher tx sent recently that
// still needs more block confirmations before being considered finalized
func (l *BatchSubmitter) waitNodeSync() error {
	ctx := l.shutdownCtx
	rollupClient, err := l.EndpointProvider.RollupClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to get rollup client: %w", err)
	}

	cCtx, cancel := context.WithTimeout(ctx, l.Config.NetworkTimeout)
	defer cancel()

	l1Tip, _, err := l.l1Tip(cCtx)
	if err != nil {
		return fmt.Errorf("failed to retrieve l1 tip: %w", err)
	}

	l1TargetBlock := l1Tip.Number
	if l.Config.CheckRecentTxsDepth != 0 {
		l.Log.Info("Checking for recently submitted batcher transactions on L1")
		recentBlock, found, err := eth.CheckRecentTxs(cCtx, l.L1Client, l.Config.CheckRecentTxsDepth, l.Txmgr.From())
		if err != nil {
			return fmt.Errorf("failed checking recent batcher txs: %w", err)
		}
		l.Log.Info("Checked for recently submitted batcher transactions on L1",
			"l1_head", l1Tip, "l1_recent", recentBlock, "found", found)
		l1TargetBlock = recentBlock
	}

	return dial.WaitRollupSync(l.shutdownCtx, l.Log, rollupClient, l1TargetBlock, time.Second*12)
}

// publishStateToL1 queues up all pending TxData to be published to the L1, returning when there is no more data to
// queue for publishing or if there was an error queing the data.
func (l *BatchSubmitter) publishStateToL1(ctx context.Context, queue *txmgr.Queue[txRef], receiptsCh chan txmgr.TxReceipt[txRef], daGroup *errgroup.Group) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		// if the txmgr is closed, we stop the transaction sending
		if l.Txmgr.IsClosed() {
			l.Log.Info("Txmgr is closed, aborting state publishing")
			return
		}
		if !l.checkTxpool(queue, receiptsCh) {
			l.Log.Info("txpool state is not good, aborting state publishing")
			return
		}

		err := l.publishTxToL1(ctx, queue, receiptsCh, daGroup)
		if err != nil {
			if err != io.EOF {
				l.Log.Error("Error publishing tx to l1", "err", err)
			}
			return
		}
	}
}

// clearState clears the state of the channel manager
func (l *BatchSubmitter) clearState(ctx context.Context) {
	l.Log.Info("Clearing state")
	defer l.Log.Info("State cleared")

	clearStateWithL1Origin := func() bool {
		l1SafeOrigin, err := l.safeL1Origin(ctx)
		if err != nil {
			l.Log.Warn("Failed to query L1 safe origin, will retry", "err", err)
			return false
		} else {
			l.Log.Info("Clearing state with safe L1 origin", "origin", l1SafeOrigin)
			l.channelMgrMutex.Lock()
			defer l.channelMgrMutex.Unlock()
			l.channelMgr.Clear(l1SafeOrigin)
			return true
		}
	}

	// Attempt to set the L1 safe origin and clear the state, if fetching fails -- fall through to an infinite retry
	if clearStateWithL1Origin() {
		return
	}

	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			if clearStateWithL1Origin() {
				return
			}
		case <-ctx.Done():
			l.Log.Warn("Clearing state cancelled")
			l.channelMgrMutex.Lock()
			defer l.channelMgrMutex.Unlock()
			l.channelMgr.Clear(eth.BlockID{})
			return
		}
	}
}

// publishTxToL1 submits a single state tx to the L1
func (l *BatchSubmitter) publishTxToL1(ctx context.Context, queue *txmgr.Queue[txRef], receiptsCh chan txmgr.TxReceipt[txRef], daGroup *errgroup.Group) error {
	// send all available transactions
	l1tip, isPectra, err := l.l1Tip(ctx)
	if err != nil {
		l.Log.Error("Failed to query L1 tip", "err", err)
		return err
	}
	l.Metr.RecordLatestL1Block(l1tip)

	// Collect next transaction data. This pulls data out of the channel, so we need to make sure
	// to put it back if ever da or txmgr requests fail, by calling l.recordFailedDARequest/recordFailedTx.
	l.channelMgrMutex.Lock()
	txdata, err := l.channelMgr.TxData(l1tip.ID(), isPectra)
	l.channelMgrMutex.Unlock()

	if err == io.EOF {
		l.Log.Trace("No transaction data available")
		return err
	} else if err != nil {
		l.Log.Error("Unable to get tx data", "err", err)
		return err
	}

	if err = l.sendTransaction(txdata, queue, receiptsCh, daGroup); err != nil {
		return fmt.Errorf("BatchSubmitter.sendTransaction failed: %w", err)
	}
	return nil
}

func (l *BatchSubmitter) safeL1Origin(ctx context.Context) (eth.BlockID, error) {
	c, err := l.EndpointProvider.RollupClient(ctx)
	if err != nil {
		log.Error("Failed to get rollup client", "err", err)
		return eth.BlockID{}, fmt.Errorf("safe l1 origin: error getting rollup client: %w", err)
	}

	cCtx, cancel := context.WithTimeout(ctx, l.Config.NetworkTimeout)
	defer cancel()

	status, err := c.SyncStatus(cCtx)
	if err != nil {
		log.Error("Failed to get sync status", "err", err)
		return eth.BlockID{}, fmt.Errorf("safe l1 origin: error getting sync status: %w", err)
	}

	// If the safe L2 block origin is 0, we are at the genesis block and should use the L1 origin from the rollup config.
	if status.LocalSafeL2.L1Origin.Number == 0 {
		return l.RollupConfig.Genesis.L1, nil
	}

	return status.LocalSafeL2.L1Origin, nil
}

// cancelBlockingTx creates an empty transaction of appropriate type to cancel out the incompatible
// transaction stuck in the txpool. In the future we might send an actual batch transaction instead
// of an empty one to avoid wasting the tx fee.
func (l *BatchSubmitter) cancelBlockingTx(queue *txmgr.Queue[txRef], receiptsCh chan txmgr.TxReceipt[txRef], isBlockedBlob bool) {
	var candidate *txmgr.TxCandidate
	var err error
	if isBlockedBlob {
		candidate = l.calldataTxCandidate([]byte{})
	} else if candidate, err = l.blobTxCandidate(emptyTxData); err != nil {
		panic(err) // this error should not happen
	}
	l.Log.Warn("sending a cancellation transaction to unblock txpool", "blocked_blob", isBlockedBlob)
	l.sendTx(txData{}, true, candidate, queue, receiptsCh)
}

// publishToAltDAAndL1 posts the txdata to the DA Provider and then sends the commitment to L1.
func (l *BatchSubmitter) publishToAltDAAndL1(txdata txData, queue *txmgr.Queue[txRef], receiptsCh chan txmgr.TxReceipt[txRef], daGroup *errgroup.Group) {
	// sanity checks
	if nf := len(txdata.frames); nf != 1 {
		l.Log.Crit("Unexpected number of frames in calldata tx", "num_frames", nf)
	}
	if txdata.asBlob {
		l.Log.Crit("Unexpected blob txdata with AltDA enabled")
	}

	// when posting txdata to an external DA Provider, we use a goroutine to avoid blocking the main loop
	// since it may take a while for the request to return.
	goroutineSpawned := daGroup.TryGo(func() error {
		// TODO: probably shouldn't be using the global shutdownCtx here, see https://go.dev/blog/context-and-structs
		// but sendTransaction receives l.killCtx as an argument, which currently is only canceled after waiting for the main loop
		// to exit, which would wait on this DA call to finish, which would take a long time.
		// So we prefer to mimic the behavior of txmgr and cancel all pending DA/txmgr requests when the batcher is stopped.
		comm, err := l.AltDA.SetInput(l.shutdownCtx, txdata.CallData())
		if err != nil {
			// Don't log context cancelled events because they are expected,
			// and can happen after tests complete which causes a panic.
			if errors.Is(err, context.Canceled) {
				l.recordFailedDARequest(txdata.ID(), nil)
			} else {
				l.Log.Error("Failed to post input to Alt DA", "error", err)
				// requeue frame if we fail to post to the DA Provider so it can be retried
				// note: this assumes that the da server caches requests, otherwise it might lead to resubmissions of the blobs
				l.recordFailedDARequest(txdata.ID(), err)
			}
			return nil
		}
		l.Log.Info("Set altda input", "commitment", comm, "tx", txdata.ID())
		candidate := l.calldataTxCandidate(comm.TxData())
		l.sendTx(txdata, false, candidate, queue, receiptsCh)
		return nil
	})
	if !goroutineSpawned {
		// We couldn't start the goroutine because the errgroup.Group limit
		// is already reached. Since we can't send the txdata, we have to
		// return it for later processing. We use nil error to skip error logging.
		l.recordFailedDARequest(txdata.ID(), nil)
	}
}

// sendTransaction creates & queues for sending a transaction to the batch inbox address with the given `txData`.
// This call will block if the txmgr queue is at the  max-pending limit.
// The method will block if the queue's MaxPendingTransactions is exceeded.
func (l *BatchSubmitter) sendTransaction(txdata txData, queue *txmgr.Queue[txRef], receiptsCh chan txmgr.TxReceipt[txRef], daGroup *errgroup.Group) error {
	var err error

	// if Alt DA is enabled we post the txdata to the DA Provider and replace it with the commitment.
	if l.Config.UseAltDA {
		l.publishToAltDAAndL1(txdata, queue, receiptsCh, daGroup)
		// we return nil to allow publishStateToL1 to keep processing the next txdata
		return nil
	}

	var candidate *txmgr.TxCandidate
	if txdata.asBlob {
		if candidate, err = l.blobTxCandidate(txdata); err != nil {
			// We could potentially fall through and try a calldata tx instead, but this would
			// likely result in the chain spending more in gas fees than it is tuned for, so best
			// to just fail. We do not expect this error to trigger unless there is a serious bug
			// or configuration issue.
			return fmt.Errorf("could not create blob tx candidate: %w", err)
		}
	} else {
		// sanity check
		if nf := len(txdata.frames); nf != 1 {
			l.Log.Crit("Unexpected number of frames in calldata tx", "num_frames", nf)
		}
		candidate = l.calldataTxCandidate(txdata.CallData())
	}

	l.sendTx(txdata, false, candidate, queue, receiptsCh)
	return nil
}

type TxSender[T any] interface {
	Send(id T, candidate txmgr.TxCandidate, receiptCh chan txmgr.TxReceipt[T])
}

// sendTx uses the txmgr queue to send the given transaction candidate after setting its
// gaslimit. It will block if the txmgr queue has reached its MaxPendingTransactions limit.
func (l *BatchSubmitter) sendTx(txdata txData, isCancel bool, candidate *txmgr.TxCandidate, queue TxSender[txRef], receiptsCh chan txmgr.TxReceipt[txRef]) {
	floorDataGas, err := core.FloorDataGas(candidate.TxData)
	if err != nil {
		// We log instead of return an error here because the txmgr will do its own gas estimation.
		l.Log.Warn("Failed to calculate floor data gas", "err", err)
	} else {
		candidate.GasLimit = floorDataGas
	}

	queue.Send(txRef{id: txdata.ID(), isCancel: isCancel, isBlob: txdata.asBlob}, *candidate, receiptsCh)
}

func (l *BatchSubmitter) blobTxCandidate(data txData) (*txmgr.TxCandidate, error) {
	blobs, err := data.Blobs()
	if err != nil {
		return nil, fmt.Errorf("generating blobs for tx data: %w", err)
	}
	size := data.Len()
	lastSize := len(data.frames[len(data.frames)-1].data)
	l.Log.Info("Building Blob transaction candidate",
		"size", size, "last_size", lastSize, "num_blobs", len(blobs))
	l.Metr.RecordBlobUsedBytes(lastSize)
	return &txmgr.TxCandidate{
		To:    &l.RollupConfig.BatchInboxAddress,
		Blobs: blobs,
	}, nil
}

func (l *BatchSubmitter) calldataTxCandidate(data []byte) *txmgr.TxCandidate {
	l.Log.Info("Building Calldata transaction candidate", "size", len(data))
	return &txmgr.TxCandidate{
		To:     &l.RollupConfig.BatchInboxAddress,
		TxData: data,
	}
}

func (l *BatchSubmitter) handleReceipt(r txmgr.TxReceipt[txRef]) {
	// Record TX Status
	if r.Err != nil {
		l.recordFailedTx(r.ID.id, r.Err)
	} else if r.Receipt != nil {
		l.recordConfirmedTx(r.ID.id, r.Receipt)
	}
	// Both r.Err and r.Receipt can be nil, in which case we do nothing.
}

func (l *BatchSubmitter) recordFailedDARequest(id txID, err error) {
	l.channelMgrMutex.Lock()
	defer l.channelMgrMutex.Unlock()
	if err != nil {
		l.Log.Warn("DA request failed", logFields(id, err)...)
	}
	l.channelMgr.TxFailed(id)
}

func (l *BatchSubmitter) recordFailedTx(id txID, err error) {
	l.channelMgrMutex.Lock()
	defer l.channelMgrMutex.Unlock()
	l.Log.Warn("Transaction failed to send", logFields(id, err)...)
	l.channelMgr.TxFailed(id)
}

func (l *BatchSubmitter) recordConfirmedTx(id txID, receipt *types.Receipt) {
	l.channelMgrMutex.Lock()
	defer l.channelMgrMutex.Unlock()
	l.Log.Info("Transaction confirmed", logFields(id, receipt)...)
	l1block := eth.ReceiptBlockID(receipt)
	l.channelMgr.TxConfirmed(id, l1block)
}

// l1Tip gets the current L1 tip as a L1BlockRef. The passed context is assumed
// to be a lifetime context, so it is internally wrapped with a network timeout.
// It also returns a boolean indicating if the tip is from a Pectra chain.
func (l *BatchSubmitter) l1Tip(ctx context.Context) (eth.L1BlockRef, bool, error) {
	tctx, cancel := context.WithTimeout(ctx, l.Config.NetworkTimeout)
	defer cancel()
	head, err := l.L1Client.HeaderByNumber(tctx, nil)
	if err != nil {
		return eth.L1BlockRef{}, false, fmt.Errorf("getting latest L1 block: %w", err)
	}
	isPectra := head.RequestsHash != nil // See https://eips.ethereum.org/EIPS/eip-7685
	return eth.InfoToL1BlockRef(eth.HeaderBlockInfo(head)), isPectra, nil
}

func (l *BatchSubmitter) checkTxpool(queue *txmgr.Queue[txRef], receiptsCh chan txmgr.TxReceipt[txRef]) bool {
	l.txpoolMutex.Lock()
	if l.txpoolState == TxpoolBlocked {
		// txpoolState is set to Blocked only if Send() is returning
		// ErrAlreadyReserved. In this case, the TxMgr nonce should be reset to nil,
		// allowing us to send a cancellation transaction.
		l.txpoolState = TxpoolCancelPending
		isBlob := l.txpoolBlockedBlob
		l.txpoolMutex.Unlock()
		l.cancelBlockingTx(queue, receiptsCh, isBlob)
		return false
	}
	r := l.txpoolState == TxpoolGood
	l.txpoolMutex.Unlock()
	return r
}

func logFields(xs ...any) (fs []any) {
	for _, x := range xs {
		switch v := x.(type) {
		case txID:
			fs = append(fs, "tx_id", v.String())
		case *types.Receipt:
			fs = append(fs, "tx", v.TxHash, "block", eth.ReceiptBlockID(v))
		case error:
			fs = append(fs, "err", v)
		default:
			fs = append(fs, "ERROR", fmt.Sprintf("logFields: unknown type: %T", x))
		}
	}
	return fs
}
