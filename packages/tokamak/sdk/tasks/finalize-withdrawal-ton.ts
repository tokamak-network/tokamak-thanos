import { task, types } from 'hardhat/config'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { predeploys } from '@eth-optimism/core-utils'
import { BytesLike, ethers } from 'ethers'

import {
  CrossChainMessenger,
  MessageStatus,
  TONBridgeAdapter,
  NumberLike,
} from '../src'

console.log('Setup task...')

const privateKey = process.env.PRIVATE_KEY as BytesLike

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

const ETH = '0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000'

const zeroAddr = '0x'.padEnd(42, '0')

let TON = process.env.TON || ''
let addressManager = process.env.ADDRESS_MANAGER || ''
let l1CrossDomainMessenger = process.env.L1_CROSS_DOMAIN_MESSENGER || ''
let l1StandardBridge = process.env.L1_STANDARD_BRIDGE || ''
let optimismPortal = process.env.OPTIMISM_PORTAL || ''
let l2OutputOracle = process.env.L1_OUTPUT_ORACLE || ''

const recipientAddress = '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266';  // 받는 주소를 새로운 주소로 변경
const amountToTransfer = ethers.utils.parseUnits('10', 1);  // 전송할 토큰의 양

// 추가된 부분: 토큰 발행 함수
const mintTokens = async () => {
  // TON 토큰 컨트랙트에 연결
  const tokenContract = new ethers.Contract(TON, erc20ABI, l1Wallet);

  // 계정에 토큰을 발행
  const mintTx = await tokenContract.faucet(amountToTransfer);

  // mintTx의 결과를 확인
  const mintTxReceipt = await mintTx.wait();

  // mintTx가 성공적으로 완료되었는지 확인
  if (mintTxReceipt.status === 1) {
    console.log(`Minting successful. Transaction hash: ${mintTxReceipt.transactionHash}`);
  } else {
    console.error('Minting failed.');
  }
}

const updateAddresses = async (hre: HardhatRuntimeEnvironment) => {
  if (TON === '') {
    const Deployment__TON = await hre.deployments.get('TON')
    TON = Deployment__TON.address
  }

  if (addressManager === '') {
    const Deployment__AddressManager = await hre.deployments.get(
      'AddressManager'
    )
    addressManager = Deployment__AddressManager.address
  }

  if (l1CrossDomainMessenger === '') {
    const Deployment__L1CrossDomainMessenger = await hre.deployments.get(
      'L1CrossDomainMessengerProxy'
    )
    l1CrossDomainMessenger = Deployment__L1CrossDomainMessenger.address
  }

  if (l1StandardBridge === '') {
    const Deployment__L1StandardBridge = await hre.deployments.get(
      'L1StandardBridgeProxy'
    )
    l1StandardBridge = Deployment__L1StandardBridge.address
  }

  if (optimismPortal === '') {
    const Deployment__OptimismPortal = await hre.deployments.get(
      'OptimismPortalProxy'
    )
    optimismPortal = Deployment__OptimismPortal.address
  }

  if (l2OutputOracle === '') {
    const Deployment__L2OutputOracle = await hre.deployments.get(
      'L2OutputOracleProxy'
    )
    l2OutputOracle = Deployment__L2OutputOracle.address
  }
}

const L2StandardTokenFactory = '0x4200000000000000000000000000000000000012';
const bridgeABI =[{"type":"event","name":"StandardL2TokenCreated","inputs":[{"type":"address","name":"_l1Token","internalType":"address","indexed":true},{"type":"address","name":"_l2Token","internalType":"address","indexed":true}],"anonymous":false},{"type":"function","stateMutability":"nonpayable","outputs":[],"name":"createStandardL2Token","inputs":[{"type":"address","name":"_l1Token","internalType":"address"},{"type":"string","name":"_name","internalType":"string"},{"type":"string","name":"_symbol","internalType":"string"}]}]

const Bridgecontract = new ethers.Contract(L2StandardTokenFactory, bridgeABI, l2Provider);

// 과거에 발생한 이벤트를 조회합니다.
const eventFilter = Bridgecontract.filters.StandardL2TokenCreated();
Bridgecontract.queryFilter(eventFilter)
  .then(logs => {
    logs.forEach((log) => {
      if (log && log.args) {
        console.log(`L1 Token: ${log.args._l1Token}`);
        console.log(`L2 Token: ${log.args._l2Token}`);
      } else {
        console.log('Invalid event log:', log);
       }
     });
  })
  .catch(console.error);


const promise = Bridgecontract.queryFilter(eventFilter);

promise
  .then(logs => {
    console.log('Promise resolved with', logs.length, 'logs');
    // ...
  })
  .catch(error => {
    console.log('Promise rejected with error:', error);
  });

console.log('Promise status:', promise);





const depositTON = async (amount: NumberLike) => {
  console.log('Deposit TON:', amount)
  console.log('TON address:', TON)

  const tonContract = new ethers.Contract(TON, erc20ABI, l1Wallet)
  const WTON2 = '0x4200000000000000000000000000000000000006'
  const tonContractl2 = new ethers.Contract(TON2, erc20ABI, l2Wallet)


  const l1Contracts = {
    StateCommitmentChain: zeroAddr,
    CanonicalTransactionChain: zeroAddr,
    BondManager: zeroAddr,
    AddressManager: addressManager,
    L1CrossDomainMessenger: l1CrossDomainMessenger,
    L1StandardBridge: l1StandardBridge,
    OptimismPortal: optimismPortal,
    L2OutputOracle: l2OutputOracle,
  }
  console.log('l1 contracts:', l1Contracts)

  const bridges = {
    TON: {
      l1Bridge: l1Contracts.L1StandardBridge,
      l2Bridge: predeploys.L2StandardBridge,
      Adapter: TONBridgeAdapter,
    },
  }

  const l1ChainId = (await l1Provider.getNetwork()).chainId
  const l2ChainId = (await l2Provider.getNetwork()).chainId

  const messenger = new CrossChainMessenger({
    bedrock: true,
    contracts: {
      l1: l1Contracts,
    },
    bridges,
    l1ChainId,
    l2ChainId,
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  let l1TONBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance:', l1TONBalance.toString())

  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())

  const approveTx = await messenger.approveERC20(TON, ETH, amount)
  await approveTx.wait()
  console.log('approveTx:', approveTx.hash)

  const depositTx = await messenger.depositERC20(TON, ETH, amount)
  await depositTx.wait()
  console.log('depositTx:', depositTx.hash)

  await messenger.waitForMessageStatus(depositTx.hash, MessageStatus.RELAYED)

  l2Balance = await l2Wallet.getBalance()
  l1TONBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance: ', l1TONBalance.toString())
  console.log('l2 native balance: ', l2Balance.toString())
  const balance = await tonContractl2.balanceOf(l2Wallet.address);
  console.log(`L2 Ton balance: ${balance}`);

}

const withdrawTON = async (amount: NumberLike) => {
  console.log('Withdraw TON:', amount)

  const tonContract = new ethers.Contract(TON, erc20ABI, l1Wallet)

  const l1Contracts = {
    StateCommitmentChain: zeroAddr,
    CanonicalTransactionChain: zeroAddr,
    BondManager: zeroAddr,
    AddressManager: addressManager,
    L1CrossDomainMessenger: l1CrossDomainMessenger,
    L1StandardBridge: l1StandardBridge,
    OptimismPortal: optimismPortal,
    L2OutputOracle: l2OutputOracle,
  }

  const bridges = {
    TON: {
      l1Bridge: l1Contracts.L1StandardBridge,
      l2Bridge: predeploys.L2StandardBridge,
      Adapter: TONBridgeAdapter,
    },
  }

  const l1ChainId = (await l1Provider.getNetwork()).chainId
  const l2ChainId = (await l2Provider.getNetwork()).chainId

  const messenger = new CrossChainMessenger({
    bedrock: true,
    contracts: {
      l1: l1Contracts,
    },
    bridges,
    l1ChainId,
    l2ChainId,
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

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
  console.log('l2 native balance:', l2Balance.toString())

  // // Check ready for prove
  await messenger.waitForMessageStatus(
    withdrawalTx.transactionHash,
    MessageStatus.READY_TO_PROVE
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
      MessageStatus.READY_FOR_RELAY
    )
  } finally {
    clearInterval(finalizeInterval)
  }

  tonBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance: ', tonBalance.ㅎ())

  const tx = await messenger.finalizeMessage(withdrawalTx.transactionHash)
  const receipt = await tx.wait()
  console.log('Finalized message tx', receipt.transactionHash)
  console.log('Finalized withdrawal')

  tonBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance: ', tonBalance.toString())
}

  task('deposit-ton', 'Deposits ERC20-TON to L2.')
    .addParam('amount', 'Deposit amount', '1', types.string)
    .setAction(async (args, hre) => {
      await updateAddresses(hre)
      await depositTON(args.amount)
    })

  task('withdraw-ton', 'Withdraw native TON from L2.')
    .addParam('amount', 'Withdrawal amount', '1', types.string)
    .setAction(async (args, hre) => {
      await updateAddresses(hre)
      await withdrawTON(args.amount)
    })
