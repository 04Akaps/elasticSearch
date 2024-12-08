package docker

import (
	"bytes"
	"fmt"
	"github.com/04Akaps/elasticSearch.git/config"
	"os"
	"os/exec"
	"strings"
)

func Initialize(cfg config.Config) {

	err := os.Chdir("../docker")

	if err != nil {
		panic(err)
	}

	if cfg.Docker.Init {
		removeImage(cfg)
	}
}

func removeImage(cfg config.Config) {
	cmd := exec.Command("docker", "ps", "--format", "{{.Names}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		panic(err)
	}

	containerNames := strings.Split(out.String(), "\n")
	targets := cfg.Docker.Targets

	for _, name := range containerNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		for _, target := range targets {
			if name == target {
				fmt.Printf("Stopping container: %s\n", name)
				stopCmd := exec.Command("docker", "stop", name)
				stopCmd.Stdout = os.Stdout
				stopCmd.Stderr = os.Stderr

				if err = stopCmd.Run(); err != nil {
					fmt.Printf("Failed to stop container %s: %v\n", name, err)
				} else {
					fmt.Printf("Successfully stopped container: %s\n", name)
				}
			}
		}
	}
}
