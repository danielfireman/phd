package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pivotal-golang/bytefmt"

	"github.com/danielfireman/phd/experiments/golang/rtsimple/docker"
	"github.com/danielfireman/phd/experiments/golang/rtsimple/load"
	"github.com/danielfireman/phd/experiments/golang/rtsimple/serverinfo"
)

var (
	logToStderr = flag.Bool("logtostderr", false, "")
	repeat      = flag.Int("repeat", 1, "Number of repetitions")

	dfTmpl    = flag.String("dockerfile_template", "Dockerfile.template", "Full path to the Dockerfile template.")
	baseImage = flag.String("base_image", "1.5.3", "Base image to run the server. Example: golang:1.5.3-alpine")
	port      = flag.Int("port", 8999, "Port to run/expose service. Example: '8999'")

	loadConcurrencyLevel = flag.Int("load_c", 10, "Number of concurrent workers generating load.")
	loadMaxQPS           = flag.Int("load_maxqps", 1000, "Maximum QPS impressed on the server. Zero means infinite.")
	loadOps              = flag.Int("load_ops", 1000, "Numer of operations done per request.")
	loadMem              = flag.Int("load_mem", 1024, "Amount of memory allocated per request (in bytes).")
	loadDuration         = flag.Duration("load_duration", 5*time.Second, "Duration of load. Example: 1m")
	loadPattern          = flag.String("load_pattern", "[{\"mem\":1024,\"ops\":1000}]", "JSON Patter of the load.Ex: \"[{\"mem\":\"1024\",\"ops\":\"1000\"}]\"")

	cpuset = flag.String("cpuset", "", "Number of CPUs used by the server container.")
	mem    = flag.String("memory", "1g", "Memory allocated in the container.")
)

type Run struct {
	Summary    *load.Summary
	ServerInfo *serverinfo.ServerInfo
}

type ExecutionSummary struct {
	ExprID string
	Runs   []Run
}

type loadConfig struct {
	Ops int    `json:"ops"`
	Mem string `json:"mem"`
}

func main() {
	flag.Parse()

	if *logToStderr {
		log.SetOutput(os.Stderr)
	}

	if *repeat < 1 {
		log.Fatalf("Invalid number of repetitions: %d", *repeat)
	}

	lc := []loadConfig{}
	if err := json.NewDecoder(bytes.NewBuffer([]byte(*loadPattern))).Decode(&lc); err != nil {
		log.Fatalf("[Stats.Inc] Problem decoding data: Pattern:%s Err:%q", *loadPattern, err)
	}
	var loadURLs []string
	for _, c := range lc {
		b, err := bytefmt.ToBytes(c.Mem)
		if err != nil {
			log.Fatalf("Invalid load specification. Mem:%d Spec:%+v", c.Mem, c)
		}
		loadURLs = append(loadURLs, serverURI(*port, fmt.Sprintf("work?cpu=%d&mem=%d", c.Ops, b)))
	}

	exprID := stripChars(*baseImage, ":.-")
	dockerfilePath, err := docker.CreateDockerfile(*dfTmpl, *baseImage, exprID, *port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Dockerfile created: ", dockerfilePath, "\n####\n")

	// Building image.
	iName := fmt.Sprintf("danielfireman/phd-experiments:%s", exprID)
	if err := docker.BuildImage(dockerfilePath, iName); err != nil {
		log.Fatal(err)
	}
	log.Println("Server image created: ", iName, "\n####\n")

	containerName := fmt.Sprintf("%s", exprID)
	// Check if container is running and stop, if it is the case.
	log.Printf("Checking if container %s is running.", containerName)
	if ok, err := docker.IsContainerRunning(containerName); err != nil {
		log.Fatal(err)
	} else if ok {
		log.Printf("Container %s is running, stopping it.", containerName)
		if err := docker.StopContainer(containerName); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("Container %s is not running, starting it.", containerName)
	}

	executionSumm := ExecutionSummary{
		ExprID: exprID,
	}
	for currentRepeat := 0; currentRepeat < *repeat; currentRepeat++ {
		log.Printf("Starting container %s", containerName)
		if err := docker.StartContainer(iName, containerName, serverURI(*port, "ping"), *port); err != nil {
			log.Fatal(err)
		}
		log.Printf("Container %s is up and server is running and healthy\n####\n", containerName)

		log.Printf("Starting server stress")
		// Load and fetching server information.
		resultSummary := load.Run(*loadMaxQPS, *loadConcurrencyLevel, *loadDuration, loadURLs...)
		serverInfo, err := serverinfo.Fetch(serverURI(*port, "debug/vars"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		executionSumm.Runs = append(executionSumm.Runs, Run{
			Summary:    resultSummary,
			ServerInfo: serverInfo,
		})
		log.Println("Finished stressing server\n####\n")

		// Stopping container.
		if err := docker.StopContainer(containerName); err != nil {
			log.Fatal(err)
		}
		log.Println("Container stopped: ", containerName, "\n####\n")
	}
	if err := os.Remove(dockerfilePath); err != nil {
		log.Fatal(err)
	}
	log.Println("Docker file successfully removed: ", dockerfilePath, "\n####\n")
	b, err := json.MarshalIndent(executionSumm, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func serverURI(port int, path string) string {
	return fmt.Sprintf("http://localhost:%d/%s", port, path)
}

func stripChars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}
