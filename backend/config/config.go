package config

import (
	"bytes"
	"embed"

	"github.com/spf13/viper"
)

var ConfEnv string = "test"

type (
	Conf struct {
		App    App    `yaml:"app"`
		Log    Log    `yaml:"log"`
		Github Github `yaml:"github"`
	}
	App struct {
		Mode            string `yaml:"mode" default:"debug"`
		DefaultPageSize int    `yaml:"defaultPageSize" default:"10"`
		MaxPageSize     int    `yaml:"maxPageSize" default:"500"`
		AppName         string `yaml:"appName" default:""`
		Version         string `yaml:"version" default:""`
		Description     string `yaml:"description" default:""`
		VersionNewMsg   string `yaml:"versionNewMsg" default:""`
		VersionOldMsg   string `yaml:"versionOldMsg" default:""`
		BtnConfirmText  string `yaml:"btnConfirmText" default:"确定"`
		BtnCancelText   string `yaml:"btnCancelText" default:"取消"`
		Width           int    `yaml:"width" default:"1024"`
		Height          int    `yaml:"height" default:"768"`
		MinWidth        int    `yaml:"minWidth" default:"1000"`
		MinHeight       int    `yaml:"minHeight" default:"700"`
		UpgradeUrl      string `yaml:"upgradeUrl" default:""`
		HomeUrl         string `yaml:"homeUrl" default:""`
		DbName          string `yaml:"dbName" default:""`
	}
	Log struct {
		Debug    bool   `yaml:"debug" default:"true"`
		FileName string `yaml:"fileName" default:"syncstock"`
		DirPath  string `yaml:"dirPath" default:"runtime/logs"`
	}
	Github struct {
		Owner string `yaml:"owner" default:""`
		Repo  string `yaml:"repo" default:""`
	}
)

//go:embed yaml
var yamlCfg embed.FS

func InitConfig() (*Conf, error) {
	var cfg *Conf
	v := viper.New()
	v.SetConfigType("yaml")
	yamlConf, _ := yamlCfg.ReadFile("yaml/config." + ConfEnv + ".yaml")
	if err := v.ReadConfig(bytes.NewBuffer(yamlConf)); err != nil {
		return nil, err
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
