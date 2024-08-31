package deployers

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum-optimism/optimism/op-chain-ops/script"
)

type DeployImplementationsInput struct {
	WithdrawalDelaySeconds          *big.Int
	MinProposalSizeBytes            *big.Int
	ChallengePeriodSeconds          *big.Int
	ProofMaturityDelaySeconds       *big.Int
	DisputeGameFinalityDelaySeconds *big.Int
	// Release version to set OPSM implementations for, of the format `op-contracts/vX.Y.Z`.
	Release               string
	SuperchainConfigProxy common.Address
	ProtocolVersionsProxy common.Address
	// TODO: Interop flag, to deploy OptimismPortalInterop / SystemConfigInterop,
	// By overriding which deploy script we use
}

type DeployImplementationsOutput struct {
	Opsm                             common.Address
	DelayedWETHImpl                  common.Address
	OptimismPortalImpl               common.Address
	PreimageOracleSingleton          common.Address
	MipsSingleton                    common.Address
	SystemConfigImpl                 common.Address
	L1CrossDomainMessengerImpl       common.Address
	L1ERC721BridgeImpl               common.Address
	L1StandardBridgeImpl             common.Address
	OptimismMintableERC20FactoryImpl common.Address
}

type DeployImplementationsScript struct {
	Run func(input, output common.Address) error
}

func DeployImplementations(l1Host *script.Host, input *DeployImplementationsInput) (*DeployImplementationsOutput, error) {
	output := &DeployImplementationsOutput{}
	inputAddr := l1Host.NewScriptAddress()
	outputAddr := l1Host.NewScriptAddress()

	cleanupInput, err := script.WithPrecompileAtAddress[*DeployImplementationsInput](l1Host, inputAddr, input)
	if err != nil {
		return nil, fmt.Errorf("failed to insert DeployImplementationsInput precompile: %w", err)
	}
	defer cleanupInput()

	cleanupOutput, err := script.WithPrecompileAtAddress[*DeployImplementationsOutput](l1Host, outputAddr, output)
	if err != nil {
		return nil, fmt.Errorf("failed to insert DeployImplementationsOutput precompile: %w", err)
	}
	defer cleanupOutput()

	deployScript, cleanupDeploy, err := script.WithScript[DeployImplementationsScript](l1Host, "DeployImplementations.s.sol", "DeployImplementations")
	if err != nil {
		return nil, fmt.Errorf("failed to load DeployImplementations script: %w", err)
	}
	defer cleanupDeploy()

	if err := deployScript.Run(inputAddr, outputAddr); err != nil {
		return nil, fmt.Errorf("failed to run DeployImplementations script: %w", err)
	}

	return output, nil
}
