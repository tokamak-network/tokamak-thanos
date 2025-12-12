import fs from 'node:fs'
import path from 'node:path'

import { task } from 'hardhat/config'
import { ethers, BigNumber } from 'ethers'
import { red, green, blue, white, yellow } from 'console-log-colors'

import {
  L2Interface,
  ERC20,
  ForceWithdrawAssetEntry,
  ForceWithdrawClaim,
} from '../utils/types'

/**
 * force-withdraw:gen
 *
 * Scans L2 events and builds generate-assets3.json (Tech Spec v2.4).
 *
 * Block range priority: CLI (l2StartBlock/l2EndBlock) > ENV (L2_START_BLOCK/L2_END_BLOCK).
 * "latest" resolves via L2 RPC, and start/end must differ (end > start).
 * Deployment json is optional (used by trh-sdk shutdown); direct runs must work with ENV only.
 */

// Configuration
const CHUNK_SIZE_DEFAULT = 10000 // block chunk size
const MAX_RETRIES = 3 // max retry attempts
const RETRY_DELAY = 2000 // retry delay (ms)

// Progress bar utility
const progressBar = (current: number, total: number, label: string = '') => {
  const width = 30
  const percent = Math.floor((current / total) * 100)
  const filled = Math.floor((current / total) * width)
  const empty = width - filled
  const bar = '█'.repeat(filled) + '░'.repeat(empty)
  process.stdout.write(
    `\r  ${label} [${bar}] ${percent}% (${current}/${total})`
  )
  if (current >= total) {
    console.log('')
  }
}

// Sleep utility
const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

// Retry wrapper
const withRetry = async <T>(
  fn: () => Promise<T>,
  maxRetries: number = MAX_RETRIES,
  delayMs: number = RETRY_DELAY,
  label: string = ''
): Promise<T> => {
  let lastError: Error | null = null

  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    try {
      const result = await fn()
      return result
    } catch (err: any) {
      lastError = err
      if (attempt < maxRetries) {
        console.log(
          yellow(
            `  ⚠️ ${label} failed (attempt ${attempt}/${maxRetries}): ${err.message}`
          )
        )
        console.log(yellow(`     retrying after ${delayMs}ms...`))
        await sleep(delayMs * attempt) // Exponential backoff
      }
    }
  }

  throw lastError
}

// Chunked event query
const queryEventsChunked = async (
  contract: ethers.Contract,
  filter: ethers.EventFilter,
  startBlock: number,
  endBlock: number,
  chunkSize: number,
  label: string = 'Events'
): Promise<ethers.Event[]> => {
  const allEvents: ethers.Event[] = []
  const totalChunks = Math.ceil((endBlock - startBlock + 1) / chunkSize)

  console.log(`  ${label}: ${startBlock} → ${endBlock} (${totalChunks} chunks)`)

  for (let i = 0; i < totalChunks; i++) {
    const chunkStart = startBlock + i * chunkSize
    const chunkEnd = Math.min(chunkStart + chunkSize - 1, endBlock)

    progressBar(i + 1, totalChunks, label)

    const events = await withRetry(
      () => contract.queryFilter(filter, chunkStart, chunkEnd),
      MAX_RETRIES,
      RETRY_DELAY,
      `Chunk ${i + 1}`
    )

    allEvents.push(...events)
  }

  return allEvents
}

const runEnvChecker = () => {
  const required = ['CONTRACT_RPC_URL_L2']

  // Bridge addresses: prefer deployment JSON, ENV overrides when provided
  const optional = ['CONTRACTS_L2BRIDGE_ADDRESS', 'CONTRACTS_L1BRIDGE_ADDRESS']

  for (const key of required) {
    if (!process.env[key]) {
      console.log(red(`env Error: ${key} is not set.\n`))
      process.exit(1)
    }
  }

  for (const key of optional) {
    if (!process.env[key]) {
      console.log(yellow(`⚠️ Warning: ${key} is not set (optional override)\n`))
    }
  }
}

const loadDeploymentAddresses = (l1ChainId: number) => {
  try {
    const deploymentPath = path.join(
      process.cwd(),
      'packages',
      'tokamak',
      'contracts-bedrock',
      'deployments',
      `${l1ChainId}.json`
    )

    if (!fs.existsSync(deploymentPath)) {
      return {
        deploymentPath: undefined,
        L1_BRIDGE: undefined,
        L2_BRIDGE: undefined,
      }
    }

    const json = JSON.parse(fs.readFileSync(deploymentPath, 'utf-8')) as {
      L1StandardBridgeProxy?: string
      L2StandardBridgeProxy?: string
    }

    return {
      deploymentPath,
      L1_BRIDGE: json.L1StandardBridgeProxy,
      L2_BRIDGE: json.L2StandardBridgeProxy,
    }
  } catch {
    return {
      deploymentPath: undefined,
      L1_BRIDGE: undefined,
      L2_BRIDGE: undefined,
    }
  }
}

task(
  'force-withdraw:gen',
  'Generate force withdrawal assets snapshot (generate-assets3.json)'
)
  .addOptionalParam(
    'l2StartBlock',
    'L2 start block number (optional, falls back to ENV)'
  )
  .addOptionalParam(
    'l2EndBlock',
    'L2 end block number (optional, falls back to ENV)'
  )
  .addOptionalParam('output', 'Output file path', 'data/generate-assets3.json')
  .addOptionalParam(
    'chunkSize',
    'Block chunk size for RPC queries',
    String(CHUNK_SIZE_DEFAULT)
  )
  .addOptionalParam('maxRetries', 'Maximum retry attempts', String(MAX_RETRIES))
  .addFlag('skipVerify', 'Skip on-chain verification after generation (faster)')
  .setAction(async (args, hre) => {
    runEnvChecker()

    console.log(
      blue.bgBlue.bold('\n🚀 ForceWithdraw Asset Generation Started\n')
    )

    // Configuration from args
    const chunkSize = parseInt(args.chunkSize, 10)
    const maxRetries = parseInt(args.maxRetries, 10)

    // Provider setup
    const L2_RPC = process.env.CONTRACT_RPC_URL_L2!

    // Resolve L1 chainId from hre config or provider
    const l1ChainId =
      hre.network.config.chainId ??
      (await hre.ethers.provider.getNetwork()).chainId

    const dep = loadDeploymentAddresses(Number(l1ChainId))

    // ENV override > deployment json (deployment is optional; ENV-only must work)
    const L1_BRIDGE =
      process.env.CONTRACTS_L1BRIDGE_ADDRESS ||
      dep.L1_BRIDGE ||
      ethers.constants.AddressZero
    const L2_BRIDGE = process.env.CONTRACTS_L2BRIDGE_ADDRESS || dep.L2_BRIDGE

    if (!L2_BRIDGE) {
      console.log(
        red(
          'L2 bridge address is required (set CONTRACTS_L2BRIDGE_ADDRESS or provide deployment json).'
        )
      )
      process.exit(1)
    }

    const l2Provider = new ethers.providers.JsonRpcProvider(L2_RPC)

    // Block range (CLI > ENV L2_START_BLOCK/L2_END_BLOCK), required
    const l2StartRaw = args.l2StartBlock ?? process.env.L2_START_BLOCK
    const l2EndRaw = args.l2EndBlock ?? process.env.L2_END_BLOCK

    const resolveBlock = async (
      raw: string | undefined,
      label: string,
      provider: ethers.providers.JsonRpcProvider
    ): Promise<number> => {
      if (raw === undefined) {
        return NaN
      }
      if (raw === 'latest') {
        return provider.getBlockNumber()
      }
      const n = parseInt(raw, 10)
      if (Number.isNaN(n)) {
        console.log(red(`${label} must be a number or "latest".`))
        process.exit(1)
      }
      return n
    }

    const l2StartBlock = await resolveBlock(
      l2StartRaw,
      'l2StartBlock',
      l2Provider
    )
    const l2EndBlock = await resolveBlock(l2EndRaw, 'l2EndBlock', l2Provider)

    if (Number.isNaN(l2StartBlock) || Number.isNaN(l2EndBlock)) {
      console.log(
        red(
          'l2StartBlock and l2EndBlock are required (CLI or ENV: L2_START_BLOCK/L2_END_BLOCK).'
        )
      )
      process.exit(1)
    }

    if (l2EndBlock <= l2StartBlock) {
      console.log(
        red(
          'l2EndBlock must be greater than l2StartBlock (use +1 if you need a single-block scan).'
        )
      )
      process.exit(1)
    }

    console.log(`📊 Scanning L2 blocks: ${l2StartBlock} → ${l2EndBlock}`)
    console.log(`📍 L2 Bridge: ${L2_BRIDGE}`)
    console.log(`📍 L1 Bridge: ${L1_BRIDGE}`)
    if (dep.deploymentPath) {
      console.log(`📄 Deployment: ${dep.deploymentPath}`)
    } else {
      console.log(
        yellow('⚠️ Deployment json not found; using ENV addresses only.')
      )
    }
    if (L1_BRIDGE === ethers.constants.AddressZero) {
      console.log(
        yellow(
          '⚠️ L1 bridge address not set; proceeding without L1 context (allowed for gen).'
        )
      )
    }
    console.log(`⚙️  Chunk Size: ${chunkSize}, Max Retries: ${maxRetries}\n`)

    // L2 Bridge contract
    const l2BridgeContract = new ethers.Contract(
      L2_BRIDGE,
      L2Interface,
      l2Provider
    )

    // Step 1: Extract token pairs from DepositFinalized events (bounded, chunked)
    console.log(
      blue.bgGreen.bold(
        'Step 1: Extracting token pairs from DepositFinalized events...'
      )
    )

    const depositFilter = l2BridgeContract.filters.DepositFinalized()

    let depositEvents: ethers.Event[] = []
    try {
      depositEvents = await queryEventsChunked(
        l2BridgeContract,
        depositFilter,
        l2StartBlock,
        l2EndBlock,
        chunkSize,
        'DepositFinalized'
      )
    } catch (err: any) {
      console.log(
        red(
          `\nFailed to fetch events after ${maxRetries} retries: ${err.message}`
        )
      )
      console.log(red('Try reducing --chunk-size or check RPC connection'))
      process.exit(1)
    }

    console.log(
      green(`\n✅ Found ${depositEvents.length} DepositFinalized events\n`)
    )

    if (depositEvents.length === 0) {
      console.log(
        yellow('⚠️ No deposit events found. Creating empty output file.')
      )
      // Create empty output file
      const emptyOutputPath = args.output as string
      const emptyOutputDir = path.dirname(emptyOutputPath)
      if (!fs.existsSync(emptyOutputDir)) {
        fs.mkdirSync(emptyOutputDir, { recursive: true })
      }
      fs.writeFileSync(emptyOutputPath, JSON.stringify([], null, 2), 'utf-8')
      console.log(green(`\n✅ Empty output file created: ${emptyOutputPath}`))
      return
    }

    // Token pair map: l1Token -> l2Token
    const tokenPairs = new Map<string, string>()
    for (const event of depositEvents) {
      if (event.args) {
        const l1Token = event.args[0] as string
        const l2Token = event.args[1] as string
        if (!tokenPairs.has(l1Token)) {
          tokenPairs.set(l1Token, l2Token)
        }
      }
    }

    console.log(
      blue.bgGreen.bold(`Found ${tokenPairs.size} unique token pairs:`)
    )
    tokenPairs.forEach((l2, l1) => {
      console.log(`  L1: ${l1} → L2: ${l2}`)
    })
    console.log('')

    // Step 2: Collect holders and balances per token
    console.log(blue.bgGreen.bold('Step 2: Collecting holders and balances...'))

    const result: ForceWithdrawAssetEntry[] = []
    let tokenIndex = 0

    for (const [l1Token, l2Token] of tokenPairs) {
      tokenIndex++
      console.log(
        `\n⦿ [${tokenIndex}/${tokenPairs.size}] Processing token: ${l2Token}`
      )

      const l2TokenContract = new ethers.Contract(l2Token, ERC20, l2Provider)

      let tokenName = 'Unknown'
      try {
        tokenName = await withRetry(
          () => l2TokenContract.name(),
          maxRetries,
          RETRY_DELAY,
          'Token name'
        )
      } catch {
        tokenName =
          l1Token === ethers.constants.AddressZero ? 'Ether' : 'Unknown'
      }

      console.log(`  Token Name: ${tokenName}`)

      // Collect holder addresses via Transfer events (bounded, chunked)
      const transferFilter = l2TokenContract.filters.Transfer()
      let transferEvents: ethers.Event[] = []

      try {
        transferEvents = await queryEventsChunked(
          l2TokenContract,
          transferFilter,
          l2StartBlock,
          l2EndBlock,
          chunkSize,
          'Transfer'
        )
      } catch (err: any) {
        console.log(red(`\n  Failed to fetch Transfer events: ${err.message}`))
        continue
      }

      // Collect unique addresses
      const addresses = new Set<string>()
      for (const event of transferEvents) {
        if (event.args) {
          addresses.add(event.args[0] as string) // from
          addresses.add(event.args[1] as string) // to
        }
      }

      console.log(`  Found ${addresses.size} unique addresses`)

      // Fetch balances and build claims
      const claims: ForceWithdrawClaim[] = []
      let totalBalance = BigNumber.from(0)
      let processedCount = 0
      const totalAddresses = addresses.size

      for (const address of addresses) {
        processedCount++
        if (processedCount % 100 === 0 || processedCount === totalAddresses) {
          progressBar(processedCount, totalAddresses, 'Balances')
        }

        try {
          const balance: BigNumber = await withRetry(
            () => l2TokenContract.balanceOf(address),
            maxRetries,
            RETRY_DELAY,
            `Balance ${address.slice(0, 10)}`
          )

          if (balance.gt(0)) {
            // hash: keccak256(abi.encodePacked(l1Token, claimer, amount))
            const hash = ethers.utils.solidityKeccak256(
              ['address', 'address', 'uint256'],
              [l1Token, address, balance]
            )

            claims.push({
              claimer: address,
              amount: balance.toString(),
              hash,
            })

            totalBalance = totalBalance.add(balance)
          }
        } catch (err: any) {
          // Warn on individual balance fetch failures; continue
          console.log(
            yellow(
              `\n  ⚠️ Failed to get balance for ${address.slice(0, 10)}...: ${
                err.message
              }`
            )
          )
        }
      }

      console.log(`  Holders with balance > 0: ${claims.length}`)
      console.log(
        `  Total balance: ${ethers.utils.formatEther(
          totalBalance
        )} (wei: ${totalBalance.toString()})`
      )

      result.push({
        l1Token,
        l2Token,
        tokenName,
        data: claims,
      })
    }

    // Step 3: Write JSON
    console.log(blue.bgGreen.bold('\n\nStep 3: Saving to JSON file...'))

    const outputPath = args.output as string
    const outputDir = path.dirname(outputPath)

    // Ensure directory exists
    if (!fs.existsSync(outputDir)) {
      fs.mkdirSync(outputDir, { recursive: true })
    }

    fs.writeFileSync(outputPath, JSON.stringify(result, null, 2), 'utf-8')

    console.log(green.bold(`\n✅ Successfully generated: ${outputPath}`))
    console.log(`   Total tokens: ${result.length}`)
    console.log(
      `   Total claims: ${result.reduce(
        (sum, entry) => sum + entry.data.length,
        0
      )}`
    )

    // Summary output
    console.log(white.bgGreen.bold('\n📋 Summary:'))
    for (const entry of result) {
      const totalAmount = entry.data.reduce(
        (sum, claim) => sum.add(BigNumber.from(claim.amount)),
        BigNumber.from(0)
      )
      console.log(
        `  ${entry.tokenName}: ${
          entry.data.length
        } holders, ${ethers.utils.formatEther(totalAmount)} total`
      )
    }

    // Step 4: On-chain verification (unless skipped)
    const skipVerify = args.skipVerify as boolean
    if (skipVerify) {
      console.log(
        yellow('\n⚠️ On-chain verification skipped (--skip-verify flag)')
      )
      console.log(
        yellow(
          'Warning: Generated data has NOT been verified against blockchain state'
        )
      )
      return
    }

    console.log(
      blue.bgGreen.bold(
        '\n\nStep 4: Verifying generated data against on-chain state...'
      )
    )

    let totalVerified = 0
    let hashErrors = 0
    let balanceErrors = 0
    const totalSupplyWarnings = 0
    let totalSupplyErrors = 0

    for (let verifyIndex = 0; verifyIndex < result.length; verifyIndex++) {
      const entry = result[verifyIndex]
      const { l1Token, l2Token, tokenName, data: claims } = entry

      console.log(
        `\n⦿ [${verifyIndex + 1}/${result.length}] Verifying: ${tokenName}`
      )

      const l2TokenContract = new ethers.Contract(l2Token, ERC20, l2Provider)

      let processedCount = 0
      const totalClaimsForToken = claims.length

      for (const claim of claims) {
        processedCount++
        totalVerified++

        if (
          processedCount % 10 === 0 ||
          processedCount === totalClaimsForToken
        ) {
          progressBar(processedCount, totalClaimsForToken, 'Verifying')
        }

        const { claimer, amount, hash } = claim

        // Hash verification
        const calculatedHash = ethers.utils.solidityKeccak256(
          ['address', 'address', 'uint256'],
          [l1Token, claimer, amount]
        )

        if (calculatedHash !== hash) {
          hashErrors++
          console.log(
            red(
              `\n  ❌ Hash mismatch for ${claimer}: expected ${hash}, got ${calculatedHash}`
            )
          )
        }

        // Balance verification
        try {
          const actualBalance: BigNumber = await withRetry(
            () => l2TokenContract.balanceOf(claimer),
            maxRetries,
            RETRY_DELAY,
            `Balance ${claimer.slice(0, 10)}`
          )

          if (actualBalance.toString() !== amount) {
            balanceErrors++
            console.log(
              red(
                `\n  ❌ Balance mismatch for ${claimer}: file=${amount}, chain=${actualBalance.toString()}`
              )
            )
          }
        } catch (err: any) {
          balanceErrors++
          console.log(
            red(
              `\n  ❌ Failed to verify balance for ${claimer}: ${err.message}`
            )
          )
        }
      }

      // Total supply check
      try {
        const totalSupply: BigNumber = await withRetry(
          () => l2TokenContract.totalSupply(),
          maxRetries,
          RETRY_DELAY,
          'Total supply'
        )

        const sumOfBalances = claims.reduce(
          (sum, claim) => sum.add(BigNumber.from(claim.amount)),
          BigNumber.from(0)
        )

        // Handle zero total supply case
        if (totalSupply.isZero()) {
          if (sumOfBalances.isZero()) {
            console.log(
              green(
                `  ✅ Supply check: Total supply is 0, sum of claims is also 0 (valid empty token)`
              )
            )
          } else {
            totalSupplyErrors++
            console.log(
              red(
                `\n  ❌ CRITICAL: Total supply is 0 but sum of claims is ${sumOfBalances.toString()}!`
              )
            )
            console.log(
              red(
                '     This indicates a severe data integrity issue - claims exist for a token with no supply.'
              )
            )
          }
        } else if (sumOfBalances.gt(totalSupply)) {
          // Sum exceeds total supply - critical error
          totalSupplyErrors++
          console.log(
            red(
              `\n  ❌ CRITICAL: Sum of claims (${sumOfBalances.toString()}) exceeds total supply (${totalSupply.toString()})`
            )
          )
          console.log(
            red('     This is impossible and indicates corrupted data.')
          )
        } else if (sumOfBalances.eq(totalSupply)) {
          // Perfect match
          console.log(
            green(
              `  ✅ Supply check: Sum of claims equals total supply (${ethers.utils.formatEther(
                totalSupply
              )}) - perfect match`
            )
          )
        } else {
          // Sum is less than total supply - normal case
          const diff = totalSupply.sub(sumOfBalances)
          const percentDiff = diff.mul(10000).div(totalSupply).toNumber() / 100
          console.log(
            green(
              `  ✅ Supply check: ${ethers.utils.formatEther(
                diff
              )} (${percentDiff.toFixed(
                2
              )}%) unaccounted (expected for contract-held tokens)`
            )
          )
        }
      } catch (err: any) {
        totalSupplyErrors++
        console.log(red(`\n  ❌ Failed to get total supply: ${err.message}`))
        console.log(
          red(
            '     Cannot verify token integrity without total supply information.'
          )
        )
      }
    }

    // Verification summary
    console.log(white.bgBlue.bold('\n\n🔍 Verification Summary:\n'))
    console.log(`  Total claims verified: ${totalVerified}`)
    console.log(`  Hash errors: ${hashErrors}`)
    console.log(`  Balance errors: ${balanceErrors}`)
    console.log(`  Total supply errors: ${totalSupplyErrors}`)
    console.log(`  Total supply warnings: ${totalSupplyWarnings}`)

    const hasErrors =
      hashErrors > 0 || balanceErrors > 0 || totalSupplyErrors > 0

    if (!hasErrors) {
      console.log(green.bold('\n✅ On-chain verification PASSED!'))
      console.log(
        green('All generated data is consistent with blockchain state.')
      )
      if (totalSupplyWarnings > 0) {
        console.log(
          yellow(
            `\n⚠️ Note: ${totalSupplyWarnings} tokens had non-critical supply warnings (this may be normal).`
          )
        )
      }
    } else {
      console.log(red.bold('\n❌ On-chain verification FAILED!'))

      const errorDetails: string[] = []
      if (hashErrors > 0) {
        errorDetails.push(`${hashErrors} hash mismatches`)
      }
      if (balanceErrors > 0) {
        errorDetails.push(`${balanceErrors} balance mismatches`)
      }
      if (totalSupplyErrors > 0) {
        errorDetails.push(`${totalSupplyErrors} total supply errors`)
      }

      console.log(red(`  Found: ${errorDetails.join(', ')}`))
      console.log(
        red(
          '\n  The generated file contains errors and MUST NOT be used for force withdrawals.'
        )
      )
      console.log(red('  Please investigate the errors above and regenerate.'))
      process.exit(1)
    }
  })
