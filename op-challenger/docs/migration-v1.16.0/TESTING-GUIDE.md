# E2E Testing Guide - Post v1.16.0 Migration

## Overview

This guide provides instructions for running E2E tests after the Optimism v1.16.0 migration, with a focus on the `TestOutputAlphabetGame_ReclaimBond` test and related fault proof tests.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Environment Setup](#environment-setup)
3. [Building Test Artifacts](#building-test-artifacts)
4. [Running Tests](#running-tests)
5. [Troubleshooting](#troubleshooting)
6. [Test Coverage](#test-coverage)

---

## Prerequisites

### Required Tools

```bash
# Verify Go version
go version  # Should be 1.21 or higher

# Verify Foundry installation
forge --version  # Should be 0.2.0 or higher

# Verify Python version
python3 --version  # Should be 3.8 or higher

# Verify Make
make --version
```

### Repository Structure

Ensure you have the following directory structure:

```
tokamak-thanos/
├── op-challenger/
├── op-e2e/
│   └── faultproofs/
│       ├── output_alphabet_test.go
│       ├── output_alphabet_bond_test.go
│       └── ... (other test files)
├── packages/
│   └── tokamak/
│       └── contracts-bedrock/
└── Makefile
```

---

## Environment Setup

### 1. Clean Previous Build Artifacts

Before starting, clean any previous build artifacts:

```bash
# From repository root
make clean

# Clean Solidity artifacts (if needed)
cd packages/tokamak/contracts-bedrock
forge clean
cd ../../..
```

### 2. Verify Go Modules

Ensure all Go modules are up to date:

```bash
# Download and verify modules
go mod download
go mod verify

# Tidy up module dependencies
go mod tidy
```

### 3. Check Environment Variables

No special environment variables are required for basic testing, but you can optionally set:

```bash
# Enable verbose test output
export VERBOSE=1

# Set test timeout (default: 10m)
export TEST_TIMEOUT=15m
```

---

## Building Test Artifacts

### Step 1: Generate Solidity Contract Artifacts

E2E tests require Solidity contract artifacts and state dumps:

```bash
# From repository root
make devnet-allocs
```

This command will:
1. Compile all Solidity contracts with multiple compiler versions (0.5.17, 0.6.12, 0.8.15, 0.8.19, 0.8.20, 0.8.25)
2. Deploy contracts to a simulated network
3. Configure 6 dispute game types (Alphabet, Fast, Cannon, PermissionedCannon, Asterisc, AsteriscKona)
4. Generate state dump files

**Expected Output**:
```
Solc 0.5.17 finished in 58.87ms
Solc 0.8.20 finished in 321.17ms
Solc 0.8.19 finished in 799.18ms
Solc 0.6.12 finished in 895.19ms
Solc 0.8.25 finished in 1.76s
Solc 0.8.15 finished in 19.60s
Compiler run successful!
...
set up op chain!
```

**Generated Files**:
- `packages/tokamak/contracts-bedrock/forge-artifacts/` - Contract artifacts (all compiled Solidity contracts)
- `packages/tokamak/contracts-bedrock/deployments/31337-deploy.json` - Deployment configuration
- `packages/tokamak/contracts-bedrock/deployments/devnetL1/` - DevnetL1 deployment artifacts

**Duration**: Approximately 30-60 seconds

### Step 2: Verify Artifact Generation

Check that key artifacts were created:

```bash
# Check for forge-artifacts directory
ls packages/tokamak/contracts-bedrock/forge-artifacts/ | head -20

# Check for deployment artifacts
ls packages/tokamak/contracts-bedrock/deployments/

# Verify dispute game factory artifact
ls packages/tokamak/contracts-bedrock/forge-artifacts/DisputeGameFactory.sol/

# Verify 31337-deploy.json was created
ls packages/tokamak/contracts-bedrock/deployments/31337-deploy.json
```

**Note**: The `make devnet-allocs` process generates state dump files internally during execution, but they are not persisted to disk. The E2E tests only require:
- `forge-artifacts/` directory with compiled contract artifacts
- `deployments/` directory with deployment configuration

### Step 3: Compile Test Package

Compile the E2E test package to verify no compilation errors:

```bash
# Compile without running tests
go build ./op-e2e/faultproofs
```

**Expected**: No output (success)

If there are errors, see [Troubleshooting](#troubleshooting).

---

## Running Tests

### Basic Test Execution

#### Run Specific Test

```bash
# Run TestOutputAlphabetGame_ReclaimBond test
go test -v \
  -run TestOutputAlphabetGame_ReclaimBond \
  -timeout 10m \
  ./op-e2e/faultproofs
```

**Options**:
- `-v`: Verbose output
- `-run`: Test name pattern (supports regex)
- `-timeout`: Maximum test duration

#### Run All Alphabet Game Tests

```bash
# Run all alphabet-related tests
go test -v \
  -run "TestOutputAlphabetGame.*" \
  -timeout 15m \
  ./op-e2e/faultproofs
```

#### Run All Fault Proof Tests

```bash
# Run entire fault proofs test suite
go test -v \
  -timeout 30m \
  ./op-e2e/faultproofs
```

### Advanced Test Options

#### Run with Race Detection

```bash
go test -v -race \
  -run TestOutputAlphabetGame_ReclaimBond \
  -timeout 15m \
  ./op-e2e/faultproofs
```

#### Run with Coverage

```bash
go test -v \
  -run TestOutputAlphabetGame_ReclaimBond \
  -coverprofile=coverage.out \
  -timeout 10m \
  ./op-e2e/faultproofs

# View coverage report
go tool cover -html=coverage.out
```

#### Run Specific Test with Debug Output

```bash
# Enable all debug logging
E2E_DEBUG=true go test -v \
  -run TestOutputAlphabetGame_ReclaimBond \
  -timeout 10m \
  ./op-e2e/faultproofs 2>&1 | tee test_output.log
```

### Background Test Execution

For long-running tests, you can run in the background:

```bash
# Run test in background, save output
nohup go test -v \
  -run TestOutputAlphabetGame_ReclaimBond \
  -timeout 10m \
  ./op-e2e/faultproofs > test_output.log 2>&1 &

# Check progress
tail -f test_output.log

# Find the process
ps aux | grep "go test"
```

---

## Test Structure

### TestOutputAlphabetGame_ReclaimBond

**Location**: `op-e2e/faultproofs/output_alphabet_bond_test.go`

**Purpose**: Tests bond reclaim functionality in alphabet dispute games after op-geth upgrade.

**Test Flow**:
1. Setup E2E environment with L1 and L2 chains
2. Deploy dispute game contracts
3. Create alphabet dispute game
4. Simulate game resolution
5. Test bond reclaim mechanism
6. Verify correct bond distribution

**Expected Duration**: 5-8 minutes

**Dependencies**:
- Solidity artifacts from `make devnet-allocs`
- op-geth v1.101601.0-rc.1
- Dispute game factory contracts
- Preimage oracle contracts

### Test Phases

1. **Initialization Phase** (30-60s):
   - Load contract artifacts
   - Initialize L1/L2 chains
   - Deploy test contracts

2. **Game Creation Phase** (60-90s):
   - Create alphabet dispute game
   - Submit initial claims
   - Establish game tree

3. **Resolution Phase** (120-180s):
   - Simulate claim challenges
   - Progress game to resolution
   - Verify game outcome

4. **Bond Reclaim Phase** (60-120s):
   - Test bond withdrawal
   - Verify bond distribution
   - Check balances

5. **Cleanup Phase** (10-30s):
   - Teardown test environment
   - Release resources

---

## Troubleshooting

### Common Issues

#### Issue 1: Missing Artifacts Error

**Error**:
```
failed to open artifact "DeployImplementations.s.sol/DeployImplementations.json":
no such file or directory
```

**Solution**:
```bash
# Generate artifacts
make devnet-allocs

# Verify artifact exists
ls packages/tokamak/contracts-bedrock/forge-artifacts/DeployImplementations.s.sol/
```

#### Issue 2: Compilation Errors

**Error**:
```
undefined: bindings.RAT
undefined: bindings.NewRAT
```

**Solution**:
These errors indicate RAT test files are present. They should have been backed up during migration:

```bash
# Verify RAT files are backed up
ls op-e2e/faultproofs/*.bak

# If not backed up, move them now
mv op-e2e/faultproofs/rat_e2e_test.go \
   op-e2e/faultproofs/rat_e2e_test.go.bak
mv op-e2e/faultproofs/rat_simple_test.go \
   op-e2e/faultproofs/rat_simple_test.go.bak
mv op-e2e/faultproofs/rat_unit_test.go \
   op-e2e/faultproofs/rat_unit_test.go.bak

# Retry compilation
go build ./op-e2e/faultproofs
```

#### Issue 3: Import Path Errors

**Error**:
```
package github.com/ethereum-optimism/optimism/... is not in GOROOT
```

**Solution**:
Import paths weren't updated correctly. Fix with:

```bash
# Find files with incorrect imports
grep -r "github.com/ethereum-optimism/optimism" \
  --include="*.go" \
  op-e2e/

# Fix import paths
find op-e2e/ -name "*.go" \
  -exec sed -i '' 's|github.com/ethereum-optimism/optimism|github.com/tokamak-network/tokamak-thanos|g' {} \;

# Verify fix
go build ./op-e2e/faultproofs
```

#### Issue 4: Test Timeout

**Error**:
```
panic: test timed out after 10m0s
```

**Solution**:
Increase timeout or optimize test:

```bash
# Increase timeout to 20 minutes
go test -v \
  -run TestOutputAlphabetGame_ReclaimBond \
  -timeout 20m \
  ./op-e2e/faultproofs
```

#### Issue 5: Port Already in Use

**Error**:
```
bind: address already in use
```

**Solution**:
```bash
# Find process using the port (usually 8545 or 9545)
lsof -i :8545
lsof -i :9545

# Kill the process
kill -9 <PID>

# Or kill all geth processes
pkill -9 geth
```

### Debug Mode

Enable debug mode for detailed logging:

```bash
# Set debug environment variables
export E2E_DEBUG=true
export LOG_LEVEL=debug

# Run test with full output
go test -v \
  -run TestOutputAlphabetGame_ReclaimBond \
  -timeout 15m \
  ./op-e2e/faultproofs 2>&1 | tee debug.log

# Search for specific errors
grep -i "error" debug.log
grep -i "panic" debug.log
```

### Clean Start

If tests continue to fail, perform a clean start:

```bash
# 1. Clean all build artifacts
make clean
cd packages/tokamak/contracts-bedrock && forge clean && cd ../../..

# 2. Clean Go cache
go clean -cache -testcache -modcache

# 3. Rebuild everything
make devnet-allocs

# 4. Verify compilation
go build ./op-e2e/faultproofs

# 5. Run test
go test -v -run TestOutputAlphabetGame_ReclaimBond ./op-e2e/faultproofs
```

---

## Test Coverage

### Current Test Status (Post-Migration)

| Test File | Status | Notes |
|-----------|--------|-------|
| `output_alphabet_test.go` | ✅ Compiles | Alphabet game tests |
| `output_alphabet_bond_test.go` | ✅ Compiles | Bond reclaim tests |
| `output_cannon_test.go` | ✅ Compiles | Cannon VM tests |
| `output_asterisc_test.go` | ✅ Compiles | Asterisc VM tests |
| `rat_e2e_test.go` | ⏸️ Excluded | Pending RAT bindings |
| `rat_simple_test.go` | ⏸️ Excluded | Pending RAT bindings |
| `rat_unit_test.go` | ⏸️ Excluded | Pending RAT bindings |

### Test Categories

#### Dispute Game Tests
- `TestOutputAlphabetGame_*` - Alphabet VM game tests
- `TestOutputCannonGame_*` - Cannon VM game tests
- `TestOutputAsteriscGame_*` - Asterisc VM game tests

#### Bond Tests
- `TestOutputAlphabetGame_ReclaimBond` - Bond reclaim verification
- `TestOutputAlphabetGame_BondDistribution` - Bond distribution modes

#### Fault Proof Tests
- `TestFaultProof_*` - General fault proof mechanisms
- `TestPreimage_*` - Preimage oracle tests
- `TestChallenge_*` - Challenge/response tests

---

## Performance Benchmarks

### Expected Test Durations (Post-Migration)

| Test | Duration | Resource Usage |
|------|----------|----------------|
| TestOutputAlphabetGame_ReclaimBond | 5-8 minutes | ~2GB RAM |
| Full alphabet suite | 15-20 minutes | ~3GB RAM |
| Full fault proofs suite | 30-45 minutes | ~4GB RAM |

### Optimization Tips

1. **Use Test Cache**:
   ```bash
   # Tests with identical code/data are cached
   go test -v ./op-e2e/faultproofs
   # Subsequent runs are faster
   ```

2. **Parallel Execution**:
   ```bash
   # Run tests in parallel (use with caution for E2E tests)
   go test -v -parallel 2 ./op-e2e/faultproofs
   ```

3. **Skip Slow Tests**:
   ```bash
   # Run only fast tests
   go test -v -short ./op-e2e/faultproofs
   ```

---

## CI/CD Integration

### GitHub Actions Example

```yaml
name: E2E Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    timeout-minutes: 60

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install Foundry
        uses: foundry-rs/foundry-toolchain@v1

      - name: Generate Artifacts
        run: make devnet-allocs

      - name: Run Bond Reclaim Test
        run: |
          go test -v \
            -run TestOutputAlphabetGame_ReclaimBond \
            -timeout 15m \
            ./op-e2e/faultproofs

      - name: Upload Logs
        if: failure()
        uses: actions/upload-artifact@v3
        with:
          name: test-logs
          path: |
            *.log
            packages/tokamak/contracts-bedrock/forge-artifacts/
```

---

## Additional Resources

### Related Documentation
- [Migration Guide](./MIGRATION-GUIDE.md) - Full migration process
- [Changelog](./CHANGELOG.md) - Detailed changes from v1.7.7 to v1.16.0
- [Optimism E2E Testing](https://github.com/ethereum-optimism/optimism/tree/develop/op-e2e) - Upstream testing docs

### Useful Commands

```bash
# List all tests
go test -list . ./op-e2e/faultproofs

# Run specific tests by pattern
go test -v -run ".*Bond.*" ./op-e2e/faultproofs

# Get test names only (no execution)
go test -v -run NonExistentTest ./op-e2e/faultproofs 2>&1 | grep "^=== RUN"

# Check test coverage
go test -cover ./op-e2e/faultproofs

# Profile test execution
go test -cpuprofile cpu.prof -memprofile mem.prof ./op-e2e/faultproofs
go tool pprof cpu.prof
```

### Getting Help

If you encounter issues not covered in this guide:

1. Check the [Changelog](./CHANGELOG.md) for known issues
2. Review the [Migration Guide](./MIGRATION-GUIDE.md) for setup verification
3. Search existing issues in the repository
4. Create a new issue with:
   - Test command used
   - Full error output
   - Environment details (`go version`, `forge --version`)
   - Steps to reproduce

---

## Summary

After v1.16.0 migration:

1. ✅ Generate artifacts: `make devnet-allocs`
2. ✅ Verify compilation: `go build ./op-e2e/faultproofs`
3. ✅ Run tests: `go test -v -run TestOutputAlphabetGame_ReclaimBond ./op-e2e/faultproofs`

The migration successfully enables E2E testing with the new op-geth v1.101601.0-rc.1 and Optimism v1.16.0 architecture.
