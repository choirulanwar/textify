package upgrade

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/choirulanwar/textify/backend/config"
	"github.com/choirulanwar/textify/backend/pkg/git"
	"github.com/choirulanwar/textify/backend/pkg/log"

	"github.com/hashicorp/go-version"
)

type Info struct {
	Cfg *config.Conf
}

type Latest struct {
	Version     string `json:"version"`     // 版本号
	VersionDes  string `json:"versionDes"`  // 版本描述
	Url         string `json:"url"`         // 安装包地址
	Md5Hash     string `json:"md5Hash"`     // hash值
	Size        int64  `json:"size"`        // 大小
	ReleaseDate string `json:"releaseDate"` // 版本时间
}

func New(cfg *config.Conf) *Info {
	return &Info{
		Cfg: cfg,
	}
}

// GetLastVersionInfo
// @Description: 获取版本信息
// @receiver i
// @return *Latest
func (i *Info) GetLastVersionInfo() (versionInfo Latest) {
	releaseList, err := git.GetGithubReleaseList(i.Cfg.Github.Owner, i.Cfg.Github.Repo)
	if err != nil {
		return
	}
	if len(releaseList) == 0 {
		return
	}
	versionInfo.Version = *releaseList[0].TagName
	versionInfo.VersionDes = *releaseList[0].Body
	versionInfo.Url += i.Cfg.App.UpgradeUrl + versionInfo.Version + "/" + i.newPackageName()
	versionInfo.ReleaseDate = releaseList[0].PublishedAt.String()
	return
}

// newPackageName
// @Description: 创建新包名
// @receiver i
// @return string
func (i *Info) newPackageName() string {
	var packageName = i.Cfg.App.AppName
	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "amd64": //64位AMD
			packageName += "-amd64-installer.exe"
		case "arm64": //64位ARM
			packageName += "-arm64-installer.exe"
		case "386": //32位Intel
			packageName += "-win32-installer.exe"
		}
	case "darwin": // 暂时只支持 amd64
		switch runtime.GOARCH {
		case "amd64":
			packageName += "-amd64-installer.dmg"
		case "arm64":
			packageName += "-arm64-installer.dmg"
		}
	}
	return packageName
}

// GetDownloadUrlInfo
// @Description: 获取下载资源信息
// @receiver i
// @param url
// @return map[string]interface{}
// @author cx
func (i *Info) GetDownloadUrlInfo(proxyUrl, downloadUrl string) (map[string]interface{}, error) {
	const DownloadSpeed = 5.0 // 下载速度为 5MB/s
	proxy := http.ProxyFromEnvironment
	if proxyUrl != "" {
		_proxyUrl, err := url.Parse(proxyUrl)
		if err != nil {
			panic(err)
		}
		proxy = http.ProxyURL(_proxyUrl)
	}
	client := &http.Client{
		Transport: &http.Transport{Proxy: proxy},
	}
	// 发送 HEAD 请求获取文件大小和文件类型等信息
	res, err := client.Head(downloadUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// 获取文件大小
	filesize := res.ContentLength
	// 计算预估下载时间
	downloadTime := float64(filesize) / (1024.0 * 1024.0 * DownloadSpeed)
	return map[string]interface{}{
		"file_size":     filesize,
		"download_time": downloadTime,
	}, nil
}

// downloadLatestVersion
// @Description: 下载最新版本
// @receiver i
// @param downloadUrl
// @return string
// @return error
// @author cx
func (i *Info) downloadLatestVersion(proxyUrl, downloadUrl string) (string, error) {
	var databaseHome string
	proxy := http.ProxyFromEnvironment
	if proxyUrl != "" {
		_proxyUrl, err := url.Parse(proxyUrl)
		if err != nil {
			panic(err)
		}
		proxy = http.ProxyURL(_proxyUrl)
	}
	client := &http.Client{
		Transport: &http.Transport{Proxy: proxy},
	}
	fmt.Println("proxyUrl", proxyUrl)
	resp, err := client.Get(downloadUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	downloadPath := strings.Split(downloadUrl, "/")
	packageName, err := url.QueryUnescape(downloadPath[len(downloadPath)-1])
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		databaseHome = os.Getenv("APPDATA")
	} else {
		databaseHome, _ = os.UserHomeDir()
	}
	pathExecutable, _ := os.Executable()
	_, pathAppname := filepath.Split(pathExecutable)
	filePath := databaseHome + "/" + pathAppname + "/" + packageName
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return filePath, err
}

// RestartApplication
// @Description: 执行安装应用
// @receiver i
// @return error
// @author cx
func (i *Info) RestartApplication(filePath string) error {
	log.PrintInfo("RestartApplication", filePath)
	// 启动新的进程
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "start "+filePath)
		return cmd.Start()
	case "darwin":
		cmd := exec.Command("open", filePath)
		return cmd.Start()
	}
	return fmt.Errorf("请自行安装 新版本存放目录：%s", filePath)
}

// DoUpgrade
// @Description: 更新操作
// @receiver i
// @return *Latest
// @return error
func (i *Info) DoUpgrade(proxyUrl string, versionInfo Latest) (string, error) {
	if versionInfo.Version == "" {
		return "", fmt.Errorf("版本获取失败")
	}
	localVer, _ := version.NewVersion(i.Cfg.App.Version)
	remoteVer, _ := version.NewVersion(versionInfo.Version)
	if !localVer.LessThan(remoteVer) {
		return "", fmt.Errorf("当前版本已是最新版本")
	}
	filePath, err := i.downloadLatestVersion(proxyUrl, versionInfo.Url)
	if err != nil {
		return "", fmt.Errorf("下载最新版本失败 err: %s", err.Error())
	}
	return filePath, nil
}
