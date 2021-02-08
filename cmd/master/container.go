package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"strconv"
)

func (s *CLIServer) GetContainerClient() *client.Client {
	if s.containerClient == nil {
		var err error
		s.containerClient, err = client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			panic(err)
		}
	}

	return s.containerClient
}

func (s *CLIServer) GenContainerName(gNBId uint32) string {
	return "nrsim-" + strconv.Itoa(int(gNBId))
}

func (s *CLIServer) NewWorker(contName string) error {
	cont, err := s.containerClient.ContainerCreate(context.Background(),
		&container.Config{
			Image: "breezestars/nrsim",
		},
		nil,
		&network.NetworkingConfig{},
		nil,
		contName,
	)

	if err != nil {
		return errors.Wrapf(err, "Create worker contaienr failed")
	}

	if err := s.containerClient.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		if err := s.containerClient.ContainerRemove(context.Background(), cont.ID, types.ContainerRemoveOptions{}); err != nil {
			panic(err)
		}
		return errors.Wrapf(err, "Start worker contaienr failed")
	}

	return nil
}
