package commands

import (
	"fmt"
	resource "github.com/mdelillo/github-workflow-resource"
	"strconv"
)

type Check struct {
	githubClient GithubClient
}

func NewCheck(githubClient GithubClient) *Check {
	return &Check{
		githubClient: githubClient,
	}
}

func (c *Check) Execute(request resource.CheckRequest) (resource.CheckResponse, error) {
	workflowRuns, err := c.githubClient.GetWorkflowRuns(request.Source.Repo, request.Source.WorkflowID)
	if err != nil {
		return resource.CheckResponse{}, fmt.Errorf("failed to get workflow runs: %w", err)
	}

	var response resource.CheckResponse

	for _, workflowRun := range workflowRuns {
		response = append([]resource.Version{{ID: strconv.Itoa(workflowRun.ID)}}, response...)

		if request.Version.ID == strconv.Itoa(workflowRun.ID) {
			break
		}
	}

	return response, nil
}
