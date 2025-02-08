import { ethers, providers } from 'ethers'
import { task, types } from 'hardhat/config'

const normalize = (str: string) => {
  return str.startsWith('0x') ? str.slice(2) : str
}

task('verify-bytecode', 'verify l1 bytecodes')
  .addParam('standardbridge', 'L1 StandardBridge addresss', '0x', types.string)
  .addParam('xdm', 'L1 StandardBridge addresss', '0x', types.string)
  .addParam('optimismportal', 'L1 StandardBridge addresss', '0x', types.string)
  .addParam('systemconfig', 'L1 StandardBridge addresss', '0x', types.string)
  .addParam('rpc', 'L1 RPC', 'http://localhost:8545', types.string)
  .addParam(
    'expectedhash',
    'Expected hash',
    '0x6c36358a52dd5aa507cd00e6377b775495c6c6e8853ff9dc2d5d06c032f380b8',
    types.string
  )
  .addFlag('debug')
  .setAction(async (args) => {
    console.log('\n[Verifying bytecodes]\n')
    console.log('L1StandardBridge address:', args.standardbridge)
    console.log('L1CrossDomainMessenger address:', args.xdm)
    console.log('OptimismPortal address:', args.optimismportal)
    console.log('SystemConfig address:', args.systemconfig)
    console.log('L1 RPC:', args.rpc)

    const l1Provider = new providers.JsonRpcProvider(args.rpc)
    try {
      await l1Provider.getBlockNumber()
    } catch (error) {
      console.error('Cannot connect to RPC:', error.message)
      process.exit(1)
    }
    console.log('Connected to RPC!')

    const standardbridgeCode = normalize(
      await l1Provider.getCode(args.standardbridge)
    )
    const xdmCode = normalize(
      await l1Provider.getCode(args.xdm)
    )
    const optimismportalCode = normalize(
      await l1Provider.getCode(args.optimismportal)
    )
    const systemconfigCode = normalize(
      await l1Provider.getCode(args.systemconfig)
    )
    const bytecode = `0x${standardbridgeCode}${xdmCode}${optimismportalCode}${systemconfigCode}`

    const result = ethers.utils.keccak256(ethers.utils.arrayify(bytecode))
    if (args.debug) {
      let totalSize = 0
      console.log('\n\n----------------------------------------\n\n')
      console.log('[L1StandardBridge code]\n', standardbridgeCode)
      totalSize += standardbridgeCode.length
      console.log('\n\n----------------------------------------\n\n')
      console.log('[L1CrossDomainMessenger code]\n', xdmCode)
      totalSize += xdmCode.length
      console.log('\n\n----------------------------------------\n\n')
      console.log('[OptimismPortal code]\n', optimismportalCode)
      totalSize += optimismportalCode.length
      console.log('\n\n----------------------------------------\n\n')
      console.log('[SystemConfig code]\n', systemconfigCode)
      totalSize += systemconfigCode.length
      console.log('\n\n----------------------------------------\n\n')
      console.log('[Bytecode]\n', bytecode)
      console.log('\n\n\n')
      console.log('Total contract bytecodes size:', totalSize)
      console.log('Concat bytecode size:', bytecode.length)
      console.log('Real hash:', result)
      console.log('Expected hash:', args.expectedhash)
    }

    if (result === args.expectedhash) {
      console.log('Verified successfully')
      process.exit(0)
    }else {
      console.log('Verified failed')
      process.exit(1)
    }
  }
)
