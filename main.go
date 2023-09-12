package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	g, gCtx := errgroup.WithContext(ctx)
	for i := 0; i < 3; i++ {
		g.Go(func() error {
			// run docker pull nginx:latest first
			containerID, err := startContainer(ctx, cli, "nginx:latest")
			if err != nil {
				log.Fatal(err)
			}
			defer cleanupContainer(gCtx, containerID, cli)

			ctrs, err := cli.ContainerList(gCtx, types.ContainerListOptions{
				Filters: filters.NewArgs(filters.Arg("id", containerID)),
			})
			if err != nil {
				log.Fatal(err)
			}
			if len(ctrs) != 1 {
				log.Fatal("could not find container")
			}
			log.Println("(Before sleeping)")
			log.Printf("Total ports: %d\n", len(ctrs[0].Ports))
			log.Printf("Ctr ports: %v\n", ctrs[0].Ports)
			log.Printf("Host port: %d\n", ctrs[0].Ports[0].PublicPort)

			time.Sleep(3 * time.Second)

			ctrs, err = cli.ContainerList(ctx, types.ContainerListOptions{
				Filters: filters.NewArgs(filters.Arg("id", containerID)),
			})
			if err != nil {
				log.Fatal(err)
			}
			if len(ctrs) != 1 {
				log.Fatal("could not find container")
			}

			log.Println("(After sleeping)")
			log.Printf("Total ports: %d\n", len(ctrs[0].Ports))
			log.Printf("Ctr ports: %v\n", ctrs[0].Ports)
			log.Printf("Host port: %d\n", ctrs[0].Ports[0].PublicPort)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

}

func startContainer(ctx context.Context, dockerClient client.APIClient, image string) (string, error) {
	res, err := dockerClient.ContainerCreate(ctx,
		&container.Config{
			Image: image,
			ExposedPorts: nat.PortSet{
				"8080/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			AutoRemove: true,
			PortBindings: nat.PortMap{
				"8080/tcp": []nat.PortBinding{
					{
						HostIP: "0.0.0.0",
						// HostPort: // // not set on purpose to get one available assigned automatically
					},
				},
			},
		},
		&network.NetworkingConfig{},
		&ocispec.Platform{OS: "linux", Architecture: runtime.GOARCH},
		fmt.Sprintf("scout_local_policy_evaluation_%s", uuid.New().String()))
	if err != nil {
		return res.ID, fmt.Errorf("could not create container: %w", err)
	}
	if err = dockerClient.ContainerStart(ctx, res.ID, types.ContainerStartOptions{}); err != nil {
		return res.ID, fmt.Errorf("could not start container: %w", err)
	}
	return res.ID, nil
}

func cleanupContainer(ctx context.Context, containerID string, dockerClient client.APIClient) {
	timeout := 0
	_ = dockerClient.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
	_ = dockerClient.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
}
