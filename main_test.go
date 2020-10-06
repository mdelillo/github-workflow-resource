package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	assertpkg "github.com/stretchr/testify/assert"
	requirepkg "github.com/stretchr/testify/require"
)

const (
	repo       = "mdelillo/github-workflow-resource"
	workflowID = "2743569"
)

func TestGithubWorkflowResource(t *testing.T) {
	if _, ok := os.LookupEnv("GITHUB_TOKEN"); !ok {
		t.Fatal("Must set GITHUB_TOKEN")
	}

	spec.Run(t, "Github Workflow Resource", testGithubWorkflowResource, spec.Report(report.Terminal{}))
}

func testGithubWorkflowResource(t *testing.T, context spec.G, it spec.S) {
	var (
		assert  = assertpkg.New(t)
		require = requirepkg.New(t)

		checkPath string
	)

	it.Before(func() {
		tempFile, err := ioutil.TempFile("", "github-workflow-resource")
		require.NoError(err)
		checkPath = tempFile.Name()
		require.NoError(tempFile.Close())

		output, err := exec.Command("go", "build", "-o", checkPath, "./cmd/check").CombinedOutput()
		require.NoError(err, string(output))
	})

	it.After(func() {
		_ = os.Remove(checkPath)
	})

	context("check", func() {
		it("lists all runs of a github workflow", func() {
			cmd := exec.Command(checkPath)
			sourceParams := fmt.Sprintf(`{"source": {"repo": "%s", "workflow_id": "%s", "github_token": "%s"}}`,
				repo,
				workflowID,
				os.Getenv("GITHUB_TOKEN"),
			)
			cmd.Stdin = strings.NewReader(sourceParams)

			output, err := cmd.CombinedOutput()
			require.NoError(err, string(output))

			assert.JSONEq(`[
  {"id": "275207818"},
  {"id": "275208078"},
  {"id": "275208338"},
  {"id": "275208364"}
]`, string(output))
		})
	})
}
