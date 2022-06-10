package structs

type Settings struct {
	Username  string   `yaml:"username"`
	Protocol  string   `yaml:"protocol"`
	GithubPAT string   `yaml:"github-pat"`
	Domain    string   `yaml:"domain"`
	Repos     []string `yaml:"repos"`
}
