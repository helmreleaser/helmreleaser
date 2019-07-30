package scm

import (
	"github.com/google/go-github/v27/github"
)

type SCMClient struct {
	GitHub *github.Client
}
