package main

import (
	"encoding/json"
	"fmt"
	"github.com/mdelillo/github-workflow-resource"
	"github.com/mdelillo/github-workflow-resource/internal/commands"
	"github.com/mdelillo/github-workflow-resource/internal/github"
	"io/ioutil"
	"os"
)

func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fatal("failed to read stdin", err)
	}

	var request resource.CheckRequest
	err = json.Unmarshal(stdin, &request)
	if err != nil {
		fatal("failed to unmarshal input", err)
	}

	check := commands.NewCheck(github.NewClient(request.Source.GithubToken))

	response, err := check.Execute(request)
	if err != nil {
		fatal("failed to execute check", err)
	}

	output, err := json.Marshal(response)
	if err != nil {
		fatal("failed to marshal response", err)
	}

	fmt.Println(string(output))
}

func fatal(message string, err error) {
	fmt.Fprintf(os.Stderr, "Error: %s: %s\n", message, err.Error())
	os.Exit(1)
}
