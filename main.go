package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	"github.com/alitto/pond"
)

type Result struct {
	Command  string        `json:"command"`
	Output   string        `json:"output"`
	Duration time.Duration `json:"duration"`
}

func execute(shell string, command string) Result {
	start := time.Now()

	out, err := exec.Command(shell, "-c", command).Output()
	if err != nil {
		log.Fatal(err)
	}

	return Result{
		Command:  command,
		Output:   strings.TrimSpace(string(out)),
		Duration: time.Since(start),
	}
}

func main() {
	description := `cute -f commands.txt
	    
	    The commands.txt file should contain lines commands like so:
	    echo hello
	    echo world
	    echo bye
	`

	parser := argparse.NewParser("cute", description)
	version := parser.Flag("V", "version",
		&argparse.Options{Help: "Print version information"})
	verbose := parser.Flag("v", "verbose",
		&argparse.Options{Help: "Verbose mode"})
	n := parser.Int("n", "parallel-tasks",
		&argparse.Options{Default: runtime.NumCPU(), Help: "Number of parallel tasks"})
	file := parser.File("f", "file", os.O_RDONLY, 0400,
		&argparse.Options{Help: "File"})
	shell := parser.String("s", "shell",
		&argparse.Options{Default: "bash", Help: "Shell to execute the commands with"})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	if *version {
		fmt.Println("cute - v1.0")
		return
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	commandCounter := 0
	var commands []string = []string{}

	for fileScanner.Scan() {
		commandCounter++
		commands = append(commands, fileScanner.Text())
	}

	file.Close()

	pool := pond.New(*n, commandCounter)
	start := time.Now()

	results := make([]Result, commandCounter)

	for i, command := range commands {
		index := i
		cmd := command

		pool.Submit(func() {
			results[index] = execute(*shell, cmd)

			if *verbose {
				log.Printf("%s returned %s in %s", cmd, results[index].Output, results[index].Duration)
			}
		})
	}

	pool.StopAndWait()

	if *verbose {
		log.Printf("All finished in %s\n", time.Since(start))
	}

	for _, result := range results {
		data, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(data))
	}
}
