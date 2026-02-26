package main

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

var auth = &http.BasicAuth{
	Username: GH_USERNAME,
	Password: GH_TOKEN,
}

type Repository struct {
	repo *git.Repository
	url  string
	dest string
}

func NewRepository(dest string, url string) *Repository {
	return &Repository{
		url:  url,
		dest: dest,
	}
}

func (r *Repository) Clone() error {
	fmt.Println("Cloning...")

	repo, err := git.PlainClone(r.dest, &git.CloneOptions{
		URL:          r.url,
		Depth:        1,
		Auth:         auth,
		SingleBranch: true,
	})
	if err != nil && err.Error() != "remote repository is empty" {
		return err
	}

	r.repo = repo

	return nil
}

func (r *Repository) Worktree() (*git.Worktree, error) {
	return r.repo.Worktree()
}

func (r *Repository) Status() (git.Status, error) {
	worktree, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (r *Repository) AddAll() error {
	worktree, err := r.Worktree()
	if err != nil {
		return err
	}

	return worktree.AddWithOptions(&git.AddOptions{
		All: true,
	})
}

func (r *Repository) Commit(message string) error {
	worktree, err := r.Worktree()
	if err != nil {
		return err
	}

	fmt.Println("Commit with message:", message)

	_, err = worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  GH_USERNAME,
			Email: GH_EMAIL,
			When:  time.Now(),
		},
	})

	return err
}

func (r *Repository) Push() error {
	fmt.Println("Pushing...")

	return r.repo.Push(&git.PushOptions{
		Auth: auth,
	})
}
