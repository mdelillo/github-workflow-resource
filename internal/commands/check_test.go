package commands

import (
	resource "github.com/mdelillo/github-workflow-resource"
	"github.com/mdelillo/github-workflow-resource/internal/github"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	assertpkg "github.com/stretchr/testify/assert"
	requirepkg "github.com/stretchr/testify/require"
	"testing"
)

func TestCheck(t *testing.T) {
	spec.Run(t, "Check", testCheck, spec.Report(report.Terminal{}))
}

func testCheck(t *testing.T, context spec.G, it spec.S) {
	var (
		check *Check

		assert  = assertpkg.New(t)
		require = requirepkg.New(t)
	)

	it.Before(func() {
		check = NewCheck(FakeCheckGithubClient{})
	})

	context("when the request specifies a version", func() {
		it("returns all versions since the given version", func() {
			request := resource.CheckRequest{Version: resource.Version{ID: "3"}}
			response, err := check.Execute(request)
			require.NoError(err)

			assert.Equal(resource.CheckResponse{
				{ID: "3"},
				{ID: "4"},
				{ID: "5"},
			}, response)
		})
	})

	context("when the request does not specify a version", func() {
		it("returns all versions", func() {
			request := resource.CheckRequest{}
			response, err := check.Execute(request)
			require.NoError(err)

			assert.Equal(resource.CheckResponse{
				{ID: "1"},
				{ID: "2"},
				{ID: "3"},
				{ID: "4"},
				{ID: "5"},
			}, response)
		})
	})
}

type FakeCheckGithubClient struct{}

func (FakeCheckGithubClient) GetWorkflowRuns(repo, workflowId string) ([]github.WorkflowRun, error) {
	return []github.WorkflowRun{
		{ID: 5},
		{ID: 4},
		{ID: 3},
		{ID: 2},
		{ID: 1},
	}, nil
}

func (FakeCheckGithubClient) GetWorkflowRun(repo, workflowRunId string) (github.WorkflowRun, error) {
	panic("implement me")
}
