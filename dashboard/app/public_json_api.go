// Copyright 2021 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package main

import "encoding/json"

// publicApiBugDescription is used to serve the /bug HTTP requests
// and provide JSON description of the BUG. Backward compatible.
type publicAPIBugDescription struct {
	Version    int         `json:"version"`
	Title      string      `json:"title,omitempty"`
	FixCommits []fixCommit `json:"fix-commits,omitempty"`
	// links to the discussions
	Discussions []string                    `json:"discussions,omitempty"`
	Crashes     []publicAPICrashDescription `json:"crashes,omitempty"`
}

type fixCommit struct {
	Title string `json:"title"`
	Link  string `json:"link,omitempty"`
	Hash  string `json:"hash,omitempty"`
}

type publicAPICrashDescription struct {
	SyzReproducer       string `json:"syz-reproducer,omitempty"`
	CReproducer         string `json:"c-reproducer,omitempty"`
	KernelConfig        string `json:"kernel-config,omitempty"`
	KernelSourceGit     string `json:"kernel-source-git,omitempty"`
	KernelSourceCommit  string `json:"kernel-source-commit,omitempty"`
	SyzkallerGit        string `json:"syzkaller-git,omitempty"`
	SyzkallerCommit     string `json:"syzkaller-commit,omitempty"`
	CompilerDescription string `json:"compiler-description,omitempty"`
	Architecture        string `json:"architecture,omitempty"`
}

func getExtAPIDescrForBugPage(bugPage *uiBugPage) *publicAPIBugDescription {
	crash := bugPage.Crashes.Crashes[0]
	return &publicAPIBugDescription{
		Version: 1,
		Title:   bugPage.Bug.Title,
		Discussions: func() []string {
			if bugPage.Bug.ExternalLink == "" {
				return nil
			}
			return []string{bugPage.Bug.ExternalLink}
		}(),
		FixCommits: func() []fixCommit {
			if len(bugPage.Bug.Commits) == 0 {
				return nil
			}
			var res []fixCommit
			for _, commit := range bugPage.Bug.Commits {
				res = append(res, fixCommit{
					Title: commit.Title,
					Hash:  commit.Hash,
					Link:  commit.Link,
				})
			}
			return res
		}(),
		Crashes: []publicAPICrashDescription{{
			SyzReproducer:      crash.ReproSyzLink,
			CReproducer:        crash.ReproCLink,
			KernelConfig:       crash.KernelConfigLink,
			KernelSourceGit:    crash.KernelCommitLink,
			KernelSourceCommit: crash.KernelCommit,
			SyzkallerGit:       crash.SyzkallerCommitLink,
			SyzkallerCommit:    crash.SyzkallerCommit,
			// TODO: add the CompilerDescription
			// TODO: add the Architecture
		}},
	}
}

type publicAPIBugGroup struct {
	Version int `json:"version"`
	Bugs    []publicAPIBug
}

type publicAPIBug struct {
	Title       string `json:"title,omitempty"`
	Link        string `json:"link"`
	LastUpdated string `json:"last-updated,omitempty"`
}

func getExtAPIDescrForBugGroups(bugGroups []*uiBugGroup) *publicAPIBugGroup {
	return &publicAPIBugGroup{
		Version: 1,
		Bugs: func() []publicAPIBug {
			var res []publicAPIBug
			for _, group := range bugGroups {
				for _, bug := range group.Bugs {
					res = append(res, publicAPIBug{
						Title: bug.Title,
						Link:  bug.Link,
					})
				}
			}
			return res
		}(),
	}
}

func GetJSONDescrFor(page interface{}) ([]byte, error) {
	var res interface{}
	switch i := page.(type) {
	case *uiBugPage:
		res = getExtAPIDescrForBugPage(i)
	case *uiTerminalPage:
		res = getExtAPIDescrForBugGroups([]*uiBugGroup{i.Bugs})
	case *uiMainPage:
		res = getExtAPIDescrForBugGroups(i.Groups)
	default:
		return nil, ErrClientNotFound
	}
	return json.MarshalIndent(res, "", "\t")
}
