package version

import (
	"github.com/choirulanwar/textify/backend/cmd/setting"
	"github.com/choirulanwar/textify/backend/models"
	"github.com/choirulanwar/textify/backend/pkg/resp"
	"github.com/choirulanwar/textify/backend/pkg/upgrade"
	"github.com/choirulanwar/textify/backend/service"
	"github.com/hashicorp/go-version"
)

type Service struct {
	App *service.App
}

func New(t *service.App) *Service {
	return &Service{
		App: t,
	}
}

func (s *Service) getLocalVersion() string {
	return s.App.Cfg.App.Version
}

func (s *Service) needUpdate(local string, remote string) bool {
	localVer, localVerErr := version.NewVersion(local)
	if localVerErr != nil {
		return false
	}
	remoteVer, remoteVerErr := version.NewVersion(remote)
	if remoteVerErr != nil {
		return false
	}
	return localVer.LessThan(remoteVer)
}

func (s *Service) CheckAndGetLatestServerVersion() *resp.Response {
	localVersion := s.getLocalVersion()
	remoteVersionInfo := upgrade.New(s.App.Cfg).GetLastVersionInfo()
	return resp.Success(map[string]interface{}{
		"needUpdate":  s.needUpdate(localVersion, remoteVersionInfo.Version),
		"versionInfo": remoteVersionInfo,
	})
}

func (s *Service) DoUpgrade() *resp.Response {
	remoteVersionInfo := upgrade.New(s.App.Cfg).GetLastVersionInfo()
	generalInfo := setting.New(s.App).GetGeneralInfo()
	settingInfo := generalInfo.Data.(models.Setting)
	filePath, err := upgrade.New(s.App.Cfg).DoUpgrade(settingInfo.ProxyUrl, remoteVersionInfo)
	if err != nil {
		return resp.Fail(err.Error())
	}
	return resp.Success(map[string]string{
		"file_path":    filePath,
		"download_url": remoteVersionInfo.Url,
	})
}

func (s *Service) GetDownloadUrlInfo(url string) *resp.Response {
	generalInfo := setting.New(s.App).GetGeneralInfo()
	settingInfo := generalInfo.Data.(models.Setting)
	info, err := upgrade.New(s.App.Cfg).GetDownloadUrlInfo(settingInfo.ProxyUrl, url)
	if err != nil {
		return resp.Fail("获取下载信息异常 err:" + err.Error())
	}
	return resp.Success(info)
}

func (s *Service) RestartApplication(filePath string) *resp.Response {
	if err := upgrade.New(s.App.Cfg).RestartApplication(filePath); err != nil {
		return resp.Fail(err.Error())
	}
	// s.App.ExitSignalChan <- true
	return resp.Success("")
}
