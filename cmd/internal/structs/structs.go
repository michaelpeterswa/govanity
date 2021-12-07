package structs

type Settings struct {
	Username string   `yaml:"username"`
	Repos    []string `yaml:"repos"`
}
