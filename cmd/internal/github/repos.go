package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/michaelpeterswa/govanity/cmd/internal/structs"
	"github.com/michaelpeterswa/govanity/cmd/internal/utils"
)

type GitHubRepositories struct {
	mu           sync.Mutex
	Repositories []GitHubRepository
}

type CondensedRepository struct {
	Name            string   `json:"name"`
	FullName        string   `json:"full_name"`
	Owner           COwner   `json:"owner"`
	HTMLURL         string   `json:"html_url"`
	Description     string   `json:"description"`
	Fork            bool     `json:"fork"`
	StargazersCount int      `json:"stargazers_count"`
	License         CLicense `json:"license"`
}

type COwner struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
}

type CLicense struct {
	Name string `json:"name"`
}

type GitHubRepository struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	HTMLURL          string      `json:"html_url"`
	Description      string      `json:"description"`
	Fork             bool        `json:"fork"`
	URL              string      `json:"url"`
	ForksURL         string      `json:"forks_url"`
	KeysURL          string      `json:"keys_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	TeamsURL         string      `json:"teams_url"`
	HooksURL         string      `json:"hooks_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	EventsURL        string      `json:"events_url"`
	AssigneesURL     string      `json:"assignees_url"`
	BranchesURL      string      `json:"branches_url"`
	TagsURL          string      `json:"tags_url"`
	BlobsURL         string      `json:"blobs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	TreesURL         string      `json:"trees_url"`
	StatusesURL      string      `json:"statuses_url"`
	LanguagesURL     string      `json:"languages_url"`
	StargazersURL    string      `json:"stargazers_url"`
	ContributorsURL  string      `json:"contributors_url"`
	SubscribersURL   string      `json:"subscribers_url"`
	SubscriptionURL  string      `json:"subscription_url"`
	CommitsURL       string      `json:"commits_url"`
	GitCommitsURL    string      `json:"git_commits_url"`
	CommentsURL      string      `json:"comments_url"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	ContentsURL      string      `json:"contents_url"`
	CompareURL       string      `json:"compare_url"`
	MergesURL        string      `json:"merges_url"`
	ArchiveURL       string      `json:"archive_url"`
	DownloadsURL     string      `json:"downloads_url"`
	IssuesURL        string      `json:"issues_url"`
	PullsURL         string      `json:"pulls_url"`
	MilestonesURL    string      `json:"milestones_url"`
	NotificationsURL string      `json:"notifications_url"`
	LabelsURL        string      `json:"labels_url"`
	ReleasesURL      string      `json:"releases_url"`
	DeploymentsURL   string      `json:"deployments_url"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	PushedAt         time.Time   `json:"pushed_at"`
	GitURL           string      `json:"git_url"`
	SSHURL           string      `json:"ssh_url"`
	CloneURL         string      `json:"clone_url"`
	SvnURL           string      `json:"svn_url"`
	Homepage         interface{} `json:"homepage"`
	Size             int         `json:"size"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Language         string      `json:"language"`
	HasIssues        bool        `json:"has_issues"`
	HasProjects      bool        `json:"has_projects"`
	HasDownloads     bool        `json:"has_downloads"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	ForksCount       int         `json:"forks_count"`
	MirrorURL        interface{} `json:"mirror_url"`
	Archived         bool        `json:"archived"`
	Disabled         bool        `json:"disabled"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	License          struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		SpdxID string `json:"spdx_id"`
		URL    string `json:"url"`
		NodeID string `json:"node_id"`
	} `json:"license"`
	AllowForking  bool          `json:"allow_forking"`
	IsTemplate    bool          `json:"is_template"`
	Topics        []interface{} `json:"topics"`
	Visibility    string        `json:"visibility"`
	Forks         int           `json:"forks"`
	OpenIssues    int           `json:"open_issues"`
	Watchers      int           `json:"watchers"`
	DefaultBranch string        `json:"default_branch"`
}

func (gh *GitHubRepositories) GetGitHubRepositories(ctx context.Context, s structs.Settings) error {
	atEnd := false
	gh.mu.Lock()

	nextRegex := regexp.MustCompile(`^<.+>; rel="next"`)

	url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=40", s.Username)

	for !atEnd {
		timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
		defer cancel()

		req, err := http.NewRequestWithContext(timeoutCtx, "GET", url, nil)
		if err == context.DeadlineExceeded {
			return errors.New("timed out")
		} else if err != nil {
			return err
		}

		c := &http.Client{}

		resp, err := c.Do(req)
		if err != nil {
			return err
		}

		link := resp.Header["Link"]
		if len(link) == 0 {
			return errors.New("link header is empty")
		}
		linkValue := link[0]

		links := strings.Split(linkValue, ",")

		found := false

		for _, l := range links {
			l = strings.TrimSpace(l)
			match := nextRegex.Match([]byte(l))
			if match {
				found = true
				bracketedURL := nextRegex.FindString(l)
				url = strings.TrimPrefix(bracketedURL, "<")
				url = strings.TrimSuffix(url, `>; rel="next"`)
				break
			}
		}
		if !found {
			atEnd = true
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var tmp []GitHubRepository

		err = json.Unmarshal(body, &tmp)
		if err != nil {
			return err
		}

		gh.Repositories = append(gh.Repositories, tmp...)
		defer resp.Body.Close()
	}

	defer gh.mu.Unlock()
	return nil
}

func (gh *GitHubRepositories) GetCurrentGoRepositories() []CondensedRepository {
	gh.mu.Lock()
	var cr []CondensedRepository
	for _, repo := range gh.Repositories {
		if repo.Language == "Go" {
			cr = append(cr, repo.CondenseGitHubRepository())
		}
	}
	defer gh.mu.Unlock()

	return cr
}

func (gh *GitHubRepository) CondenseGitHubRepository() CondensedRepository {
	return CondensedRepository{
		Name:     gh.Name,
		FullName: gh.FullName,
		Owner: COwner{
			Login:     gh.Owner.Login,
			AvatarURL: gh.Owner.AvatarURL,
			HTMLURL:   gh.Owner.HTMLURL,
		},
		HTMLURL:         gh.HTMLURL,
		Description:     gh.Description,
		Fork:            gh.Fork,
		StargazersCount: gh.StargazersCount,
		License: CLicense{
			Name: gh.License.Name,
		},
	}
}

func (gh *GitHubRepositories) GetValidRepositories(settings structs.Settings) map[string]CondensedRepository {
	existing := gh.GetCurrentGoRepositories()

	matches := make(map[string]CondensedRepository)
	for _, repo := range existing {
		if utils.Contains(settings.Repos, repo.Name) {
			matches[repo.Name] = repo
		}
	}

	return matches
}
