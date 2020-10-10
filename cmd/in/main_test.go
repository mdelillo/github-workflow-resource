package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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

func TestIn(t *testing.T) {
	if _, ok := os.LookupEnv("GITHUB_TOKEN"); !ok {
		t.Fatal("Must set GITHUB_TOKEN")
	}

	spec.Run(t, "Github Workflow Resource In", testIn, spec.Report(report.Terminal{}))
}

func testIn(t *testing.T, context spec.G, it spec.S) {
	var (
		assert  = assertpkg.New(t)
		require = requirepkg.New(t)

		inPath    string
		outputDir string
	)

	it.Before(func() {
		var err error
		outputDir, err = ioutil.TempDir("", "github-workflow-resource-in")
		require.NoError(err)

		tempFile, err := ioutil.TempFile("", "github-workflow-resource-in")
		require.NoError(err)
		inPath = tempFile.Name()
		require.NoError(tempFile.Close())

		output, err := exec.Command("go", "build", "-o", inPath).CombinedOutput()
		require.NoError(err, string(output))
	})

	it.After(func() {
		_ = os.Remove(inPath)
		_ = os.RemoveAll(outputDir)
	})

	it("writes workflow run metadata to metadata.json and stdout", func() {
		cmd := exec.Command(inPath, outputDir)
		sourceParams := fmt.Sprintf(`{"source": {"repo": "%s", "workflow_id": "%s", "github_token": "%s"}, "version": {"id": "%s"}}`,
			repo,
			workflowID,
			os.Getenv("GITHUB_TOKEN"),
			"275208338",
		)
		cmd.Stdin = strings.NewReader(sourceParams)

		output, err := cmd.CombinedOutput()
		require.NoError(err, string(output))

		assert.JSONEq(`{
  "version": {
    "id": "275208338"
  },
  "metadata": [
    {"name": "status",      "value": "completed"},
    {"name": "conclusion",  "value": "success"},
    {"name": "url",         "value": "https://api.github.com/repos/mdelillo/github-workflow-resource/actions/runs/275208338"},
    {"name": "html_url",    "value": "https://github.com/mdelillo/github-workflow-resource/actions/runs/275208338"},
    {"name": "created_at",  "value": "2020-09-27T13:57:52Z"},
    {"name": "updated_at",  "value": "2020-09-27T13:58:05Z"}
  ]
}`, string(output))

		metadata, err := ioutil.ReadFile(filepath.Join(outputDir, "metadata.json"))
		require.NoError(err)

		assert.JSONEq(`{
  "id": 275208338,
  "status": "completed",
  "conclusion": "success",
  "workflow_id": 2743569,
  "url": "https://api.github.com/repos/mdelillo/github-workflow-resource/actions/runs/275208338",
  "html_url": "https://github.com/mdelillo/github-workflow-resource/actions/runs/275208338",
  "created_at": "2020-09-27T13:57:52Z",
  "updated_at": "2020-09-27T13:58:05Z"
}`, string(metadata))
	})
}
