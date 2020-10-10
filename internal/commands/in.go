package commands

import (
	"encoding/json"
	"fmt"
	resource "github.com/mdelillo/github-workflow-resource"
	"github.com/mdelillo/github-workflow-resource/internal/github"
	"io"
	"time"
)

type In struct {
	githubClient GithubClient
	requestDelay time.Duration
}

type InOption func(*In)

func NewIn(githubClient GithubClient, options ...InOption) *In {
	in := &In{
		githubClient: githubClient,
		requestDelay: 10 * time.Second,
	}
	for _, option := range options {
		option(in)
	}
	return in
}

func WithRequestDelay(requestDelay time.Duration) func(in *In) {
	return func(i *In) {
		i.requestDelay = requestDelay
	}
}

func (i In) Execute(request resource.InRequest, metadataFile io.Writer) (resource.InResponse, error) {
	var workflowRun github.WorkflowRun

	for {
		var err error
		workflowRun, err = i.githubClient.GetWorkflowRun(request.Source.Repo, request.Version.ID)
		if err != nil {
			return resource.InResponse{}, fmt.Errorf("failed to get workflow run: %w", err)
		}

		if workflowRun.Status == "completed" {
			break
		}

		time.Sleep(i.requestDelay)
	}

	metadata, err := json.MarshalIndent(workflowRun, "", "  ")
	if err != nil {
		return resource.InResponse{}, fmt.Errorf("failed to marshal workflow run: %w", err)
	}

	_, err = metadataFile.Write(metadata)
	if err != nil {
		return resource.InResponse{}, fmt.Errorf("failed to write metadata: %w", err)
	}

	response := resource.InResponse{
		Version: request.Version,
		Metadata: []resource.Metadata{
			{Name: "status", Value: workflowRun.Status},
			{Name: "conclusion", Value: workflowRun.Conclusion},
			{Name: "url", Value: workflowRun.URL},
			{Name: "html_url", Value: workflowRun.HtmlURL},
			{Name: "created_at", Value: workflowRun.CreatedAt.Format(time.RFC3339)},
			{Name: "updated_at", Value: workflowRun.UpdatedAt.Format(time.RFC3339)},
		},
	}

	return response, nil
}
