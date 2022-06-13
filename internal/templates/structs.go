package templates

import "github.com/michaelpeterswa/govanity/internal/github"

type Homepage struct {
	Title string
	Repos []github.CondensedRepository
}

type Repo struct {
	Title string

	Repo github.CondensedRepository

	GoImport string
	GoSource string
}
