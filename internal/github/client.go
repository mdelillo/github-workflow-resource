package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Client struct {
	endpoint string
	token    string
	client   *http.Client
}

type WorkflowRun struct {
	ID         int       `json:"id"`
	WorkflowID int       `json:"workflow_id"`
	Status     string    `json:"status"`
	Conclusion string    `json:"conclusion"`
	URL        string    `json:"url"`
	HtmlURL    string    `json:"html_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Option func(*Client)

func NewClient(token string, options ...Option) *Client {
	client := &Client{
		endpoint: "https://api.github.com",
		token:    token,
		client: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func WithEndpoint(endpoint string) func(*Client) {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

func (c *Client) GetWorkflowRuns(repo, workflowID string) ([]WorkflowRun, error) {
	url := fmt.Sprintf("%s/repos/%s/actions/workflows/%s/runs", c.endpoint, repo, workflowID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "token "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("got unsuccessful response from github: %d\n%s", resp.StatusCode, string(body))
	}

	var workflowRunsResp struct {
		WorkflowRuns []WorkflowRun `json:"workflow_runs"`
	}
	err = json.Unmarshal(body, &workflowRunsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	var workflowRuns []WorkflowRun
	for _, w := range workflowRunsResp.WorkflowRuns {
		workflowRuns = append(workflowRuns, w)
	}

	return workflowRuns, nil
}

func (c *Client) GetWorkflowRun(repo, runID string) (WorkflowRun, error) {
	url := fmt.Sprintf("%s/repos/%s/actions/runs/%s", c.endpoint, repo, runID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return WorkflowRun{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "token "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return WorkflowRun{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WorkflowRun{}, fmt.Errorf("failed to read body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return WorkflowRun{}, fmt.Errorf("got unsuccessful response from github: %d\n%s", resp.StatusCode, string(body))
	}

	var workflowRun WorkflowRun
	err = json.Unmarshal(body, &workflowRun)
	if err != nil {
		return WorkflowRun{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return workflowRun, nil
}
