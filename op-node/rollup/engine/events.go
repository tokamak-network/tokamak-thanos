package engine

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum-optimism/optimism/op-node/rollup/event"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

// ReplaceBlockSource is a magic value for the "Source" attribute,
// used when a L2 block is a replacement of an invalidated block.
// After the replacement has been processed, a reset is performed to derive the next L2 blocks.
var ReplaceBlockSource = eth.L1BlockRef{
	Hash:       common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	Number:     ^uint64(0),
	ParentHash: common.Hash{},
	Time:       0,
}

type Metrics interface {
	CountSequencedTxsInBlock(txns int, deposits int)

	RecordSequencerBuildingDiffTime(duration time.Duration)
	RecordSequencerSealingTime(duration time.Duration)
}

// ForkchoiceRequestEvent signals to the engine that it should emit an artificial
// forkchoice-update event, to signal the latest forkchoice to other derivers.
// This helps decouple derivers from the actual engine state,
// while also not making the derivers wait for a forkchoice update at random.
type ForkchoiceRequestEvent struct{}

func (ev ForkchoiceRequestEvent) String() string {
	return "forkchoice-request"
}

type ForkchoiceUpdateEvent struct {
	UnsafeL2Head, SafeL2Head, FinalizedL2Head eth.L2BlockRef
}

func (ev ForkchoiceUpdateEvent) String() string {
	return "forkchoice-update"
}

// PromoteUnsafeEvent signals that the given block may now become a canonical unsafe block.
// This is pre-forkchoice update; the change may not be reflected yet in the EL.
// Note that the legacy pre-event-refactor code-path (processing P2P blocks) does fire this,
// but manually, duplicate with the newer events processing code-path.
// See EngineController.InsertUnsafePayload.
type PromoteUnsafeEvent struct {
	Ref eth.L2BlockRef
}

func (ev PromoteUnsafeEvent) String() string {
	return "promote-unsafe"
}

// RequestCrossUnsafeEvent signals that a CrossUnsafeUpdateEvent is needed.
type RequestCrossUnsafeEvent struct{}

func (ev RequestCrossUnsafeEvent) String() string {
	return "request-cross-unsafe"
}

// UnsafeUpdateEvent signals that the given block is now considered safe.
// This is pre-forkchoice update; the change may not be reflected yet in the EL.
type UnsafeUpdateEvent struct {
	Ref eth.L2BlockRef
}

func (ev UnsafeUpdateEvent) String() string {
	return "unsafe-update"
}

// PromoteCrossUnsafeEvent signals that the given block may be promoted to cross-unsafe.
type PromoteCrossUnsafeEvent struct {
	Ref eth.L2BlockRef
}

func (ev PromoteCrossUnsafeEvent) String() string {
	return "promote-cross-unsafe"
}

// CrossUnsafeUpdateEvent signals that the given block is now considered cross-unsafe.
type CrossUnsafeUpdateEvent struct {
	CrossUnsafe eth.L2BlockRef
	LocalUnsafe eth.L2BlockRef
}

func (ev CrossUnsafeUpdateEvent) String() string {
	return "cross-unsafe-update"
}

type PendingSafeUpdateEvent struct {
	PendingSafe eth.L2BlockRef
	Unsafe      eth.L2BlockRef // tip, added to the signal, to determine if there are existing blocks to consolidate
}

func (ev PendingSafeUpdateEvent) String() string {
	return "pending-safe-update"
}

// PromotePendingSafeEvent signals that a block can be marked as pending-safe, and/or safe.
type PromotePendingSafeEvent struct {
	Ref        eth.L2BlockRef
	Concluding bool // Concludes the pending phase, so can be promoted to (local) safe
	Source     eth.L1BlockRef
}

func (ev PromotePendingSafeEvent) String() string {
	return "promote-pending-safe"
}

// PromoteLocalSafeEvent signals that a block can be promoted to local-safe.
type PromoteLocalSafeEvent struct {
	Ref    eth.L2BlockRef
	Source eth.L1BlockRef
}

func (ev PromoteLocalSafeEvent) String() string {
	return "promote-local-safe"
}

// RequestCrossSafeEvent signals that a CrossSafeUpdate is needed.
type RequestCrossSafeEvent struct{}

func (ev RequestCrossSafeEvent) String() string {
	return "request-cross-safe-update"
}

type CrossSafeUpdateEvent struct {
	CrossSafe eth.L2BlockRef
	LocalSafe eth.L2BlockRef
}

func (ev CrossSafeUpdateEvent) String() string {
	return "cross-safe-update"
}

// LocalSafeUpdateEvent signals that a block is now considered to be local-safe.
type LocalSafeUpdateEvent struct {
	Ref    eth.L2BlockRef
	Source eth.L1BlockRef
}

func (ev LocalSafeUpdateEvent) String() string {
	return "local-safe-update"
}

// PromoteSafeEvent signals that a block can be promoted to cross-safe.
type PromoteSafeEvent struct {
	Ref    eth.L2BlockRef
	Source eth.L1BlockRef
}

func (ev PromoteSafeEvent) String() string {
	return "promote-safe"
}

// SafeDerivedEvent signals that a block was determined to be safe, and derived from the given L1 block.
// This is signaled upon successful processing of PromoteSafeEvent.
type SafeDerivedEvent struct {
	Safe   eth.L2BlockRef
	Source eth.L1BlockRef
}

func (ev SafeDerivedEvent) String() string {
	return "safe-derived"
}

// ProcessAttributesEvent signals to immediately process the attributes.
type ProcessAttributesEvent struct {
	Attributes *derive.AttributesWithParent
}

func (ev ProcessAttributesEvent) String() string {
	return "process-attributes"
}

type PendingSafeRequestEvent struct{}

func (ev PendingSafeRequestEvent) String() string {
	return "pending-safe-request"
}

type ProcessUnsafePayloadEvent struct {
	Envelope *eth.ExecutionPayloadEnvelope
}

func (ev ProcessUnsafePayloadEvent) String() string {
	return "process-unsafe-payload"
}

type TryBackupUnsafeReorgEvent struct{}

func (ev TryBackupUnsafeReorgEvent) String() string {
	return "try-backup-unsafe-reorg"
}

type TryUpdateEngineEvent struct {
	// These fields will be zero-value (BuildStarted,InsertStarted=time.Time{}, Envelope=nil) if
	// this event is emitted outside of engineDeriver.onPayloadSuccess
	BuildStarted  time.Time
	InsertStarted time.Time
	Envelope      *eth.ExecutionPayloadEnvelope
}

func (ev TryUpdateEngineEvent) String() string {
	return "try-update-engine"
}

// Checks for the existence of the Envelope field, which is only
// added by the PayloadSuccessEvent
func (ev TryUpdateEngineEvent) triggeredByPayloadSuccess() bool {
	return ev.Envelope != nil
}

// Returns key/value pairs that can be logged and are useful for plotting
// block build/insert time as a way to measure performance.
func (ev TryUpdateEngineEvent) getBlockProcessingMetrics() []interface{} {
	fcuFinish := time.Now()
	payload := ev.Envelope.ExecutionPayload

	logValues := []interface{}{
		"hash", payload.BlockHash,
		"number", uint64(payload.BlockNumber),
		"state_root", payload.StateRoot,
		"timestamp", uint64(payload.Timestamp),
		"parent", payload.ParentHash,
		"prev_randao", payload.PrevRandao,
		"fee_recipient", payload.FeeRecipient,
		"txs", len(payload.Transactions),
	}

	var totalTime time.Duration
	var mgasps float64
	if !ev.BuildStarted.IsZero() {
		totalTime = fcuFinish.Sub(ev.BuildStarted)
		logValues = append(logValues,
			"build_time", common.PrettyDuration(ev.InsertStarted.Sub(ev.BuildStarted)),
			"insert_time", common.PrettyDuration(fcuFinish.Sub(ev.InsertStarted)),
		)
	} else if !ev.InsertStarted.IsZero() {
		totalTime = fcuFinish.Sub(ev.InsertStarted)
	}

	// Avoid divide-by-zero for mgasps
	if totalTime > 0 {
		mgasps = float64(payload.GasUsed) * 1000 / float64(totalTime)
	}

	logValues = append(logValues,
		"total_time", common.PrettyDuration(totalTime),
		"mgas", float64(payload.GasUsed)/1000000,
		"mgasps", mgasps,
	)

	return logValues
}

type EngineResetConfirmedEvent struct {
	LocalUnsafe eth.L2BlockRef
	CrossUnsafe eth.L2BlockRef
	LocalSafe   eth.L2BlockRef
	CrossSafe   eth.L2BlockRef
	Finalized   eth.L2BlockRef
}

func (ev EngineResetConfirmedEvent) String() string {
	return "engine-reset-confirmed"
}

// PromoteFinalizedEvent signals that a block can be marked as finalized.
type PromoteFinalizedEvent struct {
	Ref eth.L2BlockRef
}

func (ev PromoteFinalizedEvent) String() string {
	return "promote-finalized"
}

// FinalizedUpdateEvent signals that a block has been marked as finalized.
type FinalizedUpdateEvent struct {
	Ref eth.L2BlockRef
}

func (ev FinalizedUpdateEvent) String() string {
	return "finalized-update"
}

// RequestFinalizedUpdateEvent signals that a FinalizedUpdateEvent is needed.
type RequestFinalizedUpdateEvent struct{}

func (ev RequestFinalizedUpdateEvent) String() string {
	return "request-finalized-update"
}

// CrossUpdateRequestEvent triggers update events to be emitted, repeating the current state.
type CrossUpdateRequestEvent struct {
	CrossUnsafe bool
	CrossSafe   bool
}

func (ev CrossUpdateRequestEvent) String() string {
	return "cross-update-request"
}

// InteropInvalidateBlockEvent is emitted when a block needs to be invalidated, and a replacement is needed.
type InteropInvalidateBlockEvent struct {
	Invalidated eth.BlockRef
	Attributes  *derive.AttributesWithParent
}

func (ev InteropInvalidateBlockEvent) String() string {
	return "interop-invalidate-block"
}

// InteropReplacedBlockEvent is emitted when a replacement is done.
type InteropReplacedBlockEvent struct {
	Ref      eth.BlockRef
	Envelope *eth.ExecutionPayloadEnvelope
}

func (ev InteropReplacedBlockEvent) String() string {
	return "interop-replaced-block"
}

type EngDeriver struct {
	metrics Metrics

	log     log.Logger
	cfg     *rollup.Config
	ec      *EngineController
	ctx     context.Context
	emitter event.Emitter
}

var _ event.Deriver = (*EngDeriver)(nil)

func NewEngDeriver(log log.Logger, ctx context.Context, cfg *rollup.Config,
	metrics Metrics, ec *EngineController,
) *EngDeriver {
	return &EngDeriver{
		log:     log,
		cfg:     cfg,
		ec:      ec,
		ctx:     ctx,
		metrics: metrics,
	}
}

func (d *EngDeriver) AttachEmitter(em event.Emitter) {
	d.emitter = em
}

func (d *EngDeriver) OnEvent(ev event.Event) bool {
	switch x := ev.(type) {
	case TryBackupUnsafeReorgEvent:
		// If we don't need to call FCU to restore unsafeHead using backupUnsafe, keep going b/c
		// this was a no-op(except correcting invalid state when backupUnsafe is empty but TryBackupUnsafeReorg called).
		fcuCalled, err := d.ec.TryBackupUnsafeReorg(d.ctx)
		// Dealing with legacy here: it used to skip over the error-handling if fcuCalled was false.
		// But that combination is not actually a code-path in TryBackupUnsafeReorg.
		// We should drop fcuCalled, and make the function emit events directly,
		// once there are no more synchronous callers.
		if !fcuCalled && err != nil {
			d.log.Crit("unexpected TryBackupUnsafeReorg error after no FCU call", "err", err)
		}
		if err != nil {
			// If we needed to perform a network call, then we should yield even if we did not encounter an error.
			if errors.Is(err, derive.ErrReset) {
				d.emitter.Emit(rollup.ResetEvent{Err: err})
			} else if errors.Is(err, derive.ErrTemporary) {
				d.emitter.Emit(rollup.EngineTemporaryErrorEvent{Err: err})
			} else {
				d.emitter.Emit(rollup.CriticalErrorEvent{Err: fmt.Errorf("unexpected TryBackupUnsafeReorg error type: %w", err)})
			}
		}
	case TryUpdateEngineEvent:
		// If we don't need to call FCU, keep going b/c this was a no-op. If we needed to
		// perform a network call, then we should yield even if we did not encounter an error.
		if err := d.ec.TryUpdateEngine(d.ctx); err != nil && !errors.Is(err, ErrNoFCUNeeded) {
			if errors.Is(err, derive.ErrReset) {
				d.emitter.Emit(rollup.ResetEvent{Err: err})
			} else if errors.Is(err, derive.ErrTemporary) {
				d.emitter.Emit(rollup.EngineTemporaryErrorEvent{Err: err})
			} else {
				d.emitter.Emit(rollup.CriticalErrorEvent{Err: fmt.Errorf("unexpected TryUpdateEngine error type: %w", err)})
			}
		} else if x.triggeredByPayloadSuccess() {
			logValues := x.getBlockProcessingMetrics()
			d.log.Info("Inserted new L2 unsafe block", logValues...)
		}
	case ProcessUnsafePayloadEvent:
		ref, err := derive.PayloadToBlockRef(d.cfg, x.Envelope.ExecutionPayload)
		if err != nil {
			d.log.Error("failed to decode L2 block ref from payload", "err", err)
			return true
		}
		if err := d.ec.InsertUnsafePayload(d.ctx, x.Envelope, ref); err != nil {
			d.log.Info("failed to insert payload", "ref", ref,
				"txs", len(x.Envelope.ExecutionPayload.Transactions), "err", err)
			// yes, duplicate error-handling. After all derivers are interacting with the engine
			// through events, we can drop the engine-controller interface:
			// unify the events handler with the engine-controller,
			// remove a lot of code, and not do this error translation.
			if errors.Is(err, derive.ErrReset) {
				d.emitter.Emit(rollup.ResetEvent{Err: err})
			} else if errors.Is(err, derive.ErrTemporary) {
				d.emitter.Emit(rollup.EngineTemporaryErrorEvent{Err: err})
			} else {
				d.emitter.Emit(rollup.CriticalErrorEvent{Err: fmt.Errorf("unexpected InsertUnsafePayload error type: %w", err)})
			}
		} else {
			d.log.Info("successfully processed payload", "ref", ref, "txs", len(x.Envelope.ExecutionPayload.Transactions))
		}
	case ForkchoiceRequestEvent:
		d.emitter.Emit(ForkchoiceUpdateEvent{
			UnsafeL2Head:    d.ec.UnsafeL2Head(),
			SafeL2Head:      d.ec.SafeL2Head(),
			FinalizedL2Head: d.ec.Finalized(),
		})
	case rollup.ForceResetEvent:
		ForceEngineReset(d.ec, x)

		// Time to apply the changes to the underlying engine
		d.emitter.Emit(TryUpdateEngineEvent{})

		v := EngineResetConfirmedEvent{
			LocalUnsafe: d.ec.LocalSafeL2Head(),
			CrossUnsafe: d.ec.CrossUnsafeL2Head(),
			LocalSafe:   d.ec.LocalSafeL2Head(),
			CrossSafe:   d.ec.SafeL2Head(),
			Finalized:   d.ec.Finalized(),
		}
		// We do not emit the original event values, since those might not be set (optional attributes).
		d.emitter.Emit(v)
		d.log.Info("Reset of Engine is completed",
			"local_unsafe", v.LocalUnsafe,
			"cross_unsafe", v.CrossUnsafe,
			"local_safe", v.LocalSafe,
			"cross_safe", v.CrossSafe,
			"finalized", v.Finalized,
		)
	case PromoteUnsafeEvent:
		// Backup unsafeHead when new block is not built on original unsafe head.
		if d.ec.unsafeHead.Number >= x.Ref.Number {
			d.ec.SetBackupUnsafeL2Head(d.ec.unsafeHead, false)
		}
		d.ec.SetUnsafeHead(x.Ref)
		d.emitter.Emit(UnsafeUpdateEvent(x))
	case UnsafeUpdateEvent:
		// pre-interop everything that is local-unsafe is also immediately cross-unsafe.
		if !d.cfg.IsInterop(x.Ref.Time) {
			d.emitter.Emit(PromoteCrossUnsafeEvent(x))
		}
		// Try to apply the forkchoice changes
		d.emitter.Emit(TryUpdateEngineEvent{})
	case PromoteCrossUnsafeEvent:
		d.ec.SetCrossUnsafeHead(x.Ref)
		d.emitter.Emit(CrossUnsafeUpdateEvent{
			CrossUnsafe: x.Ref,
			LocalUnsafe: d.ec.UnsafeL2Head(),
		})
	case RequestCrossUnsafeEvent:
		d.emitter.Emit(CrossUnsafeUpdateEvent{
			CrossUnsafe: d.ec.CrossUnsafeL2Head(),
			LocalUnsafe: d.ec.UnsafeL2Head(),
		})
	case RequestCrossSafeEvent:
		d.emitter.Emit(CrossSafeUpdateEvent{
			CrossSafe: d.ec.SafeL2Head(),
			LocalSafe: d.ec.LocalSafeL2Head(),
		})
	case PendingSafeRequestEvent:
		d.emitter.Emit(PendingSafeUpdateEvent{
			PendingSafe: d.ec.PendingSafeL2Head(),
			Unsafe:      d.ec.UnsafeL2Head(),
		})
	case PromotePendingSafeEvent:
		// Only promote if not already stale.
		// Resets/overwrites happen through engine-resets, not through promotion.
		if x.Ref.Number > d.ec.PendingSafeL2Head().Number {
			d.log.Debug("Updating pending safe", "pending_safe", x.Ref, "local_safe", d.ec.LocalSafeL2Head(), "unsafe", d.ec.UnsafeL2Head(), "concluding", x.Concluding)
			d.ec.SetPendingSafeL2Head(x.Ref)
			d.emitter.Emit(PendingSafeUpdateEvent{
				PendingSafe: d.ec.PendingSafeL2Head(),
				Unsafe:      d.ec.UnsafeL2Head(),
			})
		}
		if x.Concluding && x.Ref.Number > d.ec.LocalSafeL2Head().Number {
			d.emitter.Emit(PromoteLocalSafeEvent{
				Ref:    x.Ref,
				Source: x.Source,
			})
		}
	case PromoteLocalSafeEvent:
		d.log.Debug("Updating local safe", "local_safe", x.Ref, "safe", d.ec.SafeL2Head(), "unsafe", d.ec.UnsafeL2Head())
		d.ec.SetLocalSafeHead(x.Ref)
		d.emitter.Emit(LocalSafeUpdateEvent(x))
	case LocalSafeUpdateEvent:
		// pre-interop everything that is local-safe is also immediately cross-safe.
		if !d.cfg.IsInterop(x.Ref.Time) {
			d.emitter.Emit(PromoteSafeEvent(x))
		}
	case PromoteSafeEvent:
		d.log.Debug("Updating safe", "safe", x.Ref, "unsafe", d.ec.UnsafeL2Head())
		d.ec.SetSafeHead(x.Ref)
		// Finalizer can pick up this safe cross-block now
		d.emitter.Emit(SafeDerivedEvent{Safe: x.Ref, Source: x.Source})
		d.emitter.Emit(CrossSafeUpdateEvent{
			CrossSafe: d.ec.SafeL2Head(),
			LocalSafe: d.ec.LocalSafeL2Head(),
		})
		if x.Ref.Number > d.ec.crossUnsafeHead.Number {
			d.log.Debug("Cross Unsafe Head is stale, updating to match cross safe", "cross_unsafe", d.ec.crossUnsafeHead, "cross_safe", x.Ref)
			d.ec.SetCrossUnsafeHead(x.Ref)
			d.emitter.Emit(CrossUnsafeUpdateEvent{
				CrossUnsafe: x.Ref,
				LocalUnsafe: d.ec.UnsafeL2Head(),
			})
		}
		// Try to apply the forkchoice changes
		d.emitter.Emit(TryUpdateEngineEvent{})
	case PromoteFinalizedEvent:
		if x.Ref.Number < d.ec.Finalized().Number {
			d.log.Error("Cannot rewind finality,", "ref", x.Ref, "finalized", d.ec.Finalized())
			return true
		}
		if x.Ref.Number > d.ec.SafeL2Head().Number {
			d.log.Error("Block must be safe before it can be finalized", "ref", x.Ref, "safe", d.ec.SafeL2Head())
			return true
		}
		d.ec.SetFinalizedHead(x.Ref)
		d.emitter.Emit(FinalizedUpdateEvent(x))
		// Try to apply the forkchoice changes
		d.emitter.Emit(TryUpdateEngineEvent{})
	case RequestFinalizedUpdateEvent:
		d.emitter.Emit(FinalizedUpdateEvent{Ref: d.ec.Finalized()})
	case CrossUpdateRequestEvent:
		if x.CrossUnsafe {
			d.emitter.Emit(CrossUnsafeUpdateEvent{
				CrossUnsafe: d.ec.CrossUnsafeL2Head(),
				LocalUnsafe: d.ec.UnsafeL2Head(),
			})
		}
		if x.CrossSafe {
			d.emitter.Emit(CrossSafeUpdateEvent{
				CrossSafe: d.ec.SafeL2Head(),
				LocalSafe: d.ec.LocalSafeL2Head(),
			})
		}
	case InteropInvalidateBlockEvent:
		d.emitter.Emit(BuildStartEvent{Attributes: x.Attributes})
	case BuildStartEvent:
		d.onBuildStart(x)
	case BuildStartedEvent:
		d.onBuildStarted(x)
	case BuildSealEvent:
		d.onBuildSeal(x)
	case BuildSealedEvent:
		d.onBuildSealed(x)
	case BuildInvalidEvent:
		d.onBuildInvalid(x)
	case BuildCancelEvent:
		d.onBuildCancel(x)
	case PayloadProcessEvent:
		d.onPayloadProcess(x)
	case PayloadSuccessEvent:
		d.onPayloadSuccess(x)
	case PayloadInvalidEvent:
		d.onPayloadInvalid(x)
	default:
		return false
	}
	return true
}

type ResetEngineControl interface {
	SetUnsafeHead(eth.L2BlockRef)
	SetCrossUnsafeHead(ref eth.L2BlockRef)
	SetLocalSafeHead(ref eth.L2BlockRef)
	SetSafeHead(eth.L2BlockRef)
	SetFinalizedHead(eth.L2BlockRef)
	SetBackupUnsafeL2Head(block eth.L2BlockRef, triggerReorg bool)
	SetPendingSafeL2Head(eth.L2BlockRef)
}

func ForceEngineReset(ec ResetEngineControl, x rollup.ForceResetEvent) {
	// local-unsafe is an optional attribute, empty to preserve the existing latest chain
	if x.LocalUnsafe != (eth.L2BlockRef{}) {
		ec.SetUnsafeHead(x.LocalUnsafe)
	}
	// cross-safe is fine to revert back, it does not affect engine logic, just sync-status
	ec.SetCrossUnsafeHead(x.CrossUnsafe)

	// derivation continues at local-safe point
	ec.SetLocalSafeHead(x.LocalSafe)
	ec.SetPendingSafeL2Head(x.LocalSafe)

	// "safe" in RPC terms is cross-safe
	ec.SetSafeHead(x.CrossSafe)

	// finalized head
	ec.SetFinalizedHead(x.Finalized)

	ec.SetBackupUnsafeL2Head(eth.L2BlockRef{}, false)
}
