const optimism = require('@tokamak-network/titan2-sdk')
const ethers = require('ethers')

const privateKey = process.env.PRIVATE_KEY
const TON = process.env.TON

const l1Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L1_URL
)
const l2Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L2_URL
)
const l1Wallet = new ethers.Wallet(privateKey, l1Provider)
const l2Wallet = new ethers.Wallet(privateKey, l2Provider)

const erc20ABI = [
  {
    inputs: [
      { internalType: 'address', name: '_spender', type: 'address' },
      { internalType: 'uint256', name: '_value', type: 'uint256' },
    ],
    name: 'approve',
    outputs: [{ internalType: 'bool', name: '', type: 'bool' }],
    stateMutability: 'nonpayable',
    type: 'function',
  },
  {
    constant: true,
    inputs: [{ name: '_owner', type: 'address' }],
    name: 'balanceOf',
    outputs: [{ name: 'balance', type: 'uint256' }],
    type: 'function',
  },
  {
    inputs: [{ internalType: 'uint256', name: 'amount', type: 'uint256' }],
    name: 'faucet',
    outputs: [],
    stateMutability: 'nonpayable',
    type: 'function',
  },
]
const tonContract = new ethers.Contract(TON, erc20ABI, l1Wallet)

const ETH = '0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000'

const zeroAddr = '0x'.padEnd(42, '0')
const l1Contracts = {
  StateCommitmentChain: zeroAddr,
  CanonicalTransactionChain: zeroAddr,
  BondManager: zeroAddr,
  AddressManager: process.env.ADDRESS_MANAGER, // Lib_AddressManager.json
  L1CrossDomainMessenger: process.env.L1_CROSS_DOMAIN_MESSENGER, // Proxy__OVM_L1CrossDomainMessenger.json
  L1StandardBridge: process.env.L1_STANDARD_BRIDGE, // Proxy__OVM_L1StandardBridge.json
  OptimismPortal: process.env.OPTIMISM_PORTAL, // OptimismPortalProxy.json
  L2OutputOracle: process.env.L1_OUTPUT_ORACLE, // L2OutputOracleProxy.json
}

const bridges = {
  Standard: {
    l1Bridge: l1Contracts.L1StandardBridge,
    l2Bridge: '0x4200000000000000000000000000000000000010',
    Adapter: optimism.TONBridgeAdapter,
  },
}

const messenger = new optimism.CrossChainMessenger({
  bedrock: true,
  contracts: {
    l1: l1Contracts,
  },
  bridges,
  l1ChainId: process.env.L1_CHAIN_ID,
  l2ChainId: process.env.L2_CHAIN_ID,
  l1SignerOrProvider: l1Wallet,
  l2SignerOrProvider: l2Wallet,
})

const depositTON = async (amount) => {
  console.log('Deposit TON:', amount)

  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())

  const l1TONBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance:', l1TONBalance.toString())

  const approveTx = await messenger.approveERC20(TON, ETH, amount)
  await approveTx.wait()
  console.log('approveTx:', approveTx.hash)

  const depositTx = await messenger.depositERC20(TON, ETH, amount)
  await depositTx.wait()
  console.log('depositTx:', depositTx.hash)

  await messenger.waitForMessageStatus(
    depositTx.hash,
    optimism.MessageStatus.RELAYED
  )

  l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())
}

const withdrawTON = async (amount) => {
  console.log('Withdraw TON:', amount)
  let tonBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance: ', tonBalance.toString())

  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())

  const withdrawal = await messenger.withdrawETH(amount)
  const withdrawalTx = await withdrawal.wait()
  console.log(
    'withdrawal Tx:',
    withdrawalTx.transactionHash,
    ' Block',
    withdrawalTx.blockNumber,
    ' hash',
    withdrawal.hash
  )

  l2Balance = await l2Wallet.getBalance()
  console.log('Updated l2 native balance:', l2Balance.toString())

  // // Check ready for prove
  await messenger.waitForMessageStatus(
    withdrawalTx.transactionHash,
    optimism.MessageStatus.READY_TO_PROVE
  )

  console.log('Prove the message')
  const proveTx = await messenger.proveMessage(withdrawalTx.transactionHash)
  const proveReceipt = await proveTx.wait(3)
  console.log('Proved the message: ', proveReceipt.transactionHash)

  const finalizeInterval = setInterval(async () => {
    const currentStatus = await messenger.getMessageStatus(
      withdrawalTx.transactionHash
    )
    console.log('Message status: ', currentStatus)
  }, 3000)

  try {
    await messenger.waitForMessageStatus(
      withdrawalTx.transactionHash,
      optimism.MessageStatus.READY_FOR_RELAY
    )
  } finally {
    clearInterval(finalizeInterval)
  }

  tonBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance: ', tonBalance.toString())

  const tx = await messenger.finalizeMessage(withdrawalTx.transactionHash)
  const receipt = await tx.wait()
  console.log('Finalized message tx', receipt.transactionHash)
  console.log('Finalized withdrawal')

  tonBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance: ', tonBalance.toString())
}

const main = async () => {
  console.log('L1 Address ', l1Wallet.address)
  console.log('L2 Address ', l2Wallet.address)

  await depositTON(5000)
  await withdrawTON(4000)
}

main()
