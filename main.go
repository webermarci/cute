package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/akamensky/argparse"
	"github.com/alitto/pond"
)

type CommandGroup struct {
	Commands []Command
}

type Command struct {
	Command string
	Alias   string
}

type Result struct {
	Command  string
	Alias    string
	Output   string
	Duration time.Duration
}

func parseGroups(input string) map[string]string {
	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}

	items := strings.FieldsFunc(input, f)

	m := make(map[string]string)
	for _, item := range items {
		x := strings.Split(item, "=")
		m[x[0]] = strings.Trim(x[1], "\"")
	}

	return m
}

func execute(shell string, command Command) Result {
	start := time.Now()

	out, err := exec.Command(shell, "-c", command.Command).Output()
	if err != nil {
		log.Fatal(err)
	}

	return Result{
		Command:  command.Command,
		Alias:    command.Alias,
		Output:   strings.TrimSpace(string(out)),
		Duration: time.Since(start),
	}
}

func main() {
	description := `You need to pipe the commands into it:
	    "commands.txt | cute" or "cute < commands.txt"
	    
	    The commands.txt file should contain lines of command groups like so:
	    first="echo hello" second="echo world"
	    thrid="echo bye"

	    Command groups may contain one or more commands. Command groups are executed in parallel, but commands in a group are sync.
	`

	parser := argparse.NewParser("cute", description)
	version := parser.Flag("V", "version",
		&argparse.Options{Help: "Print version"})
	verbose := parser.Flag("v", "verbose",
		&argparse.Options{Help: "Verbose mode"})
	n := parser.Int("n", "parallel-tasks",
		&argparse.Options{Default: runtime.NumCPU(), Help: "Number of parallel tasks"})
	shell := parser.String("s", "shell",
		&argparse.Options{Default: "bash", Help: "Shell to execute the commands with"})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	if *version {
		fmt.Println("cute - v0.3")
		return
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		log.Fatal("Nothing is piped into")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	var pipe []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		pipe = append(pipe, input)
	}

	if len(pipe) == 0 {
		log.Fatal("Nothing is piped into")
	}

	var groups []CommandGroup
	commandCounter := 0

	for _, groupString := range strings.Split(string(pipe), "\n") {
		group := CommandGroup{}
		parsed := parseGroups(groupString)
		for alias, command := range parsed {
			group.Commands = append(group.Commands, Command{
				Alias:   alias,
				Command: command,
			})
			commandCounter++
		}
		groups = append(groups, group)
	}

	if len(groups) == 0 {
		log.Fatal("Failed to parse any command groups")
	}

	pool := pond.New(*n, commandCounter)
	groupResults := make([]map[string]any, len(groups))
	start := time.Now()

	for i, g := range groups {
		index := i
		group := g

		pool.Submit(func() {
			results := make([]Result, len(group.Commands))
			for j, cmd := range group.Commands {
				results[j] = execute(*shell, cmd)

				if *verbose {
					log.Printf("%s=\"%s\" returned %s in %s",
						cmd.Alias, cmd.Command, results[j].Output, results[j].Duration)
				}
			}

			groupResult := make(map[string]any)
			for _, result := range results {
				parsedFloat, err := strconv.ParseFloat(result.Output, 64)
				if err == nil {
					groupResult[result.Alias] = parsedFloat
				} else {
					groupResult[result.Alias] = result.Output
				}

				groupResult[result.Alias+"_duration"] = int64(result.Duration)
			}
			groupResults[index] = groupResult
		})
	}

	pool.StopAndWait()

	if *verbose {
		log.Printf("All finished in %s\n", time.Since(start))
	}

	for _, result := range groupResults {
		data, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(data))
	}
}
