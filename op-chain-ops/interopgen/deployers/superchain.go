package deployers

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum-optimism/optimism/op-chain-ops/script"
)

type DeploySuperchainInput struct {
	ProxyAdminOwner            common.Address
	ProtocolVersionsOwner      common.Address
	Guardian                   common.Address
	Paused                     bool
	RequiredProtocolVersion    params.ProtocolVersion
	RecommendedProtocolVersion params.ProtocolVersion
}

type DeploySuperchainOutput struct {
	SuperchainProxyAdmin  common.Address
	SuperchainConfigImpl  common.Address
	SuperchainConfigProxy common.Address
	ProtocolVersionsImpl  common.Address
	ProtocolVersionsProxy common.Address
}

type DeploySuperchainScript struct {
	Run func(input, output common.Address) error
}

func DeploySuperchain(l1Host *script.Host, input *DeploySuperchainInput) (*DeploySuperchainOutput, error) {
	output := &DeploySuperchainOutput{}
	inputAddr := l1Host.NewScriptAddress()
	outputAddr := l1Host.NewScriptAddress()

	cleanupInput, err := script.WithPrecompileAtAddress[*DeploySuperchainInput](l1Host, inputAddr, input)
	if err != nil {
		return nil, fmt.Errorf("failed to insert DeploySuperchainInput precompile: %w", err)
	}
	defer cleanupInput()

	cleanupOutput, err := script.WithPrecompileAtAddress[*DeploySuperchainOutput](l1Host, outputAddr, output)
	if err != nil {
		return nil, fmt.Errorf("failed to insert DeploySuperchainOutput precompile: %w", err)
	}
	defer cleanupOutput()

	deployScript, cleanupDeploy, err := script.WithScript[DeploySuperchainScript](l1Host, "DeploySuperchain.s.sol", "DeploySuperchain")
	if err != nil {
		return nil, fmt.Errorf("failed to load DeploySuperchain script: %w", err)
	}
	defer cleanupDeploy()

	if err := deployScript.Run(inputAddr, outputAddr); err != nil {
		return nil, fmt.Errorf("failed to run DeploySuperchain script: %w", err)
	}

	return output, nil
}
