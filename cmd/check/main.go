package main

import (
	"encoding/json"
	"fmt"
	"github.com/mdelillo/github-workflow-resource/internal/github"
	"io/ioutil"
	"os"
	"strconv"
)

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type CheckResponse []Version

type Source struct {
	Repo        string `json:"repo"`
	WorkflowID  string `json:"workflow_id"`
	GithubToken string `json:"github_token"`
}

type Version struct {
	ID string `json:"id"`
}

func main() {
	output, err := check()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(output)
}

func check() (string, error) {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read stdin: %w", err)
	}

	var request CheckRequest
	err = json.Unmarshal(stdin, &request)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal input: %w\n", err)
	}

	githubClient := github.NewClient(request.Source.GithubToken)

	workflowRuns, err := githubClient.GetWorkflowRuns(request.Source.Repo, request.Source.WorkflowID)
	if err != nil {
		return "", fmt.Errorf("failed to get workflow runs: %w\n", err)
	}

	var response CheckResponse

	returnVersions := true
	if request.Version.ID != "" {
		returnVersions = false
	}

	for _, workflowRun := range workflowRuns {
		if request.Version.ID == strconv.Itoa(workflowRun.ID) {
			returnVersions = true
		}

		if returnVersions {
			response = append(response, Version{ID: strconv.Itoa(workflowRun.ID)})
		}
	}

	output, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(output), nil
}
