package resource

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type Source struct {
	Repo        string `json:"repo"`
	WorkflowID  string `json:"workflow_id"`
	GithubToken string `json:"github_token"`
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
