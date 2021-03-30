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
	"time"
)

const (
	ContainerImageName  = "breezestars/nrsim:dev"
	MasterServerAddress = "172.17.0.1:50051"
)

type GnbConfig struct {
	ContainerId string
	Config      *api.GnbConfig
	Ip          string
	Registered  bool

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

func (s *CLIServer) NewWorker(contName string) (string, string, error) {
	client := s.GetContainerClient()
	cont, err := client.ContainerCreate(context.Background(),
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
		return "", "", errors.Wrapf(err, "Create worker contaienr failed")
	}

	if err := client.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		if err := client.ContainerRemove(context.Background(), cont.ID, types.ContainerRemoveOptions{}); err != nil {
			panic(err)
		}
		return "", "", errors.Wrapf(err, "Start worker contaienr failed")
	}

	//TODO: Get network IP and return so can bind in map.
	inspect, err := client.ContainerInspect(context.Background(), cont.ID)
	if err != nil {
		errLog.Printf("%+v", err)
		//return "", "", errors.Wrapf(err, "Get IP from worker contaienr failed")
	}

	for i := 0; i < 5; i++ {
		if inspect.NetworkSettings.IPAddress != "" {
			return cont.ID, inspect.NetworkSettings.IPAddress, nil
		}
		time.Sleep(time.Millisecond)
	}

	return cont.ID, "", errors.New("Did not get IP address")
}

func (s *CLIServer) DelWorker(contId, contName string) error {
	if err := s.GetContainerClient().ContainerRemove(context.Background(), contId, types.ContainerRemoveOptions{Force: true}); err != nil {
		return errors.Wrapf(err, "Delete worker container failed")
	}

	return nil
}
