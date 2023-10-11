package setting

import (
	"fmt"
	"net/url"
	rt "runtime"

	setting_repository "github.com/choirulanwar/textify/backend/internal/repositories/setting"
	"github.com/choirulanwar/textify/backend/models"
	"github.com/choirulanwar/textify/backend/pkg/resp"
	"github.com/choirulanwar/textify/backend/service"
	"github.com/google/go-github/v51/github"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Service struct {
	App *service.App
}

func New(a *service.App) *Service {
	return &Service{
		App: a,
	}
}

func (s *Service) GetGeneralInfo() *resp.Response {
	settingRepository := setting_repository.NewSettingRepository(s.App.DB)

	if result, err := settingRepository.Find(); err != nil {
		return resp.Fail(err.Error())
	} else {
		return resp.Success(result)
	}
}

func (s *Service) SetGeneralData(data models.Setting) *resp.Response {
	settingRepository := setting_repository.NewSettingRepository(s.App.DB)

	if err := settingRepository.Update(&data); err != nil {
		return resp.Fail(err.Error())
	} else {
		return resp.Success(data)
	}
}

func (s *Service) FeedBack(data models.FeedBack) *resp.Response {
	client := github.NewClient(nil)
	orgs, githubResp, err := client.Issues.Create(s.App.Ctx, s.App.Cfg.Github.Owner, s.App.Cfg.Github.Repo, &github.IssueRequest{
		Title:       &data.Title,
		Body:        &data.Body,
		Labels:      &data.Labels,
		Assignees:   &data.Assignees,
		State:       &data.State,
		StateReason: &data.StateReason,
		Milestone:   &data.Milestone,
	})
	fmt.Println(orgs, githubResp, err)
	if err != nil {
		return resp.Fail("feedback failed err:" + err.Error())
	}
	return resp.Success("")
}

func (s *Service) GetFeedBackUrl(data models.FeedbackReq) *resp.Response {
	body := "- [ ] I'm sure this does not appear in [the issue list of the repository](https://github.com/choirulanwar/textify/issues) "
	if data.IssueType == 1 {
		body += fmt.Sprintf("%s ## Basic Info:%s - Version: %s ## Steps to reproduce: %s", "%0A", "%0A", data.Version+"%0A", "%0A"+data.Body+"%0A")
	} else {
		body += fmt.Sprintf("%s ## Basic Info:%s - Version: %s ## What is expected?: %s", "%0A", "%0A", data.Version+"%0A", "%0A"+data.Body+"%0A")
	}
	parseUrl, _ := url.Parse("https://github.com/" + s.App.Cfg.Github.Owner + "/" + s.App.Cfg.Github.Repo + "/issues/new?title=" + data.Title + "&body=" + body)
	return resp.Success(parseUrl.String())
}

func (s *Service) GetGithubReleaseList() *resp.Response {
	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(s.App.Ctx, s.App.Cfg.Github.Owner, s.App.Cfg.Github.Repo, nil)
	if err != nil {
		return resp.Fail(err.Error())
	}
	return resp.Success(releases)
}

func (s *Service) Minimise() {
	runtime.WindowMinimise(s.App.Ctx)
}

func (s *Service) Fullscreen() {
	runtime.WindowFullscreen(s.App.Ctx)
}

func (s *Service) NormalScreen() {
	runtime.WindowSetSize(s.App.Ctx, s.App.Cfg.App.Width, s.App.Cfg.App.Height)
}

func (s *Service) Hide() {
	runtime.Hide(s.App.Ctx)
}

func (s *Service) Quit() {
	runtime.Quit(s.App.Ctx)
}

func (s *Service) ReloadApp() {
	runtime.WindowReloadApp(s.App.Ctx)
}

func (s *Service) IsWindows() bool {
	return rt.GOOS == "darwin"
}
