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

	"github.com/choirulanwar/textify/backend/config"
	"github.com/choirulanwar/textify/backend/pkg/log"
	"github.com/choirulanwar/textify/backend/service"
	"github.com/choirulanwar/textify/backend/service/command"
	"github.com/choirulanwar/textify/backend/service/setting"
)

//go:embed all:frontend/dist
var assets embed.FS
var (
	env       = "test"
	app       *service.App
	frameless = true
)

func main() {
	if env != "test" && env != "prod" {
		panic(errors.New("环境变量异常，只能设置test、prod"))
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
		MinWidth:          app.Cfg.App.MinWidth,    // 最小宽度
		MinHeight:         app.Cfg.App.MinHeight,   // 最小高度
		MaxWidth:          app.Cfg.App.Width * 10,  // 最大宽度
		MaxHeight:         app.Cfg.App.Height * 10, // 最大高度
		DisableResize:     false,                   // 调整窗口尺寸
		Frameless:         frameless,               // 无边框
		StartHidden:       false,                   // 启动后隐藏
		HideWindowOnClose: false,                   // 关闭窗口将隐藏而不退出应用程序
		LogLevel:          logger.DEBUG,            // 日志级别
		OnStartup:         app.OnStartUp,
		OnDomReady:        app.OnDomReady,
		OnBeforeClose:     app.OnBeforeClose, //
		OnShutdown:        app.OnShutdown,    //
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
		Bind: []interface{}{
			command.New(app),
			setting.New(app),
		},
	}); err != nil {
		app.LogErrorf("Error: %s", err.Error())
		panic(err)
	}
}
