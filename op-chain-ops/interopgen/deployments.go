package interopgen

import (
	"github.com/ethereum/go-ethereum/common"
)

type L1Deployment struct {
	// No global deployed contracts that aren't part of the superchain, yet.
}

type Implementations struct {
	Opsm                             common.Address `json:"OPSM"` // not proxied
	DelayedWETHImpl                  common.Address `json:"DelayedWETH"`
	OptimismPortalImpl               common.Address `json:"OptimismPortal"`
	PreimageOracleSingleton          common.Address `json:"PreimageOracleSingleton"`
	MipsSingleton                    common.Address `json:"MipsSingleton"`
	SystemConfigImpl                 common.Address `json:"SystemConfig"`
	L1CrossDomainMessengerImpl       common.Address `json:"L1CrossDomainMessenger"`
	L1ERC721BridgeImpl               common.Address `json:"L1ERC721Bridge"`
	L1StandardBridgeImpl             common.Address `json:"L1StandardBridge"`
	OptimismMintableERC20FactoryImpl common.Address `json:"OptimismMintableERC20Factory"`
}

type SuperchainDeployment struct {
	Implementations

	ProxyAdmin common.Address `json:"ProxyAdmin"`

	ProtocolVersions      common.Address `json:"ProtocolVersions"`
	ProtocolVersionsProxy common.Address `json:"ProtocolVersionsProxy"`

	SuperchainConfig      common.Address `json:"SuperchainConfig"`
	SuperchainConfigProxy common.Address `json:"SuperchainConfigProxy"`
}

type L2OpchainDeployment struct {
	OpChainProxyAdmin                 common.Address `json:"OpChainProxyAdmin"`
	AddressManager                    common.Address `json:"AddressManager"`
	L1ERC721BridgeProxy               common.Address `json:"L1ERC721BridgeProxy"`
	SystemConfigProxy                 common.Address `json:"SystemConfigProxy"`
	OptimismMintableERC20FactoryProxy common.Address `json:"OptimismMintableERC20FactoryProxy"`
	L1StandardBridgeProxy             common.Address `json:"L1StandardBridgeProxy"`
	L1CrossDomainMessengerProxy       common.Address `json:"L1CrossDomainMessengerProxy"`
	// Fault proof contracts below.
	OptimismPortalProxy                common.Address `json:"OptimismPortalProxy"`
	DisputeGameFactoryProxy            common.Address `json:"DisputeGameFactoryProxy"`
	DisputeGameFactoryImpl             common.Address `json:"DisputeGameFactoryImpl"`
	AnchorStateRegistryProxy           common.Address `json:"AnchorStateRegistryProxy"`
	AnchorStateRegistryImpl            common.Address `json:"AnchorStateRegistryImpl"`
	FaultDisputeGame                   common.Address `json:"FaultDisputeGame"`
	PermissionedDisputeGame            common.Address `json:"PermissionedDisputeGame"`
	DelayedWETHPermissionedGameProxy   common.Address `json:"DelayedWETHPermissionedGameProxy"`
	DelayedWETHPermissionlessGameProxy common.Address `json:"DelayedWETHPermissionlessGameProxy"`
}

type L2Deployment struct {
	L2OpchainDeployment

	// In the future this may contain optional extras,
	// e.g. a Safe that will own the L2 chain contracts
}

type WorldDeployment struct {
	L1         *L1Deployment            `json:"L1"`
	Superchain *SuperchainDeployment    `json:"Superchain"`
	L2s        map[string]*L2Deployment `json:"L2s"`
}
