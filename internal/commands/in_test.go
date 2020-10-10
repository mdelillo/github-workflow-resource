package commands

import (
	"bytes"
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
		in               *In
		fakeGithubClient FakeInGithubClient

		assert  = assertpkg.New(t)
		require = requirepkg.New(t)
	)

	it.Before(func() {
		fakeGithubClient = FakeInGithubClient{}
		in = NewIn(&fakeGithubClient, WithRequestDelay(0))
	})

	it("returns the workflow run with metadata", func() {
		var buffer bytes.Buffer

		request := resource.InRequest{Version: resource.Version{ID: "123"}}
		response, err := in.Execute(request, &buffer)
		require.NoError(err)

		assert.Equal(resource.InResponse{
			Version: resource.Version{
				ID: "123",
			},
			Metadata: []resource.Metadata{
				{Name: "status", Value: "completed"},
				{Name: "conclusion", Value: "some-conclusion"},
				{Name: "url", Value: "some-url"},
				{Name: "html_url", Value: "some-html-url"},
				{Name: "created_at", Value: "2020-01-01T00:00:00Z"},
				{Name: "updated_at", Value: "2020-01-01T00:01:00Z"},
			},
		}, response)

		assert.JSONEq(`{
  "id": 123,
  "status": "completed",
  "conclusion": "some-conclusion",
  "workflow_id": 456,
  "url": "some-url",
  "html_url": "some-html-url",
  "created_at": "2020-01-01T00:00:00Z",
  "updated_at": "2020-01-01T00:01:00Z"
}`, buffer.String())
	})

	context("when wait_for_completion is true", func() {
		it("waits until the run is complete before returning", func() {
			fakeGithubClient.inProgressRunRequests = 2

			var buffer bytes.Buffer

			request := resource.InRequest{
				Params: resource.InParams{
					WaitForCompletion: true,
				},
				Version: resource.Version{ID: "123"},
			}
			response, err := in.Execute(request, &buffer)
			require.NoError(err)

			assert.Equal(resource.InResponse{
				Version: resource.Version{
					ID: "123",
				},
				Metadata: []resource.Metadata{
					{Name: "status", Value: "completed"},
					{Name: "conclusion", Value: "some-conclusion"},
					{Name: "url", Value: "some-url"},
					{Name: "html_url", Value: "some-html-url"},
					{Name: "created_at", Value: "2020-01-01T00:00:00Z"},
					{Name: "updated_at", Value: "2020-01-01T00:01:00Z"},
				},
			}, response)
		})
	})
}

type FakeInGithubClient struct {
	inProgressRunRequests int
}

func (f *FakeInGithubClient) GetWorkflowRuns(repo, workflowId string) ([]github.WorkflowRun, error) {
	panic("implement me")
}

func (f *FakeInGithubClient) GetWorkflowRun(repo, workflowRunId string) (github.WorkflowRun, error) {
	if f.inProgressRunRequests > 0 {
		f.inProgressRunRequests--

		return github.WorkflowRun{
			Status: "in_progress",
		}, nil
	}

	return github.WorkflowRun{
		ID:         123,
		WorkflowID: 456,
		Status:     "completed",
		Conclusion: "some-conclusion",
		URL:        "some-url",
		HtmlURL:    "some-html-url",
		CreatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:  time.Date(2020, 1, 1, 0, 1, 0, 0, time.UTC),
	}, nil
}
