package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/danielfireman/phd/experiments/golang/rtsimple/docker"
	"github.com/danielfireman/phd/experiments/golang/rtsimple/load"
)

var (
	dfTmpl    = flag.String("dockerfile_template", "Dockerfile.template", "Full path to the Dockerfile template.")
	baseImage = flag.String("base_image", "1.5.3", "Base image to run the server. Example: golang:1.5.3-alpine")
	port      = flag.Int("port", 8999, "Port to run/expose service. Example: '8999'")

	loadN                = flag.Int("load_n", 5000, "Number of requests to generate.")
	loadConcurrencyLevel = flag.Int("load_c", 10, "Number of concurrent workers generating load.")
	loadOps              = flag.Int("load_ops", 1000000, "Numer of operations done per request.")
	loadMem              = flag.Int("load_mem", 1024, "Amount of memory allocated per request (in bytes).")

	cpus       = flag.Int("cpus", 1, "Number of CPUs used by the server container.")
	gomaxprocs = flag.Int("gomaxprocs", 1, "Number of max processors used by golang runtime. Example: 2")
	mem        = flag.String("memory", "1g", "Memory allocated in the container.")
)

const instanceNamePrefix = "restserver_"

func main() {
	flag.Parse()

	exprID := stripChars(*baseImage, ":.-")

	// Docker likes absolute paths for Dockerfiles.
	dockerfilePath := filepath.Join(filepath.Dir(*dfTmpl), fmt.Sprintf("Dockerfile_%s", exprID))
	if !strings.HasPrefix(dockerfilePath, "/") {
		v, ok := os.LookupEnv("PWD")
		if !ok {
			log.Fatalf("Envvar PWD not set.")
		}
		dockerfilePath = filepath.Join(v, dockerfilePath)
	}
	if err := createDockerfile(*dfTmpl, dockerfilePath, *baseImage, *port, *gomaxprocs); err != nil {
		log.Fatal(err)
	}
	log.Println("Dockerfile created: ", dockerfilePath, "\n####\n")

	// Building image.
	iName := fmt.Sprintf("danielfireman/phd-experiments:%s%s", instanceNamePrefix, exprID)
	if err := docker.BuildImage(dockerfilePath, iName); err != nil {
		log.Fatal(err)
	}
	log.Println("Server image created: ", iName, "\n####\n")

	// Dealing with containers.
	cName := fmt.Sprintf("%s%s", instanceNamePrefix, exprID)

	// Check if container is running and stop, if it is the case.
	log.Printf("Checking if container %s is running.", cName)
	if ok, err := docker.IsContainerRunning(cName); err != nil {
		log.Fatal(err)
	} else if ok {
		log.Printf("Container %s is running, stopping it.", cName)
		if err := docker.StopContainer(cName); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("Container %s is not running, starting it.", cName)
	}

	// Starting container.
	cErrChan := make(chan error)
	go func() {
		cErrChan <- docker.StartContainer(iName, cName, *port, *cpus, *mem)
	}()
	if err := waitHealth(serverURI(*port, "ping"), cErrChan); err != nil {
		log.Fatal(err)
	}
	log.Printf("Server is up, running and healthy: %s\n####\n", cName)

	// Stressing service.
	g := load.Generator{
		N:                *loadN,
		ConcurrencyLevel: *loadConcurrencyLevel,
		URL:              serverURI(*port, fmt.Sprintf("work&cpu=%d?mem=%d", *loadOps, *loadMem)),
	}
	log.Printf("Start stressing server. #req:%d concurrencyLevel:%d url:%s", g.N, g.ConcurrencyLevel, g.URL)
	log.Printf("%+v", g.Run())
	log.Println("Finish stressing server\n####\n")

	// Stopping container.
	if err := docker.StopContainer(cName); err != nil {
		log.Fatal(err)
	}
	log.Println("Container stopped: ", cName, "\n####\n")

	if err := os.Remove(dockerfilePath); err != nil {
		log.Fatal(err)
	}
	log.Println("Docker file successfully removed: ", dockerfilePath, "\n####\n")
}

func createDockerfile(tmplPath, path, baseImage string, port, gomaxprocs int) error {
	t, err := template.New(filepath.Base(tmplPath)).ParseFiles(tmplPath)
	if err != nil {
		return err
	}
	dockerfile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dockerfile.Close()
	err = t.Execute(dockerfile, struct {
		BaseImage  string
		Port       string
		GOMAXPROCS int
	}{
		baseImage,
		fmt.Sprintf("%d", port),
		gomaxprocs,
	})
	if err != nil {
		return err
	}
	return nil
}

func serverURI(port int, path string) string {
	return fmt.Sprintf("http://localhost:%d/%s", port, path)
}

const (
	pingInterval      = 5 * time.Second
	waitHealthTimeout = 1 * time.Minute
)

func waitHealth(pingUrl string, cErrChan <-chan error) error {
	ticker := time.Tick(pingInterval)
	timeout := time.After(waitHealthTimeout)
	log.Printf("Waiting for server to start. Ping URL: %s\n", pingUrl)
	for {
		select {
		case <-timeout:
			return fmt.Errorf("Timed out waiting for server to be health.")
		case err := <-cErrChan:
			return err
		case <-ticker:
			log.Printf("Sending ping to URL: %s\n", pingUrl)
			resp, err := http.Get(pingUrl)
			if err == nil && resp.StatusCode == http.StatusOK {
				return nil
			}
		}
	}
}

func stripChars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}
