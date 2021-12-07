package templates

type Homepage struct {
	Title string
}

type Repo struct {
	Title       string
	Name        string
	Description string

	GoImport string
	GoSource string
}
