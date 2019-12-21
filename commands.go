package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

const (
	commitFlag = "commit"
	branchFlag = "branch"
)

type constError string

const (
	errPathNotPassed = constError("path not passed")
)

func (e constError) Error() string { return string(e) }

func formatCommit(c interface{}) string {
	switch t := c.(type) {
	case *plumbing.Reference:
		return t.Name().Short()
	case *object.Commit:
		return "v0.0.0-" +
			t.Committer.When.UTC().Format("20060102150405") +
			"-" +
			t.Hash.String()[:12]
	default:
		return "could not find version"
	}
}

func latestRevision() (err error) {
	var (
		tmp  string
		res  interface{}
		lc   *object.Commit
		repo *git.Repository
		ctx  = newGracefulContext()
		path = viper.GetString("path")
	)

	if path == "" {
		return errPathNotPassed
	}

	defer func() {
		if err := os.RemoveAll(tmp); err != nil && tmp != "" {
			fmt.Println("Removing temp: ", err.Error())
		}

		if err != nil && lc != nil {
			res = lc
		}

		fmt.Println()

		if res != nil {
			fmt.Println("go.mod version: " + formatCommit(res))
			os.Exit(0)
		}
	}()

	if _, err = os.Stat(path); err == nil {
		if repo, err = git.PlainOpen(path); err != nil {
			return err
		}
	} else {
		sha := sha256.Sum256([]byte(path))

		if tmp, err = ioutil.TempDir("", hex.EncodeToString(sha[:])); err != nil {
			return err
		}

		var u *url.URL
		if u, err = url.Parse(path); err == nil {
			if u.Scheme == "" {
				u.Scheme = "https"
			}

			path = u.String()
		}

		repo, err = git.PlainCloneContext(ctx, tmp, true, &git.CloneOptions{
			URL:           path,
			Depth:         1,
			Progress:      os.Stdout,
			Tags:          git.AllTags,
			ReferenceName: plumbing.NewBranchReferenceName(viper.GetString(branchFlag)),
		})

		if err != nil {
			return err
		}
	}

	var (
		iter object.CommitIter
		ref  *plumbing.Reference
		tags storer.ReferenceIter
	)

	if ref, err = repo.Head(); err != nil {
		return err
	} else if iter, err = repo.Log(&git.LogOptions{From: ref.Hash()}); err != nil {
		return err
	} else if lc, err = iter.Next(); err != nil {
		return err
	} else if viper.GetBool(commitFlag) {
		res = lc
		return
	} else if tags, err = repo.Tags(); err != nil {
		return
	}

	var latestTag *object.Commit

	return tags.ForEach(func(ref *plumbing.Reference) error {
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
			res = ref
		}

		if commit.Committer.When.After(latestTag.Committer.When) {
			latestTag = commit
			res = ref
		}

		return nil
	})
}
