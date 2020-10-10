package resource

type InRequest struct {
	Source  Source   `json:"source"`
	Version Version  `json:"version"`
	Params  InParams `json:"params"`
}

type InResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type InParams struct {
	WaitForCompletion bool `json:"wait_for_completion"`
}

type Source struct {
	Repo        string `json:"repo"`
	WorkflowID  string `json:"workflow_id"`
	GithubToken string `json:"github_token"`
	Status      string `json:"status"`
	Conclusion  string `json:"conclusion"`
}

type Version struct {
	ID string `json:"id"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type CheckResponse []Version
