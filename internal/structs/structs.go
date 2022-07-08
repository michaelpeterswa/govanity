package structs

type Settings struct {
	Username  string   `yaml:"username"`
	GithubPAT string   `yaml:"github-pat"`
	Domain    string   `yaml:"domain"`
	Repos     []string `yaml:"repos"`
}

type Healthcheck struct {
	IsHealthy bool `json:"is_healthy"`
}
