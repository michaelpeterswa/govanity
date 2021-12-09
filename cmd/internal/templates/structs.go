package templates

import "github.com/michaelpeterswa/govanity/cmd/internal/github"

type Homepage struct {
	Title string
	Repos []github.CondensedRepository
}

type Repo struct {
	Repo github.CondensedRepository

	GoImport string
	GoSource string
}
