package main

import (
	"embed"
	"errors"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"

	"github.com/choirulanwar/textify/backend/cmd/command"
	"github.com/choirulanwar/textify/backend/cmd/file"
	"github.com/choirulanwar/textify/backend/cmd/keyword"
	"github.com/choirulanwar/textify/backend/cmd/keyword_trend_explorer"
	"github.com/choirulanwar/textify/backend/cmd/setting"
	"github.com/choirulanwar/textify/backend/cmd/version"
	"github.com/choirulanwar/textify/backend/config"
	"github.com/choirulanwar/textify/backend/pkg/log"
	"github.com/choirulanwar/textify/backend/service"
)

//go:embed all:frontend/dist
var assets embed.FS
var (
	env       = "test"
	frameless = true
)

func main() {
	if env != "test" && env != "prod" {
		panic(errors.New("the environment variable is abnormal, it can only be set to test or prod"))
	}

	config.ConfEnv = env
	app := service.NewApp()
	app.LogInfo("env:", env)
	app.LogInfo("sqlite db:", app.GetDatabasePath())
	app.LogInfo("runtime path:", log.GetRuntimePath())

	if runtime.GOOS == "darwin" {
		frameless = false
	}

	if err := wails.Run(&options.App{
		Title:             app.Cfg.App.AppName,
		Width:             app.Cfg.App.Width,
		Height:            app.Cfg.App.Height,
		MinWidth:          app.Cfg.App.MinWidth,
		MinHeight:         app.Cfg.App.MinHeight,
		MaxWidth:          app.Cfg.App.Width * 10,
		MaxHeight:         app.Cfg.App.Height * 10,
		DisableResize:     false,
		Frameless:         frameless,
		StartHidden:       false,
		HideWindowOnClose: false,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.OnStartUp,
		OnDomReady:        app.OnDomReady,
		OnBeforeClose:     app.OnBeforeClose,
		OnShutdown:        app.OnShutdown,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
		Bind: []interface{}{
			file.New(app),
			command.New(app),
			setting.New(app),
			version.New(app),
			keyword_trend_explorer.New(app),
			keyword.New(app),
		},
	}); err != nil {
		app.LogErrorf("Error: %s", err.Error())
		panic(err)
	}
}
