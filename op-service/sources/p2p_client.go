package sources

import (
	"context"

	"github.com/tokamak-network/tokamak-thanos/op-service/apis"
	"github.com/tokamak-network/tokamak-thanos/op-service/client"
)

const P2PNamespaceRPC = "opp2p"

// P2PClient wraps an RPC client for P2P API calls.
type P2PClient struct {
	client client.RPC
}

// NewP2PClient creates a new P2PClient.
func NewP2PClient(client client.RPC) *P2PClient {
	return &P2PClient{client: client}
}

func (c *P2PClient) Self(ctx context.Context) (*apis.PeerInfo, error) {
	var out apis.PeerInfo
	err := c.client.CallContext(ctx, &out, P2PNamespaceRPC+"_self")
	return &out, err
}

func (c *P2PClient) Peers(ctx context.Context, connected bool) (*apis.PeerDump, error) {
	var out apis.PeerDump
	err := c.client.CallContext(ctx, &out, P2PNamespaceRPC+"_peers", connected)
	return &out, err
}

func (c *P2PClient) PeerStats(ctx context.Context) (*apis.PeerStats, error) {
	var out apis.PeerStats
	err := c.client.CallContext(ctx, &out, P2PNamespaceRPC+"_peerStats")
	return &out, err
}

func (c *P2PClient) DiscoveryTable(ctx context.Context) ([]*apis.PeerInfo, error) {
	var out []*apis.PeerInfo
	err := c.client.CallContext(ctx, &out, P2PNamespaceRPC+"_discoveryTable")
	return out, err
}

func (c *P2PClient) BlockAddr(ctx context.Context, ip string) error {
	return c.client.CallContext(ctx, nil, P2PNamespaceRPC+"_blockAddr", ip)
}

func (c *P2PClient) UnblockAddr(ctx context.Context, ip string) error {
	return c.client.CallContext(ctx, nil, P2PNamespaceRPC+"_unblockAddr", ip)
}

func (c *P2PClient) BlockPeer(ctx context.Context, id string) error {
	return c.client.CallContext(ctx, nil, P2PNamespaceRPC+"_blockPeer", id)
}

func (c *P2PClient) UnblockPeer(ctx context.Context, id string) error {
	return c.client.CallContext(ctx, nil, P2PNamespaceRPC+"_unblockPeer", id)
}

func (c *P2PClient) BlockSubnet(ctx context.Context, subnet string) error {
	return c.client.CallContext(ctx, nil, P2PNamespaceRPC+"_blockSubnet", subnet)
}

func (c *P2PClient) UnblockSubnet(ctx context.Context, subnet string) error {
	return c.client.CallContext(ctx, nil, P2PNamespaceRPC+"_unblockSubnet", subnet)
}

func (c *P2PClient) ProtectPeer(ctx context.Context, id string) error {
	return c.client.CallContext(ctx, nil, P2PNamespaceRPC+"_protectPeer", id)
}

func (c *P2PClient) UnprotectPeer(ctx context.Context, id string) error {
	return c.client.CallContext(ctx, nil, P2PNamespaceRPC+"_unprotectPeer", id)
}

func (c *P2PClient) ConnectPeer(ctx context.Context, addr string) error {
	return c.client.CallContext(ctx, nil, P2PNamespaceRPC+"_connectPeer", addr)
}

func (c *P2PClient) DisconnectPeer(ctx context.Context, id string) error {
	return c.client.CallContext(ctx, nil, P2PNamespaceRPC+"_disconnectPeer", id)
}
