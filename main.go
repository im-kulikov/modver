package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func halt(msg string) {
	fmt.Println(msg)
	os.Exit(2)
}

func check(err error) {
	if err == nil {
		return
	}
	halt(err.Error())
}

func formatCommit(c *object.Commit) string {
	return "v0.0.0-" +
		c.Author.When.UTC().Format("20060102150405") +
		"-" +
		c.Hash.String()[:12]
}

func latestRevision(path string) (string, error) {
	var (
		err  error
		repo *git.Repository
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-ch
		cancel()
	}()

	if _, err = os.Stat(path); err == nil {
		if repo, err = git.PlainOpen(path); err != nil {
			return "", err
		}
	} else {
		sha := sha256.Sum256([]byte(path))

		tmp, err := ioutil.TempDir("", hex.EncodeToString(sha[:]))
		if err != nil {
			return "", err
		}

		defer func() {
			if err := os.RemoveAll(tmp); err != nil {
				fmt.Println(err.Error())
			}
		}()

		if u, err := url.Parse(path); err == nil {
			if u.Scheme == "" {
				u.Scheme = "https"
			}

			path = u.String()
		}

		repo, err = git.PlainCloneContext(ctx, tmp, true, &git.CloneOptions{
			URL:      path,
			Depth:    1,
			Progress: os.Stdout,
			Tags:     git.AllTags,
		})

		if err != nil {
			return "", err
		}
	}

	ref, err := repo.Head()
	check(err)

	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	check(err)

	latestCommit, err := cIter.Next()
	check(err)

	tags, err := repo.Tags()
	if err != nil {
		return formatCommit(latestCommit), nil
	}

	var latestTag *object.Commit

	var result string

	err = tags.ForEach(func(ref *plumbing.Reference) error {
		rev := plumbing.Revision(ref.Name().String())
		tag, err := repo.ResolveRevision(rev)
		if err != nil {
			return err
		}

		commit, err := repo.CommitObject(*tag)
		if err != nil {
			return err
		}

		if latestTag == nil {
			latestTag = commit
			result = ref.Name().Short()
		}

		if commit.Committer.When.After(latestTag.Committer.When) {
			latestTag = commit
			result = ref.Name().Short()
		}

		return nil
	})

	if err != nil || result == "" {
		return formatCommit(latestCommit), nil
	}

	return result, nil
}

func main() {
	if len(os.Args) != 2 {
		halt("You should provide path to repo.")
	}

	rev, err := latestRevision(os.Args[1])
	check(err)

	fmt.Println("go.mod version: " + rev)
}
