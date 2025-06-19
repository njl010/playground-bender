package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type inspectResult struct {
	NetworkSettings struct {
		Networks map[string]struct {
			IPAddress string `json:"IPAddress"`
		} `json:"Networks"`
	} `json:"NetworkSettings"`
}

type Container struct {
	ID     string
	Name   string
	Status string
}

func RunCustomContainer(ctx context.Context, containerName, imageName, networkName string) (string, string, error) {
	cmd := exec.CommandContext(ctx, "docker", "run",
		"--name", containerName,
		"--privileged",
		"--volume", "/dev/kvm:/dev/kvm",
		"--device", "/dev/kvm",
		"--cap-add", "NET_ADMIN",
		"--cap-add", "SYS_ADMIN",
		"--network", networkName,
		"--volume", "/home/john/resourses:/assets",
		// "--ip", ip,
		"-d", // Run in detached mode
		imageName,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("failed to run container: %v\nOutput: %s", err, string(output))
	}

	containerID := strings.TrimSpace(string(output))

	inspectCmd := exec.CommandContext(ctx, "docker", "inspect", containerName)
	inspectOut, err := inspectCmd.Output()
	if err != nil {
		return containerID, "", fmt.Errorf("failed to inspect container: %v", err)
	}

	var result []inspectResult
	if err := json.Unmarshal(inspectOut, &result); err != nil {
		return containerID, "", fmt.Errorf("failed to parse docker inspect output: %v", err)
	}

	ip := result[0].NetworkSettings.Networks[networkName].IPAddress

	fmt.Printf("Container %s started with IP %s on network %s\n", containerName, ip, networkName)
	return containerID, ip, nil
}

func DeleteContainerByName(ctx context.Context, containerName string) error {
	// Step 1: Stop the container (optional, safe cleanup)
	stopCmd := exec.CommandContext(ctx, "docker", "stop", containerName)
	stopOut, stopErr := stopCmd.CombinedOutput()
	if stopErr != nil {
		return fmt.Errorf("failed to stop container: %v\nOutput: %s", stopErr, string(stopOut))
	}

	// Step 2: Remove the container
	rmCmd := exec.CommandContext(ctx, "docker", "rm", containerName)
	rmOut, rmErr := rmCmd.CombinedOutput()
	if rmErr != nil {
		return fmt.Errorf("failed to remove container: %v\nOutput: %s", rmErr, string(rmOut))
	}

	fmt.Printf("Container %s stopped and removed successfully\n", containerName)
	return nil
}

func ListPlaygroundContainers() ([]Container, error) {
	ctx := context.Background()

	// Run docker ps with filter and custom format
	cmd := exec.CommandContext(ctx, "docker", "ps",
		"--filter", "name=^/playground",
		"--format", "{{.ID}}|{{.Names}}|{{.Status}}",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v\nOutput: %s", err, string(output))
	}

	lines := strings.Split(string(output), "\n")
	containers := []Container{}

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) >= 3 {
			containers = append(containers, Container{
				ID:     parts[0],
				Name:   parts[1],
				Status: parts[2],
			})
		}
	}

	return containers, nil
}
