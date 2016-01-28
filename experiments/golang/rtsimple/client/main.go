package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

var (
	dfTmpl    = flag.String("dockerfile_template", "", "Full path to the Dockerfile template.")
	goVersion = flag.String("go_version", "1.5.3", "Golang version. Example: 1.5.3")
	port      = flag.Int("port", 8999, "Port to run/expose service. Example: '8999'")

	loadN                = flag.Int("load_n", 500, "Number of requests to generate.")
	loadConcurrencyLevel = flag.Int("load_c", 10, "Number of concurrent workers generating load.")

	cpus = flag.Int("cpus", 1, "Number of CPUs used by the server.")
	mem  = flag.String("memory", "1g", "Memory allocated in the container.")
)

func main() {
	flag.Parse()
	log.Println(*dfTmpl)
	dockerfile, err := createDockerfile(*dfTmpl, *goVersion, *port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Dockerfile created: ", dockerfile, "\n####\n")

	image, err := buildServerImage(dockerfile, *goVersion)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server image created: ", image, "\n####\n")

	errChan := startServer(image, *port, *cpus, *mem, *goVersion)
	if err := <-errChan; err != nil {
		log.Fatal(err)
	}
	log.Println("Server container started.")
}

func createDockerfile(tmplPath, goVersion string, port int) (string, error) {
	t, err := template.New(filepath.Base(tmplPath)).ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}
	path := filepath.Join(filepath.Dir(tmplPath), fmt.Sprintf("Dockerfile_%s", goVersion))
	dockerfile, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dockerfile.Close()
	err = t.Execute(dockerfile, struct {
		Version,
		Port string
	}{
		goVersion,
		fmt.Sprintf("%d", port),
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

func buildServerImage(dockerfile, goVersion string) (string, error) {
	name := fmt.Sprintf("danielfireman/phd-experiments:restserver_go%s", goVersion)
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
	return name, c.Run()
}

func startServer(name string, port, cpus int, mem, goVersion string) <-chan error {
	errChan := make(chan error)
	args := []string{
		"run",
		fmt.Sprintf("--publish=%d:%d", port, port),
		fmt.Sprintf("--memory=%s", mem),
		fmt.Sprintf("--cpuset-cpus=%d", cpus),
		name,
	}
	log.Printf("Starting server: docker %s\n", args)
	cmd := exec.Command("docker", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	go func() {
		if err := cmd.Start(); err != nil {
			errChan <- err
			return
		}
	}()
	return errChan
}
