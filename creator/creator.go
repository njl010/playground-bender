package creator

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	converters "github.com/devdevaraj/bender/converter"
	"github.com/devdevaraj/bender/docker"
)

type Network struct {
	ID     string `json:"id"`
	CID    string `json:"cid"`
	Name   string `json:"name"`
	CName  string `json:"cname"`
	Status string `json:"status"`
	Driver string `json:"driver"`
	Subnet string `json:"subnet,omitempty"`
}

type NetworkInspect []struct {
	IPAM struct {
		Config []struct {
			Subnet  string `json:"Subnet"`
			Gateway string `json:"Gateway"`
		} `json:"Config"`
	} `json:"IPAM"`
}

func CreateDockerBridge(name string, image string) (string, string, string, string, string, error) {
	ctx := context.Background()

	cmd := exec.CommandContext(ctx, "docker", "network", "create",
		"--driver", "bridge",
		"bg-"+name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to create network: %v, output: %s", err, string(output))
	}

	inspectCmd := exec.CommandContext(ctx, "docker", "network", "inspect", "bg-"+name)
	outputi, errI := inspectCmd.Output()
	if errI != nil {
		return "", "", "", "", "", fmt.Errorf("failed to inspect network: %v, output: %s", errI, string(outputi))
	}

	var result NetworkInspect
	err = json.Unmarshal(outputi, &result)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to parse json: %v, output: %s", err, string(outputi))
	}

	ip, err := converters.NextIP(result[0].IPAM.Config[0].Gateway)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to parse ip: %v, output: %s", err, string(ip))
	}

	var cid = ""
	cid, ip, err = docker.RunCustomContainer(ctx, "playground-"+name, image, "bg-"+name)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to start container: %v, output: %s", err, string(ip))
	}

	networkID := strings.TrimSpace(string(output))

	return networkID, cid, ip, result[0].IPAM.Config[0].Subnet, result[0].IPAM.Config[0].Gateway, nil
}

func DeleteDockerBridge(name string) error {
	ctx := context.Background()

	err := docker.DeleteContainerByName(ctx, "playground-"+name)
	if err != nil {
		return fmt.Errorf("failed to delete playground: %v", err)
	}

	// Build the docker network rm command
	cmd := exec.CommandContext(ctx, "docker", "network", "rm", "bg-"+name)

	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete network: %v, output: %s", err, string(output))
	}

	return nil
}

func ListDockerBridges() ([]Network, error) {
	ctx := context.Background()

	// Build the docker network ls command with format option for json output
	cmd := exec.CommandContext(ctx, "docker", "network", "ls", "--filter", "name=^bg", "--format", "{{.ID}}|{{.Name}}|{{.Driver}}")

	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %v, output: %s", err, string(output))
	}

	coutput, cerr := docker.ListPlaygroundContainers()
	if cerr != nil {
		return nil, fmt.Errorf("failed to list networks: %v, output: %s", cerr, string(output))
	}

	// Parse the output
	networks := []Network{}
	lines := strings.Split(string(output), "\n")

	for i, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) >= 3 {
			networks = append(networks, Network{
				ID:     parts[0],
				CID:    coutput[len(coutput)-i-1].ID,
				Name:   parts[1],
				CName:  coutput[len(coutput)-i-1].Name,
				Status: coutput[len(coutput)-i-1].Status,
				Driver: parts[2],
			})
		}
	}

	return networks, nil
}
