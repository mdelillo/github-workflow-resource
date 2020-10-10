package commands

import (
	"fmt"
	resource "github.com/mdelillo/github-workflow-resource"
	"github.com/mdelillo/github-workflow-resource/internal/github"
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
		if !c.workflowRunMatches(workflowRun, request.Source) {
			continue
		}

		response = append([]resource.Version{{ID: strconv.Itoa(workflowRun.ID)}}, response...)

		if request.Version.ID == strconv.Itoa(workflowRun.ID) {
			break
		}
	}

	return response, nil
}

func (c *Check) workflowRunMatches(workflowRun github.WorkflowRun, source resource.Source) bool {
	if source.Status != "" && workflowRun.Status != source.Status {
		return false
	}

	if source.Conclusion != "" && workflowRun.Conclusion != source.Conclusion {
		return false
	}

	return true
}
