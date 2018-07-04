package pica

import (
	"fmt"
	"time"

	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/utils/diff"
)

const (
	Version = "0.0.1"
)

type VersionChange struct {
	Commit *object.Commit
	Diffs  []diffmatchpatch.Diff
}

type VersionNote struct {
	Changes []VersionChange
}

type ApiVersionController struct {
	rep      *git.Repository
	FileName string
}

func NewApiVersionController(filename string) *ApiVersionController {
	fmt.Print(filename)
	r, err := git.PlainOpen(".")
	if err != nil {
		panic(err)
	}
	return &ApiVersionController{
		rep:      r,
		FileName: filename,
	}
}

func (v *ApiVersionController) Commit(msg string) (string, error) {
	w, err := v.rep.Worktree()
	if err != nil {
		return "", err
	}
	_, err = w.Add(v.FileName)
	if err != nil {
		return "", err
	}
	msg = fmt.Sprintf("[Pica] %s", msg)
	hash, err := w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "jeremaihloo",
			Email: "jeremaihloo1024@gmail.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return "", err
	}
	return hash.String(), nil
}

func (v *ApiVersionController) Diff(src, dst string) []diffmatchpatch.Diff {
	return diff.Do(src, dst)
}

func (v *ApiVersionController) GetCommits() ([]*object.Commit, error) {
	ref, err := v.rep.Head()
	if err != nil {
		return nil, err
	}
	logs, err := v.rep.Log(&git.LogOptions{
		From: ref.Hash(),
	})
	if err != nil {
		return nil, err
	}
	var commits []*object.Commit
	err = logs.ForEach(func(commit *object.Commit) error {
		// find [Pica] and filename
		tree, err := commit.Tree()
		if err != nil {
			return err
		}
		flag := false
		// ... get the files iterator and print the file
		tree.Files().ForEach(func(f *object.File) error {
			if f.Name == v.FileName {
				flag = true
				return err
			}
			return nil
		})
		if flag && strings.Index(commit.Message, "[Pica]") > -1 {
			commits = append(commits, commit)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return commits, nil
}

func (v *ApiVersionController) Notes() (*VersionNote, error) {
	vn := &VersionNote{}
	commits, err := v.GetCommits()
	if err != nil {
		return nil, err
	}
	if len(commits) < 2 {
		for _, commit := range commits {
			commit.Message = strings.Replace(commit.Message, "[Pica]", "", -1)
			vc := VersionChange{
				Commit: commit,
			}
			vn.Changes = append(vn.Changes, vc)
		}
	}
	for index := 0; index < len(commits)-1; index++ {
		commit := commits[index]
		getContent := func(commit *object.Commit) (string, error) {
			content := ""
			tree, err := commit.Tree()
			if err != nil {
				return "", err
			}
			// ... get the files iterator and print the file
			tree.Files().ForEach(func(f *object.File) error {
				if f.Name == v.FileName {
					content, err = f.Contents()
					return err
				}
				return nil
			})
			return content, nil
		}
		currentContext, err := getContent(commit)
		if err != nil {
			return nil, err
		}
		lastContent, err := getContent(commits[index+1])
		if err != nil {
			return nil, err
		}

		vc := VersionChange{
			Commit: commit,
			Diffs:  v.Diff(currentContext, lastContent),
		}
		vn.Changes = append(vn.Changes, vc)
	}
	return vn, nil
}
