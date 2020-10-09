package commands

import (
	resource "github.com/mdelillo/github-workflow-resource"
	"github.com/mdelillo/github-workflow-resource/internal/github"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	assertpkg "github.com/stretchr/testify/assert"
	requirepkg "github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestIn(t *testing.T) {
	spec.Run(t, "In", testIn, spec.Report(report.Terminal{}))
}

func testIn(t *testing.T, context spec.G, it spec.S) {
	var (
		in *In

		assert  = assertpkg.New(t)
		require = requirepkg.New(t)
	)

	it.Before(func() {
		in = NewIn(FakeInGithubClient{})
	})

	it("returns the workflow run with metadata", func() {
		request := resource.InRequest{Version: resource.Version{ID: "123"}}
		response, err := in.Execute(request)
		require.NoError(err)

		assert.Equal(resource.InResponse{
			Version: resource.Version{
				ID: "123",
			},
			Metadata: []resource.Metadata{
				{Name: "status", Value: "some-status"},
				{Name: "conclusion", Value: "some-conclusion"},
				{Name: "url", Value: "some-url"},
				{Name: "html_url", Value: "some-html-url"},
				{Name: "created_at", Value: time.Unix(1, 0).Format(time.RFC3339)},
				{Name: "updated_at", Value: time.Unix(2, 0).Format(time.RFC3339)},
			},
		}, response)
	})
}

type FakeInGithubClient struct{}

func (FakeInGithubClient) GetWorkflowRuns(repo, workflowId string) ([]github.WorkflowRun, error) {
	panic("implement me")
}

func (FakeInGithubClient) GetWorkflowRun(repo, workflowRunId string) (github.WorkflowRun, error) {
	return github.WorkflowRun{
		ID:         123,
		Status:     "some-status",
		Conclusion: "some-conclusion",
		URL:        "some-url",
		HtmlURL:    "some-html-url",
		CreatedAt:  time.Unix(1, 0),
		UpdatedAt:  time.Unix(2, 0),
	}, nil
}
