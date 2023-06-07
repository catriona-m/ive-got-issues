package gh

import (
	"net/http"

	"github.com/google/go-github/v52/github"
	"github.com/gregjones/httpcache"
	"golang.org/x/oauth2"
)

type Token struct {
	Token *string
}

type Repo struct {
	Owner string
	Name  string
	Token
}

func NewRepo(owner, repo, token string) Repo {
	r := Repo{
		Owner: owner,
		Name:  repo,
		Token: Token{
			Token: nil,
		},
	}

	if token != "" {
		r.Token.Token = &token
	}

	return r
}

func (t Token) NewClient() *github.Client {

	tc := &http.Client{
		Transport: &oauth2.Transport{
			Base: httpcache.NewMemoryCacheTransport(),
		},
	}

	if t.Token != nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: *t.Token},
		)
		tc = &http.Client{
			Transport: &oauth2.Transport{
				Base:   httpcache.NewMemoryCacheTransport(),
				Source: ts,
			},
		}
	}

	return github.NewClient(tc)
}
