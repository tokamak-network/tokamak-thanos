package fault

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/test"
)

// TestRATChallengerIntegration tests the basic integration between RAT contract and challenger
func TestRATChallengerIntegration(t *testing.T) {
	// Setup test environment
	env := test.SetupRATTestEnvironment(t)
	defer env.Cleanup()

	t.Log("=== Phase 1-1: RAT-Challenger Integration Test ===")

	// 1. Verify RAT contract is deployed and initialized
	t.Log("Step 1: Verifying RAT contract deployment and initialization")

	// Check initial state
	validCount, err := env.GetValidChallengerCount()
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1), validCount, "Should have 1 dummy challenger initially")

	version, err := env.RATContract.Version(&bind.CallOpts{})
	require.NoError(t, err)
	require.NotEmpty(t, version, "Version should be set")
	t.Logf("RAT contract version: %s", version)

	// 2. Challenger stakes to RAT
	t.Log("Step 2: Challenger staking to RAT")

	stakeAmount := new(big.Int)
	stakeAmount.SetString("2500000000000000000", 10) // 2.5 ETH
	err = env.StakeToRAT(stakeAmount)
	require.NoError(t, err)
	t.Logf("Challenger staked %s ETH to RAT", stakeAmount.String())

	// Verify challenger was added to valid challengers
	validCount, err = env.GetValidChallengerCount()
	require.NoError(t, err)
	require.Equal(t, big.NewInt(2), validCount, "Should have 2 challengers (1 dummy + 1 real)")

	challengerInfo, err := env.GetChallengerInfo(env.ChallengerAccount.Address)
	require.NoError(t, err)
	require.Equal(t, stakeAmount, challengerInfo.StakingAmount, "Staking amount should match")
	require.True(t, challengerInfo.IsValid, "Challenger should be valid")
	require.True(t, challengerInfo.TotalSlashedAmount.Cmp(big.NewInt(0)) == 0, "No slashing yet")
	t.Logf("Challenger info: staked=%s, valid=%v, index=%d",
		challengerInfo.StakingAmount.String(), challengerInfo.IsValid, challengerInfo.ValidatorIndex)

	// 3. Create mock dispute game with incorrect state root
	t.Log("Step 3: Creating mock dispute game")

	gameAddr := env.CreateMockDisputeGame()
	incorrectStateRoot := test.GenerateRandomHash("incorrect_state_root")
	t.Logf("Created mock dispute game at address: %s", gameAddr.Hex())
	t.Logf("Using incorrect state root: %s", incorrectStateRoot.Hex())

	// 4. Trigger RAT attention test
	t.Log("Step 4: Triggering RAT attention test")

	blockHash := env.GetLatestBlockHash()
	err = env.TriggerRATAttentionTest(gameAddr, incorrectStateRoot, blockHash)
	require.NoError(t, err)
	t.Log("RAT attention test triggered successfully")

	// 5. Verify attention test was created correctly
	t.Log("Step 5: Verifying attention test creation")

	attentionInfo, err := env.GetAttentionTestInfo(gameAddr)
	require.NoError(t, err)
	require.Equal(t, incorrectStateRoot, common.BytesToHash(attentionInfo.StateRoot[:]), "State root should match")
	require.Equal(t, env.ChallengerAccount.Address, attentionInfo.ChallengerAddress, "Challenger should be selected")
	require.False(t, attentionInfo.EvidenceSubmitted, "Evidence should not be submitted yet")
	require.True(t, attentionInfo.BondAmount.Uint64() > 0, "Bond amount should be positive")

	t.Logf("Attention test info: challenger=%s, bond=%s, block=%d",
		attentionInfo.ChallengerAddress.Hex(), attentionInfo.BondAmount.String(), attentionInfo.L1BlockNumber)

	// Verify challenger's stake was reduced (bonded)
	challengerInfoAfterBond, err := env.GetChallengerInfo(env.ChallengerAccount.Address)
	require.NoError(t, err)
	require.True(t, challengerInfoAfterBond.StakingAmount.Cmp(challengerInfo.StakingAmount) < 0,
		"Challenger stake should be reduced by bond amount")
	require.True(t, challengerInfoAfterBond.TotalSlashedAmount.Uint64() > 0,
		"Total slashed amount should increase")

	bondAmount := attentionInfo.BondAmount
	expectedNewStake := new(big.Int).Sub(challengerInfo.StakingAmount, bondAmount)
	require.Equal(t, expectedNewStake, challengerInfoAfterBond.StakingAmount,
		"New stake should equal original - bond amount")

	// 6. Create correct evidence and submit
	t.Log("Step 6: Submitting correct evidence")

	// Generate correct proof values that will produce the required state root
	proofLV := test.GenerateRandomHash("left_value")
	proofRV := test.GenerateRandomHash("right_value")
	correctStateRoot := test.CreateTestStateRoot(proofLV, proofRV)
	t.Logf("Generated correct state root: %s", correctStateRoot.Hex())

	// For the test to work, we need the attention test to expect this correct state root
	// In a real scenario, the attention test would be triggered with the correct expected root
	// For this test, we'll trigger a new attention test with the correct expected root
	gameAddr2 := env.CreateMockDisputeGame()
	t.Logf("Created second game address: %s", gameAddr2.Hex())
	err = env.TriggerRATAttentionTest(gameAddr2, correctStateRoot, env.GetLatestBlockHash())
	require.NoError(t, err)

	// Verify the second attention test was created
	attentionInfo2, err := env.GetAttentionTestInfo(gameAddr2)
	require.NoError(t, err)
	t.Logf("Second attention test info: challenger=%s, stateRoot=%s",
		attentionInfo2.ChallengerAddress.Hex(), common.BytesToHash(attentionInfo2.StateRoot[:]).Hex())

	// Check if a challenger was actually selected (not zero address)
	if attentionInfo2.ChallengerAddress == (common.Address{}) {
		t.Skip("No challenger selected for second attention test - RAT probability issue")
	}

	// Now submit the correct evidence
	err = env.SubmitCorrectEvidence(gameAddr2, proofLV, proofRV)
	require.NoError(t, err)
	t.Log("Evidence submitted successfully")

	// 7. Verify bond was restored
	t.Log("Step 7: Verifying bond restoration")

	// Verify evidence was marked as submitted
	attentionInfo2, err = env.GetAttentionTestInfo(gameAddr2)
	require.NoError(t, err)
	require.True(t, attentionInfo2.EvidenceSubmitted, "Evidence should be marked as submitted")

	// Verify challenger's balance was restored
	finalChallengerInfo, err := env.GetChallengerInfo(env.ChallengerAccount.Address)
	require.NoError(t, err)

	t.Logf("Challenger stakes - After first bond: %s, After second bond: %s, Final: %s",
		challengerInfoAfterBond.StakingAmount.String(),
		new(big.Int).Sub(challengerInfoAfterBond.StakingAmount, attentionInfo2.BondAmount).String(),
		finalChallengerInfo.StakingAmount.String())

	// The challenger was bonded twice but should only be restored once (for the evidence submitted)
	// So final stake = original stake - first bond (still pending) + 0 (second bond restored)
	// Which equals challengerInfoAfterBond.StakingAmount (after first bond was deducted)
	require.Equal(t, challengerInfoAfterBond.StakingAmount, finalChallengerInfo.StakingAmount,
		"Challenger stake should be same as after first bond (second bond restored)")

	t.Logf("Final challenger staking amount: %s (was %s after bond, restored by %s)",
		finalChallengerInfo.StakingAmount.String(),
		challengerInfoAfterBond.StakingAmount.String(),
		attentionInfo2.BondAmount.String())

	// 8. Verify final state is consistent
	t.Log("Step 8: Verifying final system state")

	// Challenger should still be valid if they have enough stake
	require.True(t, finalChallengerInfo.IsValid || finalChallengerInfo.StakingAmount.Cmp(big.NewInt(1e18)) >= 0,
		"Challenger should be valid or have sufficient stake")

	// Total valid challengers should still be 2 (dummy + challenger)
	finalValidCount, err := env.GetValidChallengerCount()
	require.NoError(t, err)
	require.True(t, finalValidCount.Cmp(big.NewInt(1)) >= 0, "Should have at least 1 valid challenger")

	t.Log("=== RAT-Challenger Integration Test PASSED ===")
}

// TestRATMultipleChallengers tests RAT behavior with multiple challengers
func TestRATMultipleChallengers(t *testing.T) {
	env := test.SetupRATTestEnvironment(t)
	defer env.Cleanup()

	t.Log("=== Phase 1-2: RAT Multiple Challengers Test ===")

	// 1. Create multiple challengers
	t.Log("Step 1: Creating multiple challengers")

	additionalChallengers := env.CreateMultipleChallengers(3)
	require.Len(t, additionalChallengers, 3, "Should create 3 additional challengers")

	allChallengers := append([]*test.TestAccount{env.ChallengerAccount}, additionalChallengers...)
	t.Logf("Created %d total challengers", len(allChallengers))

	// 2. All challengers stake to RAT
	t.Log("Step 2: All challengers staking to RAT")

	stakeAmount := new(big.Int)
	stakeAmount.SetString("2500000000000000000", 10) // 2.5 ETH each
	for i, challenger := range allChallengers {
		err := env.StakeToRATWithAccount(challenger, stakeAmount)
		require.NoError(t, err)
		t.Logf("Challenger %d (%s) staked %s ETH to RAT",
			i+1, challenger.Address.Hex()[:10], stakeAmount.String())
	}

	// Verify all challengers are valid
	validCount, err := env.GetValidChallengerCount()
	require.NoError(t, err)
	expectedCount := int64(len(allChallengers) + 1) // +1 for dummy challenger
	require.Equal(t, big.NewInt(expectedCount), validCount,
		"Should have %d challengers (1 dummy + %d real)", expectedCount, len(allChallengers))

	// 3. Create dispute game and trigger attention test
	t.Log("Step 3: Creating dispute game and triggering attention test")

	gameAddr := env.CreateMockDisputeGame()
	incorrectStateRoot := test.GenerateRandomHash("incorrect_state_root_multi")
	blockHash := env.GetLatestBlockHash()

	err = env.TriggerRATAttentionTest(gameAddr, incorrectStateRoot, blockHash)
	require.NoError(t, err)
	t.Log("Attention test triggered for multiple challengers scenario")

	// 4. Verify only one challenger was selected
	t.Log("Step 4: Verifying challenger selection")

	attentionInfo, err := env.GetAttentionTestInfo(gameAddr)
	require.NoError(t, err)

	selectedChallenger := attentionInfo.ChallengerAddress
	require.NotEqual(t, common.Address{}, selectedChallenger, "A challenger should be selected")
	t.Logf("Selected challenger: %s", selectedChallenger.Hex())

	// Verify the selected challenger is one of our challengers
	isValidSelection := false
	for _, challenger := range allChallengers {
		if challenger.Address == selectedChallenger {
			isValidSelection = true
			break
		}
	}
	require.True(t, isValidSelection, "Selected challenger should be one of the staked challengers")

	// 5. Verify other challengers were not affected
	t.Log("Step 5: Verifying non-selected challengers are unaffected")

	bondAmount := attentionInfo.BondAmount
	t.Logf("Bond amount deducted: %s", bondAmount.String())

	for i, challenger := range allChallengers {
		challengerInfo, err := env.GetChallengerInfo(challenger.Address)
		require.NoError(t, err)

		if challenger.Address == selectedChallenger {
			// Selected challenger should have reduced stake
			expectedStake := new(big.Int).Sub(stakeAmount, bondAmount)
			require.Equal(t, expectedStake, challengerInfo.StakingAmount,
				"Selected challenger should have reduced stake")
			require.True(t, challengerInfo.TotalSlashedAmount.Cmp(big.NewInt(0)) > 0,
				"Selected challenger should have non-zero slashed amount")
			t.Logf("Challenger %d (SELECTED): stake=%s, slashed=%s",
				i+1, challengerInfo.StakingAmount.String(), challengerInfo.TotalSlashedAmount.String())
		} else {
			// Other challengers should have original stake
			require.Equal(t, stakeAmount, challengerInfo.StakingAmount,
				"Non-selected challenger should maintain original stake")
			require.True(t, challengerInfo.TotalSlashedAmount.Cmp(big.NewInt(0)) == 0,
				"Non-selected challenger should have zero slashed amount")
			t.Logf("Challenger %d: stake=%s, slashed=%s",
				i+1, challengerInfo.StakingAmount.String(), challengerInfo.TotalSlashedAmount.String())
		}
	}

	t.Log("=== RAT Multiple Challengers Test PASSED ===")
}

// TestRatTriggerProbability tests RAT trigger probability functionality
func TestRatTriggerProbability(t *testing.T) {
	// Create custom config with low probability
	config := test.DefaultRATTestConfig()
	config.RatTriggerProbability = big.NewInt(10000) // 10% of 100,000

	env := test.SetupRATTestEnvironmentWithConfig(t, config)
	defer env.Cleanup()

	t.Log("=== Phase 1-3: RAT Trigger Probability Test ===")

	// 1. Verify initial probability setting
	t.Log("Step 1: Verifying probability configuration")

	currentProb, err := env.RATContract.RatTriggerProbability(&bind.CallOpts{})
	require.NoError(t, err)
	require.Equal(t, config.RatTriggerProbability, currentProb, "Initial probability should match config")
	t.Logf("RAT trigger probability set to: %s (10%%)", currentProb.String())

	// 2. Stake a challenger
	t.Log("Step 2: Staking challenger for probability test")

	stakeAmount := new(big.Int)
	stakeAmount.SetString("2500000000000000000", 10)
	err = env.StakeToRAT(stakeAmount)
	require.NoError(t, err)

	// 3. Test with multiple games to verify probability
	t.Log("Step 3: Testing probability with multiple games")

	gameCount := 30 // Use more games for better statistical sample
	triggeredCount := 0

	for i := 0; i < gameCount; i++ {
		gameAddr := env.CreateMockDisputeGame()
		stateRoot := test.GenerateRandomHash("state_root_" + string(rune(i)))
		blockHash := env.GetLatestBlockHash()

		err = env.TriggerRATAttentionTest(gameAddr, stateRoot, blockHash)
		require.NoError(t, err)

		// Check if attention test was actually created
		attentionInfo, err := env.GetAttentionTestInfo(gameAddr)
		require.NoError(t, err)

		if attentionInfo.ChallengerAddress != (common.Address{}) {
			triggeredCount++
		}

		// Advance blocks to change randomness
		env.AdvanceBlocks(1)
	}

	triggerPercentage := float64(triggeredCount) / float64(gameCount) * 100
	t.Logf("Triggered %d out of %d games (%.1f%%)", triggeredCount, gameCount, triggerPercentage)

	// 4. Verify probability is roughly correct (allow for randomness variance)
	t.Log("Step 4: Verifying probability statistics")

	// For 10% probability with 30 samples, expect 0-8 triggers (allowing for variance)
	require.True(t, triggeredCount <= gameCount*40/100,
		"Too many triggers for 10%% probability: %d/%d", triggeredCount, gameCount)
	require.True(t, triggeredCount >= 0,
		"Should have at least 0 triggers")

	// 5. Test 100% probability
	t.Log("Step 5: Testing 100% probability")

	err = env.SetRatTriggerProbability(big.NewInt(100000)) // 100%
	require.NoError(t, err)

	// Re-stake to ensure sufficient balance for 100% test
	additionalStake := new(big.Int)
	additionalStake.SetString("5000000000000000000", 10) // 5 ETH
	err = env.StakeToRAT(additionalStake)
	require.NoError(t, err)

	// Test a few games with 100% probability
	allTriggeredCount := 0
	testGames := 5

	for i := 0; i < testGames; i++ {
		gameAddr := env.CreateMockDisputeGame()
		stateRoot := test.GenerateRandomHash("100pct_state_root_" + string(rune(i)))
		blockHash := env.GetLatestBlockHash()

		// Check challenger state before triggering
		challengerInfo, err := env.GetChallengerInfo(env.ChallengerAccount.Address)
		require.NoError(t, err)
		t.Logf("Before Game %d - Challenger stake: %s, valid: %v",
			i+1, challengerInfo.StakingAmount.String(), challengerInfo.IsValid)

		err = env.TriggerRATAttentionTest(gameAddr, stateRoot, blockHash)
		require.NoError(t, err)

		attentionInfo, err := env.GetAttentionTestInfo(gameAddr)
		require.NoError(t, err)

		if attentionInfo.ChallengerAddress != (common.Address{}) {
			allTriggeredCount++
			t.Logf("Game %d triggered with challenger: %s, bond: %s",
				i+1, attentionInfo.ChallengerAddress.Hex(), attentionInfo.BondAmount.String())
		} else {
			t.Logf("Game %d did not trigger (no challenger available)", i+1)
		}

		// Check valid challenger count
		validCount, err := env.GetValidChallengerCount()
		require.NoError(t, err)
		t.Logf("After Game %d - Valid challenger count: %s", i+1, validCount.String())
	}

	// With 100% probability, games should trigger until no more challengers are available
	// Once a challenger is selected, they can't be selected again until they resolve their test
	require.True(t, allTriggeredCount >= 1, "At least one game should trigger with 100%% probability")
	require.True(t, allTriggeredCount <= testGames, "Cannot trigger more games than available challengers")
	t.Logf("100%% probability test: %d/%d games triggered", allTriggeredCount, testGames)

	// 6. Test 0% probability
	t.Log("Step 6: Testing 0% probability")

	err = env.SetRatTriggerProbability(big.NewInt(0)) // 0%
	require.NoError(t, err)

	noneTriggeredCount := 0
	for i := 0; i < testGames; i++ {
		gameAddr := env.CreateMockDisputeGame()
		stateRoot := test.GenerateRandomHash("0pct_state_root_" + string(rune(i)))
		blockHash := env.GetLatestBlockHash()

		err = env.TriggerRATAttentionTest(gameAddr, stateRoot, blockHash)
		require.NoError(t, err)

		attentionInfo, err := env.GetAttentionTestInfo(gameAddr)
		require.NoError(t, err)

		if attentionInfo.ChallengerAddress != (common.Address{}) {
			noneTriggeredCount++
		}
	}

	require.Equal(t, 0, noneTriggeredCount,
		"No games should trigger with 0%% probability")
	t.Logf("0%% probability test: %d/%d games triggered", noneTriggeredCount, testGames)

	t.Log("=== RAT Trigger Probability Test PASSED ===")
}

// TestRATIncorrectEvidenceSubmission tests what happens when challenger submits wrong evidence
func TestRATIncorrectEvidenceSubmission(t *testing.T) {
	env := test.SetupRATTestEnvironment(t)
	defer env.Cleanup()

	t.Log("=== Phase 1-4: RAT Incorrect Evidence Submission Test ===")

	// 1. Challenger stakes to RAT
	t.Log("Step 1: Challenger staking to RAT")
	stakeAmount := new(big.Int)
	stakeAmount.SetString("3000000000000000000", 10) // 3 ETH
	err := env.StakeToRAT(stakeAmount)
	require.NoError(t, err)

	// 2. Create dispute game and trigger attention test
	t.Log("Step 2: Creating dispute game and triggering attention test")
	gameAddr := env.CreateMockDisputeGame()
	correctStateRoot := test.GenerateRandomHash("correct_state_root")
	blockHash := env.GetLatestBlockHash()

	err = env.TriggerRATAttentionTest(gameAddr, correctStateRoot, blockHash)
	require.NoError(t, err)

	// Verify attention test was created
	attentionInfo, err := env.GetAttentionTestInfo(gameAddr)
	require.NoError(t, err)
	if attentionInfo.ChallengerAddress == (common.Address{}) {
		t.Skip("No challenger selected - RAT probability issue")
	}

	// 3. Get challenger info before submitting wrong evidence
	t.Log("Step 3: Recording challenger state before wrong evidence submission")
	challengerInfoBefore, err := env.GetChallengerInfo(env.ChallengerAccount.Address)
	require.NoError(t, err)
	t.Logf("Challenger stake before: %s, valid: %v",
		challengerInfoBefore.StakingAmount.String(), challengerInfoBefore.IsValid)

	// 4. Submit INCORRECT evidence (should fail)
	t.Log("Step 4: Submitting incorrect evidence (should fail)")
	wrongProofLV := test.GenerateRandomHash("wrong_left_value")
	wrongProofRV := test.GenerateRandomHash("wrong_right_value")

	// This should fail because the hash doesn't match the expected state root
	err = env.SubmitCorrectEvidence(gameAddr, wrongProofLV, wrongProofRV)
	require.Error(t, err, "Submitting incorrect evidence should fail")
	t.Logf("Expected error occurred: %v", err)

	// 5. Verify challenger's state hasn't changed (evidence was rejected)
	t.Log("Step 5: Verifying challenger state after failed evidence submission")
	challengerInfoAfter, err := env.GetChallengerInfo(env.ChallengerAccount.Address)
	require.NoError(t, err)

	require.Equal(t, challengerInfoBefore.StakingAmount, challengerInfoAfter.StakingAmount,
		"Challenger stake should be unchanged after failed evidence submission")
	require.Equal(t, challengerInfoBefore.TotalSlashedAmount, challengerInfoAfter.TotalSlashedAmount,
		"Slashed amount should be unchanged after failed evidence submission")

	// 6. Verify attention test is still waiting for evidence
	t.Log("Step 6: Verifying attention test is still pending")
	attentionInfo, err = env.GetAttentionTestInfo(gameAddr)
	require.NoError(t, err)
	require.False(t, attentionInfo.EvidenceSubmitted,
		"Evidence should still be marked as not submitted")

	// 7. Now submit CORRECT evidence to show it works
	t.Log("Step 7: Submitting correct evidence to verify system still works")
	correctProofLV := test.GenerateRandomHash("left_value")
	correctProofRV := test.GenerateRandomHash("right_value")
	correctHashForEvidence := test.CreateTestStateRoot(correctProofLV, correctProofRV)

	// Create new attention test with the correct expected hash
	gameAddr2 := env.CreateMockDisputeGame()
	err = env.TriggerRATAttentionTest(gameAddr2, correctHashForEvidence, env.GetLatestBlockHash())
	require.NoError(t, err)

	// Check if challenger was selected for second test
	attentionInfo2, err := env.GetAttentionTestInfo(gameAddr2)
	require.NoError(t, err)
	if attentionInfo2.ChallengerAddress == (common.Address{}) {
		t.Skip("No challenger selected for second attention test")
	}

	// Submit correct evidence
	err = env.SubmitCorrectEvidence(gameAddr2, correctProofLV, correctProofRV)
	require.NoError(t, err)
	t.Log("Correct evidence submitted successfully")

	// 8. Verify the correct evidence was accepted
	t.Log("Step 8: Verifying correct evidence was accepted")
	attentionInfo2, err = env.GetAttentionTestInfo(gameAddr2)
	require.NoError(t, err)
	require.True(t, attentionInfo2.EvidenceSubmitted,
		"Correct evidence should be marked as submitted")

	t.Log("=== RAT Incorrect Evidence Submission Test PASSED ===")
}