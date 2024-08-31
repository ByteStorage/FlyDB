package service

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Container struct {
	testcontainers.Container
	URI string
}

func (c *Container) setUp(ctx context.Context, serviceContainer testcontainers.Container,
	exposedPort nat.Port) {
	var (
		host string
		port nat.Port
		err  error
	)

	if host, err = serviceContainer.Host(ctx); err != nil {
		fmt.Println("get grpc server container host error: ", err)
		return
	}

	if port, err = serviceContainer.MappedPort(ctx, exposedPort); err != nil {
		fmt.Println("get grpc server container port error: ", err)
		return
	}

	c.URI = fmt.Sprintf("%s:%s", host, port.Port())
	c.Container = serviceContainer
}

func (c *Container) getCloseServiceContainer(ctx context.Context, serviceContainer testcontainers.Container) func() {
	return func() {
		fmt.Println("close grpc server container")
		var (
			err error
		)

		if err = serviceContainer.Terminate(ctx); err != nil {
			fmt.Println("close grpc server container error: ", err)
		}
	}
}

func (c *Container) StartServiceContainer(ctx context.Context, servicePort string) (func(),
	error) {
	fmt.Println("step up flyDB grpc server container")
	var (
		serviceContainer testcontainers.Container
		err              error
		req              = testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "../..",
				Dockerfile: "./docker/Dockerfile",
			},
			ExposedPorts: []string{servicePort},
			WaitingFor:   wait.ForListeningPort(nat.Port(servicePort)),
		}
	)

	if serviceContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}); err != nil {
		return nil, err
	}

	c.setUp(ctx, serviceContainer, nat.Port(servicePort))

	return c.getCloseServiceContainer(ctx, serviceContainer), nil
}
