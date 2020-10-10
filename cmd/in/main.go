package main

import (
	"encoding/json"
	"fmt"
	"github.com/mdelillo/github-workflow-resource"
	"github.com/mdelillo/github-workflow-resource/internal/commands"
	"github.com/mdelillo/github-workflow-resource/internal/github"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fatal("failed to read stdin", err)
	}

	var request resource.InRequest
	err = json.Unmarshal(stdin, &request)
	if err != nil {
		fatal("failed to unmarshal input", err)
	}

	metadataFile, err := os.OpenFile(filepath.Join(os.Args[1], "metadata.json"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fatal("failed to open metadata.json", err)
	}
	defer metadataFile.Close()

	in := commands.NewIn(github.NewClient(request.Source.GithubToken))
	response, err := in.Execute(request, metadataFile)
	if err != nil {
		fatal("failed to execute in", err)
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
