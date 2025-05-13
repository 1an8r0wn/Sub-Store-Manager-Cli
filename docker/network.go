package docker

import (
	"context"
	"fmt"

	networkType "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var (
	ctx = context.Background()
	cli *client.Client
)

func (c *Container) networkIsExist() (networkType.Inspect, error) {
	fmt.Println("Checking network if exist...")
	networks, err := cli.NetworkList(ctx, networkType.ListOptions{})
	if err != nil {
		return networkType.Inspect{}, err
	}

	for _, network := range networks {
		if network.Name == c.Network {
			return network, nil
		}
	}

	return networkType.Inspect{}, nil
}

func (c *Container) createNetwork() (string, error) {
	fmt.Println("Creating network...")
	network, err := cli.NetworkCreate(ctx, c.Network, networkType.CreateOptions{
		Driver: "bridge",
	})
	if err != nil {
		return "", err
	}

	return network.ID, nil
}

func (c *Container) GetNetworkID() (string, error) {
	fmt.Println("Getting network ID...")

	dClient, err := client.NewClientWithOpts()
	if err != nil {
		fmt.Println("Failed to create docker client:", err)
		return "", err
	}
	cli = dClient

	network, err := c.networkIsExist()
	if err != nil {
		fmt.Println("Failed to get network:", err)
		return "", err
	}

	if network.ID == "" {
		fmt.Println("Network not exist, creating...")
		id, err := c.createNetwork()
		if err != nil {
			fmt.Println("Failed to create network:", err)
			return "", err
		}
		return id, nil
	}

	fmt.Println("Network exist, using...")
	return network.ID, nil
}
