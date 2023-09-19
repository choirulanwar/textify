package service

import (
	"context"
	"fmt"

	"github.com/choirulanwar/textify/backend/config"
	"github.com/choirulanwar/textify/backend/pkg/log"
	"github.com/choirulanwar/textify/backend/pkg/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PushResp struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

type App struct {
	Ctx      context.Context
	Cfg      *config.Conf
	Log      *logrus.Logger
	DB       *gorm.DB
	ClientId string
}

func NewApp() *App {
	conf, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	logger := log.NewLogger(conf.Log.DirPath, conf.Log.FileName, conf.Log.Debug)
	db, err := sqlite.WithConnect(logger, conf)
	if err != nil {
		panic(err)
	}
	return &App{
		Ctx: context.Background(),
		Cfg: conf,
		Log: logger,
		DB:  db,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.Ctx = ctx
}

func (a *App) OnStartUp(ctx context.Context) {
	a.Ctx = ctx
	return
}

func (a *App) OnDomReady(ctx context.Context) {
	a.Log.Info("DomReady")
	return
}

func (a *App) OnShutdown(ctx context.Context) {
	a.Log.Info("Shutdown")
	return
}

func (a *App) OnBeforeClose(ctx context.Context) bool {
	a.Log.Info("BeforeClose")
	// 返回true将阻止程序关闭
	return false
}

func (a *App) GetDatabasePath() string {
	return sqlite.GetDatabasePath()
}

func (a *App) GetIsAutoMigrate() bool {
	return sqlite.GetIsAutoMigrate()
}

// LogInfo ...
func (a *App) LogInfo(args ...interface{}) {
	a.Log.Info(args...)
	log.PrintInfo(args...)
}

// LogInfof ...
func (a *App) LogInfof(format string, args ...interface{}) {
	a.Log.Infof(format, args...)
	log.PrintInfo(fmt.Sprintf(format, args...))
}

// LogError ...
func (a *App) LogError(args ...interface{}) {
	a.Log.Error(args...)
	log.PrintError(args...)
}

// LogErrorf ...
func (a *App) LogErrorf(format string, args ...interface{}) {
	a.Log.Errorf(format, args...)
	log.PrintError(fmt.Sprintf(format, args...))
}

// Println ...
func (a *App) Println(level logrus.Level, args ...interface{}) {
	if level == logrus.InfoLevel {
		a.LogInfo(args...)
	} else {
		a.Log.Log(level, args...)
		log.PrintError(args...)
	}
}
