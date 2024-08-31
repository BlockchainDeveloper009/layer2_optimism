package deployers

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum-optimism/optimism/op-chain-ops/genesis"
	"github.com/ethereum-optimism/optimism/op-chain-ops/script"
)

var (
	// address(uint160(uint256(keccak256(abi.encode("optimism.deployconfig"))))) - not a simple hash, due to ABI encode
	deployConfigAddr       = common.HexToAddress("0x9568d36E291c2C4c34fa5593fcE73715abEf6F9c")
	deploymentRegistryAddr = common.Address(crypto.Keccak256([]byte("optimism.deploymentregistry"))[12:])
)

type L2GenesisInput struct {
	L1Deployments struct {
		L1CrossDomainMessengerProxy common.Address
		L1StandardBridgeProxy       common.Address
		L1ERC721BridgeProxy         common.Address
	}
	L2Config genesis.L2InitializationConfig
}

type L2GenesisScript struct {
	RunWithEnv func() error
}

func L2Genesis(l2Host *script.Host, input *L2GenesisInput) error {
	deploymentRegistry := &DeploymentRegistryPrecompile{
		Deployments: map[string]common.Address{
			"L1CrossDomainMessengerProxy": input.L1Deployments.L1CrossDomainMessengerProxy,
			"L1StandardBridgeProxy":       input.L1Deployments.L1StandardBridgeProxy,
			"L1ERC721BridgeProxy":         input.L1Deployments.L1ERC721BridgeProxy,
		},
	}

	cleanupDeploymentRegistry, err := script.WithPrecompileAtAddress[*DeploymentRegistryPrecompile](
		l2Host, deploymentRegistryAddr, deploymentRegistry)
	if err != nil {
		return fmt.Errorf("failed to insert DeploymentRegistry precompile: %w", err)
	}
	defer cleanupDeploymentRegistry()

	deployConfig := &genesis.DeployConfig{
		L2InitializationConfig: input.L2Config,
	}
	cleanupDeployConfig, err := script.WithPrecompileAtAddress[*genesis.DeployConfig](l2Host, deployConfigAddr, deployConfig, script.WithFieldsOnly[*genesis.DeployConfig])
	if err != nil {
		return fmt.Errorf("failed to insert DeployConfig precompile: %w", err)
	}
	defer cleanupDeployConfig()

	l2GenesisScript, cleanupL2Genesis, err := script.WithScript[L2GenesisScript](l2Host, "L2Genesis.s.sol", "L2Genesis")
	if err != nil {
		return fmt.Errorf("failed to load L2Genesis script: %w", err)
	}
	defer cleanupL2Genesis()

	if err := l2GenesisScript.RunWithEnv(); err != nil {
		return fmt.Errorf("failed to run L2 genesis script: %w", err)
	}
	return nil
}
