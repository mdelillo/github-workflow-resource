package github_test

import (
	"github.com/mdelillo/github-workflow-resource/internal/github"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	assertpkg "github.com/stretchr/testify/assert"
	requirepkg "github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGithub(t *testing.T) {
	spec.Run(t, "Github Workflow Resource", testGithub, spec.Report(report.Terminal{}))
}

func testGithub(t *testing.T, context spec.G, it spec.S) {
	var (
		assert  = assertpkg.New(t)
		require = requirepkg.New(t)

		client *github.Client
		server *httptest.Server
	)

	it.Before(func() {
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/repos/some-repo/actions/workflows/123/runs" {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"message": "Not Found"}`))
			}

			if r.Header.Get("Authorization") != "token some-token" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"message": "Unauthorized"}`))
			}

			_, _ = w.Write([]byte(`
{
  "total_count": 4,
  "workflow_runs": [
    {
      "id": 4,
      "status": "queued",
      "conclusion": null,
      "workflow_id": 123,
      "url": "some-url-4",
      "html_url": "some-html-url-4",
      "created_at": "2020-01-04T00:00:00Z",
      "updated_at": "2020-01-04T01:00:00Z"
    },
    {
      "id": 3,
      "status": "in_progress",
      "conclusion": null,
      "workflow_id": 123,
      "url": "some-url-3",
      "html_url": "some-html-url-3",
      "created_at": "2020-01-03T00:00:00Z",
      "updated_at": "2020-01-03T01:00:00Z"
    },
    {
      "id": 2,
      "status": "completed",
      "conclusion": "success",
      "workflow_id": 123,
      "url": "some-url-2",
      "html_url": "some-html-url-2",
      "created_at": "2020-01-02T00:00:00Z",
      "updated_at": "2020-01-02T01:00:00Z"
    },
    {
      "id": 1,
      "status": "completed",
      "conclusion": "failure",
      "workflow_id": 123,
      "url": "some-url-1",
      "html_url": "some-html-url-1",
      "created_at": "2020-01-01T00:00:00Z",
      "updated_at": "2020-01-01T01:00:00Z"
    }
  ]
}
`))
		}))

		client = github.NewClient("some-token", github.WithEndpoint(server.URL))
	})

	context("GetWorkflowRuns", func() {
		it("returns all runs of the workflow", func() {
			workflowRuns, err := client.GetWorkflowRuns("some-repo", "123")
			require.NoError(err)

			assert.Equal([]github.WorkflowRun{
				{
					ID:         4,
					WorkflowID: 123,
					Status:     "queued",
					Conclusion: "",
					URL:        "some-url-4",
					HtmlURL:    "some-html-url-4",
					CreatedAt:  time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
					UpdatedAt:  time.Date(2020, 1, 4, 1, 0, 0, 0, time.UTC),
				},
				{
					ID:         3,
					WorkflowID: 123,
					Status:     "in_progress",
					Conclusion: "",
					URL:        "some-url-3",
					HtmlURL:    "some-html-url-3",
					CreatedAt:  time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					UpdatedAt:  time.Date(2020, 1, 3, 1, 0, 0, 0, time.UTC),
				},
				{
					ID:         2,
					WorkflowID: 123,
					Status:     "completed",
					Conclusion: "success",
					URL:        "some-url-2",
					HtmlURL:    "some-html-url-2",
					CreatedAt:  time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					UpdatedAt:  time.Date(2020, 1, 2, 1, 0, 0, 0, time.UTC),
				},
				{
					ID:         1,
					WorkflowID: 123,
					Status:     "completed",
					Conclusion: "failure",
					URL:        "some-url-1",
					HtmlURL:    "some-html-url-1",
					CreatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:  time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
				},
			}, workflowRuns)
		})
	})
}
