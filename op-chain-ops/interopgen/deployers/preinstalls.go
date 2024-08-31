package deployers

import (
	"fmt"

	"github.com/ethereum-optimism/optimism/op-chain-ops/genesis"
	"github.com/ethereum-optimism/optimism/op-chain-ops/script"
)

type PreinstallsScript struct {
	SetPreinstalls func() error
}

func InsertPreinstalls(host *script.Host) error {
	// Init L2Genesis script. Yes, this might be L1. Hack to deploy all preinstalls.
	l2GenesisScript, cleanupL2Genesis, err := script.WithScript[PreinstallsScript](host, "L2Genesis.s.sol", "L2Genesis")
	if err != nil {
		return fmt.Errorf("failed to load L2Genesis script for preinstalls work: %w", err)
	}
	defer cleanupL2Genesis()

	// We need the Chain ID for the preinstalls setter to work
	deployConfig := &genesis.DeployConfig{}
	chainID := host.ChainID()
	if !chainID.IsUint64() {
		return fmt.Errorf("preinstalls script expects uint64 chainID, but got %d (bitlen %d)", chainID, chainID.BitLen())
	}
	deployConfig.L2ChainID = chainID.Uint64()
	cleanupDeployConfig, err := script.WithPrecompileAtAddress[*genesis.DeployConfig](host, deployConfigAddr, deployConfig, script.WithFieldsOnly[*genesis.DeployConfig])
	if err != nil {
		return fmt.Errorf("failed to insert DeployConfig precompile: %w", err)
	}
	defer cleanupDeployConfig()

	if err := l2GenesisScript.SetPreinstalls(); err != nil {
		return fmt.Errorf("failed to set preinstalls: %w", err)
	}
	return nil
}
