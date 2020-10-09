package commands

import "github.com/mdelillo/github-workflow-resource/internal/github"

type GithubClient interface {
	GetWorkflowRuns(repo, workflowId string) ([]github.WorkflowRun, error)
	GetWorkflowRun(repo, workflowRunId string) (github.WorkflowRun, error)
}
