package docker

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func BuildImage(dockerfile, name string) error {
	args := []string{
		"build",
		fmt.Sprintf("--tag=%s", name),
		fmt.Sprintf("--file=%s", dockerfile),
		"--rm",
		filepath.Dir(dockerfile),
	}
	log.Printf("Building server image: docker %s\n", args)
	c := exec.Command("docker", args...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
}

func StartContainer(iName, cName string, port int, cpuset, mem string) error {
	args := []string{
		"run",
		"--rm",
		fmt.Sprintf("--name=%s", cName),
		fmt.Sprintf("--publish=%d:%d", port, port),
		fmt.Sprintf("--memory=%s", mem),
	}
	if len(cpuset) > 0 {
		args = append(args, fmt.Sprintf("--cpuset-cpus=%s", cpuset))
	}
	args = append(args, iName)
	log.Printf("Running: docker %s\n", args)
	c := exec.Command("docker", args...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
}

func StopContainer(cName string) error {
	args := []string{"stop", cName}
	log.Printf("Running: docker %s\n", args)
	c := exec.Command("docker", args...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
}

func IsContainerRunning(cName string) (bool, error) {
	args := []string{
		"ps",
		"--format={{.Names}}",
		fmt.Sprintf("--filter=name=%s", cName),
	}
	log.Printf("Running: docker %s\n", args)
	o, err := exec.Command("docker", args...).Output()
	if err != nil {
		return false, err
	}
	return len(o) != 0, nil
}
