package module

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"net/http"
	"teamide/internal/base"
)

const (
	GithubReleasesURL       = `https://github.com/team-ide/teamide/releases`
	GithubReleaseHistoryURL = `https://github.com/team-ide/teamide/blob/main/CHANGELOG.md`
	GiteeReleaseHistoryURL  = `https://gitee.com/teamide/teamide/blob/main/CHANGELOG.md`
)

type UpdateCheckRequest struct {
}

type UpdateCheckResponse struct {
	CurrentVersion    string `json:"currentVersion,omitempty"`
	ReleaseHistory    string `json:"releaseHistory,omitempty"`
	GithubReleasesURL string `json:"githubReleasesURL,omitempty"`
}

func (this_ *Api) apiUpdateCheck(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &UpdateCheckRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &UpdateCheckResponse{}

	response.CurrentVersion = this_.Version
	response.GithubReleasesURL = GithubReleasesURL
	response.ReleaseHistory, err = releasesCheck()

	res = response
	return
}

func releasesCheck() (releaseHtml string, err error) {
	releaseHtml, err = releasesCheckGitee()
	if err != nil || releaseHtml == "" {
		releaseHtml, err = releasesCheckGithub()
		return
	}
	return
}

func releasesCheckGithub() (releaseHtml string, err error) {
	res, err := http.Get(GithubReleaseHistoryURL)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	selection := doc.Find(".markdown-body")

	releaseHtml, err = selection.Html()
	return
}

func releasesCheckGitee() (releaseHtml string, err error) {
	res, err := http.Get(GiteeReleaseHistoryURL)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	selection := doc.Find(".markdown-body")

	releaseHtml, err = selection.Html()
	return
}
