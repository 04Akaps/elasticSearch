package docker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Initialize() {

	err := os.Chdir("../docker")

	if err != nil {
		panic(err)
	}

	removeImage()
}

func removeImage() {
	cmd := exec.Command("docker", "ps", "--format", "{{.Names}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		panic(err)
	}

	containerNames := strings.Split(out.String(), "\n")
	targets := []string{"es01", "kibana"}

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
