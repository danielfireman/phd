package docker

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func CreateDockerfile(tmplPath, baseImage, exprID string, serverPort int) (string, error) {
	// Docker likes absolute paths for Dockerfiles.
	dockerfilePath := filepath.Join(filepath.Dir(tmplPath), fmt.Sprintf("Dockerfile_%s", exprID))
	if !strings.HasPrefix(dockerfilePath, "/") {
		v, ok := os.LookupEnv("PWD")
		if !ok {
			return "", fmt.Errorf("Envvar PWD not set.")
		}
		dockerfilePath = filepath.Join(v, dockerfilePath)
	}
	t, err := template.New(filepath.Base(tmplPath)).ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}
	dockerfile, err := os.Create(dockerfilePath)
	if err != nil {
		return "", err
	}
	defer dockerfile.Close()
	if err = t.Execute(dockerfile, struct {
		BaseImage string
		Port      string
	}{
		baseImage,
		fmt.Sprintf("%d", serverPort),
	}); err != nil {
		return "", err
	}
	return dockerfilePath, nil
}

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

const (
	pingInterval      = 5 * time.Second
	waitHealthTimeout = 1 * time.Minute
)

func StartContainer(iName, cName, pingURL string, port int) error {
	errChan := make(chan error)
	go func() {
		args := []string{
			"run",
			"--rm",
			fmt.Sprintf("--name=%s", cName),
			fmt.Sprintf("--publish=%d:%d", port, port),
			iName,
		}
		log.Printf("Running: docker %s\n", args)
		c := exec.Command("docker", args...)
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		if err := c.Run(); err != nil {
			errChan <- err
		}
	}()

	// Waiting for server to be healthy.
	ticker := time.Tick(pingInterval)
	timeout := time.After(waitHealthTimeout)
	log.Printf("Waiting for server to start. Ping URL: %s\n", pingURL)
	for {
		select {
		case <-timeout:
			return fmt.Errorf("Timed out waiting for server to be health.")
		case err := <-errChan:
			return err
		case <-ticker:
			log.Printf("Sending ping to URL: %s\n", pingURL)
			resp, err := http.Get(pingURL)
			if err == nil && resp.StatusCode == http.StatusOK {
				return nil
			}
		}
	}
	return nil
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
