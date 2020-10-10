package commands

import (
	"encoding/json"
	"fmt"
	resource "github.com/mdelillo/github-workflow-resource"
	"io/ioutil"
	"path/filepath"
	"time"
)

type In struct {
	githubClient GithubClient
}

func NewIn(githubClient GithubClient) *In {
	return &In{
		githubClient: githubClient,
	}
}

func (i *In) Execute(request resource.InRequest, outputDir string) (resource.InResponse, error) {
	workflowRun, err := i.githubClient.GetWorkflowRun(request.Source.Repo, request.Version.ID)
	if err != nil {
		return resource.InResponse{}, fmt.Errorf("failed to get workflow run: %w", err)
	}

	metadata, err := json.Marshal(workflowRun)
	if err != nil {
		return resource.InResponse{}, fmt.Errorf("failed to marshal workflow run: %w", err)
	}

	metadataPath := filepath.Join(outputDir, "metadata.json")

	err = ioutil.WriteFile(metadataPath, metadata, 0644)
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
