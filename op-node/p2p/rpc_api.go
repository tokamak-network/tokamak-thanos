package p2p

import (
	"github.com/tokamak-network/tokamak-thanos/op-service/apis"
	"github.com/tokamak-network/tokamak-thanos/op-service/sources"
)

// backwards compatibility
var NamespaceRPC = sources.P2PNamespaceRPC

type API = apis.P2PClient
type PeerInfo = apis.PeerInfo
type PeerDump = apis.PeerDump
type PeerStats = apis.PeerStats
