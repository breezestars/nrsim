package main

import (
	"context"
	"github.com/cmingou/nrsim/internal/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"strconv"
)

const (
	ContainerImageName  = "breezestars/nrsim:dev"
	MasterServerAddress = "172.17.0.1:50051"
)

type GnbConfig struct {
	ContainerId string
	Config      *api.GnbConfig
}

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

func (s *CLIServer) NewWorker(contName string) (string, error) {
	cont, err := s.GetContainerClient().ContainerCreate(context.Background(),
		&container.Config{
			Image: ContainerImageName,
			Cmd:   strslice.StrSlice{"-masterSrvIp", MasterServerAddress},
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		nil,
		contName,
	)

	if err != nil {
		return "", errors.Wrapf(err, "Create worker contaienr failed")
	}

	if err := s.GetContainerClient().ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		if err := s.GetContainerClient().ContainerRemove(context.Background(), cont.ID, types.ContainerRemoveOptions{}); err != nil {
			panic(err)
		}
		return "", errors.Wrapf(err, "Start worker contaienr failed")
	}

	return cont.ID, nil
}

func (s *CLIServer) DelWorker(contId, contName string) error {
	if err := s.GetContainerClient().ContainerRemove(context.Background(), contId, types.ContainerRemoveOptions{Force: true}); err != nil {
		return errors.Wrapf(err, "Delete worker container failed")
	}

	return nil
}
