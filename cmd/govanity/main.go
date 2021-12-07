package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/michaelpeterswa/govanity/cmd/internal/github"
	"github.com/michaelpeterswa/govanity/cmd/internal/settings"
	"github.com/michaelpeterswa/govanity/cmd/internal/structs"
	"github.com/michaelpeterswa/govanity/cmd/internal/templates"
	"go.uber.org/zap"
)

type Handler struct {
	gh        *github.GitHubRepositories
	settings  structs.Settings
	valids    map[string]github.CondensedRepository
	templates map[string]*template.Template
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("govanity is initializing...")

	background := context.Background()

	// TODO: set environment variable
	settings, err := settings.InitSettings()
	if err != nil {
		logger.Error("init settings failed", zap.Error(err))
	}

	var h Handler

	h.templates = make(map[string]*template.Template)
	homeTmpl, err := template.ParseFiles("cmd/internal/templates/home.gotmpl")
	if err != nil {
		logger.Error("error parsing home template", zap.Error(err))
	}
	h.templates["home"] = homeTmpl

	repoTmpl, err := template.ParseFiles("cmd/internal/templates/repo.gotmpl")
	if err != nil {
		logger.Error("error parsing repo template", zap.Error(err))
	}
	h.templates["repo"] = repoTmpl

	h.settings = *settings
	h.gh = &github.GitHubRepositories{}

	err = h.gh.GetGitHubRepositories(background, *settings)
	if err != nil {
		logger.Error("getGitHubRepositories failed", zap.Error(err))
	}

	h.valids = h.gh.GetValidRepositories(*settings)

	r := mux.NewRouter()
	r.HandleFunc("/", h.HomeHandler)
	r.HandleFunc("/{repo}", h.VanityHandler)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("cmd/internal/static"))))

	http.ListenAndServe("localhost:8080", nil)
}

func (h Handler) HomeHandler(writer http.ResponseWriter, request *http.Request) {
	homeData := templates.Homepage{
		Title: "nw.codes",
	}
	h.templates["home"].Execute(writer, homeData)
}

func (h Handler) VanityHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	repo := vars["repo"]

	if v, ok := h.valids[repo]; ok {
		repoData := templates.Repo{
			Title:       v.FullName,
			Name:        v.Name,
			Description: v.Description,
			GoImport:    "importlol",
			GoSource:    "sourcelol",
		}

		h.templates["repo"].Execute(writer, repoData)
		return
	}
	writer.WriteHeader(http.StatusNotFound)
	fmt.Fprint(writer, `{"found":false}`)
	return
}