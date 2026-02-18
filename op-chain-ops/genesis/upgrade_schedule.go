package genesis

// UpgradeScheduleDeployConfig configures which hardforks to activate.
type UpgradeScheduleDeployConfig struct{}

// ActivateForkAtGenesis activates the given fork at genesis time (time=0).
func (u *UpgradeScheduleDeployConfig) ActivateForkAtGenesis(fork interface{}) {
	// stub: hardfork activation is configured elsewhere for devnet
}
