package network

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	tcnetwork "github.com/testcontainers/testcontainers-go/network"
)

type Network struct {
	network *testcontainers.DockerNetwork
}

func NewNetwork(ctx context.Context, projectName string) (*Network, error) {
	net, err := tcnetwork.New(ctx,
		tcnetwork.WithDriver(testcontainers.Bridge),
		tcnetwork.WithAttachable(),
		tcnetwork.WithLabels(map[string]string{
			"project": projectName,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker network: %w", err)
	}

	return &Network{network: net}, nil
}

func (n *Network) Name() string {
	return n.network.Name
}

func (n *Network) Remove(ctx context.Context) error {
	return n.network.Remove(ctx)
}
