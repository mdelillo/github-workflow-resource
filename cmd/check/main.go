package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type checkInput struct {
	Source struct {
		Repo        string `json:"repo"`
		WorkflowID  string `json:"workflow_id"`
		GithubToken string `json:"github_token"`
	} `json:"source"`
}

type workflowRunsResponse struct {
	WorkflowRuns []struct {
		ID int `json:"id"`
	} `json:"workflow_runs"`
}

type workflowRun struct {
	ID string `json:"id"`
}

func main() {
	output, err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(output)
}

func run() (string, error) {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read stdin: %w", err)
	}

	var input checkInput
	err = json.Unmarshal(stdin, &input)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal input: %w\n", err)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/actions/workflows/%s/runs", input.Source.Repo, input.Source.WorkflowID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "token "+input.Source.GithubToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("got unsuccessful response from github: %d\n%s", resp.StatusCode, string(body))
	}

	var workflowRunsResp workflowRunsResponse
	err = json.Unmarshal(body, &workflowRunsResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	var workflowRuns []workflowRun
	for _, w := range workflowRunsResp.WorkflowRuns {
		workflowRuns = append([]workflowRun{{ID: strconv.Itoa(w.ID)}}, workflowRuns...)
	}

	output, err := json.Marshal(workflowRuns)
	if err != nil {
		return "", fmt.Errorf("failed to marshal workflow runs: %w", err)
	}

	return string(output), nil
}
