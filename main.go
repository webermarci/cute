package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Result struct {
	Command  string        `json:"command"`
	Output   string        `json:"output"`
	Duration time.Duration `json:"duration"`
}

var (
	verbose  bool
	parallel bool
	format   string
	shell    string
	commands []string
)

func parseArgs() {
	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.BoolVar(&parallel, "p", false, "Execute the commands in parallel")
	flag.StringVar(&format, "f", "json", "Format of the output")
	flag.StringVar(&shell, "s", "sh", "Shell to use")
	flag.Parse()
	commands = flag.Args()

	switch format {
	case "json", "csv", "logfmt":
	default:
		log.Fatal("Invalid format type")
	}
}

func execute(command string) Result {
	start := time.Now()

	out, err := exec.Command(shell, "-c", command).Output()
	if err != nil {
		log.Fatal(err)
	}

	result := Result{
		Command:  command,
		Output:   strings.TrimSpace(string(out)),
		Duration: time.Since(start),
	}

	if verbose {
		log.Printf("\"%s\" returned %s in %s", command, result.Output, result.Duration)
	}

	return result
}

func main() {
	start := time.Now()

	parseArgs()
	var results []Result = make([]Result, len(commands))

	if parallel {
		wg := sync.WaitGroup{}
		wg.Add(len(commands))
		for i, command := range commands {
			go func(i int, cmd string) {
				results[i] = execute(cmd)
				wg.Done()
			}(i, command)
		}
		wg.Wait()
	} else {
		for i, command := range commands {
			results[i] = execute(command)
		}
	}

	if verbose {
		log.Printf("All finished in %s\n", time.Since(start))
	}

	switch format {
	case "json":
		for _, result := range results {
			bytes, err := json.Marshal(result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(bytes))
		}

	case "csv":
		fmt.Println("command,output,duration")
		for _, result := range results {
			fmt.Printf("%s,%s,%d\n", result.Command, result.Output, result.Duration)
		}

	case "logfmt":
		for _, result := range results {
			fmt.Printf("command=\"%s\" output=\"%s\" duration=%d\n", result.Command, result.Output, result.Duration)
		}
	}
}
