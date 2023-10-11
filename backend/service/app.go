package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/choirulanwar/textify/backend/config"
	handlers_queues "github.com/choirulanwar/textify/backend/internal/queues/handlers"
	tasks_queues "github.com/choirulanwar/textify/backend/internal/queues/tasks"
	keyword_repository "github.com/choirulanwar/textify/backend/internal/repositories/keyword"
	keyword_trend_explorer_repository "github.com/choirulanwar/textify/backend/internal/repositories/keyword_trend_explorer"
	"github.com/choirulanwar/textify/backend/pkg/log"
	"github.com/choirulanwar/textify/backend/pkg/queue"
	"github.com/choirulanwar/textify/backend/pkg/sqlite"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PushResp struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

type App struct {
	Ctx         context.Context
	Cfg         *config.Conf
	Log         *logrus.Logger
	DB          *gorm.DB
	ClientId    string
	AsynqClient *asynq.Client
	RedisOpt    *asynq.RedisClientOpt
	exDir       string
	cmdPrefix   string
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

	redisOpt := asynq.RedisClientOpt{
		Addr: "localhost:6379",
	}

	asynqClient := queue.NewAsynqClient(redisOpt)

	return &App{
		Ctx:         context.Background(),
		Cfg:         conf,
		Log:         logger,
		DB:          db,
		AsynqClient: asynqClient,
		RedisOpt:    &redisOpt,
		exDir:       "",
		cmdPrefix:   "",
	}
}

func (a *App) OnStartUp(ctx context.Context) {
	a.Ctx = ctx
	a.RunWorker(*a.RedisOpt)
	if runtime.GOOS == "darwin" {
		ex, _ := os.Executable()
		a.exDir = filepath.Dir(ex) + "/../../../"
		a.cmdPrefix = "cd " + a.exDir + " && "
	}

}

func (a *App) OnDomReady(ctx context.Context) {
	a.Log.Info("DomReady")
}

func (a *App) OnShutdown(ctx context.Context) {
	a.Log.Info("Shutdown")
}

func (a *App) OnBeforeClose(ctx context.Context) bool {
	a.Log.Info("BeforeClose")
	return false
}

func (a *App) GetDatabasePath() string {
	return sqlite.GetDatabasePath()
}

func (a *App) GetIsAutoMigrate() bool {
	return sqlite.GetIsAutoMigrate()
}

func (a *App) LogInfo(args ...interface{}) {
	a.Log.Info(args...)
	log.PrintInfo(args...)
}

func (a *App) LogInfof(format string, args ...interface{}) {
	a.Log.Infof(format, args...)
	log.PrintInfo(fmt.Sprintf(format, args...))
}

func (a *App) LogError(args ...interface{}) {
	a.Log.Error(args...)
	log.PrintError(args...)
}

func (a *App) LogErrorf(format string, args ...interface{}) {
	a.Log.Errorf(format, args...)
	log.PrintError(fmt.Sprintf(format, args...))
}

func (a *App) Println(level logrus.Level, args ...interface{}) {
	if level == logrus.InfoLevel {
		a.LogInfo(args...)
	} else {
		a.Log.Log(level, args...)
		log.PrintError(args...)
	}
}

func (a *App) RunWorker(redisOpt asynq.RedisClientOpt) {
	a.LogInfo("[+] Worker start...")

	server := queue.NewAsynqServer(redisOpt, 50, map[string]int{
		tasks_queues.KeywordTrendExplorerQueueName: 3,
		tasks_queues.TrendMetricQueueName:          3,
		tasks_queues.SerpExplorerQueueName:         3,
	})

	mux := asynq.NewServeMux()
	mux.Use(a.workerMiddleware)

	addWorkerHandler(mux, a)

	go func(mux *asynq.ServeMux, server *asynq.Server) {
		if err := server.Run(mux); err != nil {
			a.LogErrorf("Queue job server failed %v", err)
		}
	}(mux, server)
}

func (a *App) workerMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
		start := time.Now()
		a.Log.WithFields(logrus.Fields{
			"Type":    task.Type(),
			"Payload": task.Payload(),
		}).Info("Start processing")
		err := h.ProcessTask(ctx, task)
		if err != nil {
			return err
		}
		a.Log.WithFields(logrus.Fields{
			"Elapsed Time": time.Since(start),
		}).Info("Finished processing")
		return nil
	})
}

func addWorkerHandler(mux *asynq.ServeMux, app *App) {
	keywordTrendExplorerRepository := keyword_trend_explorer_repository.NewKeywordTrendExplorerRepository(app.DB)
	keywordRepository := keyword_repository.NewKeywordRepository(app.DB)
	handlerQueues := handlers_queues.NewHandlerQueues(app.Ctx, app.Log, app.DB, app.AsynqClient, keywordTrendExplorerRepository, keywordRepository)

	handOnKeywordTrendExplorerRequestedTask := func(ctx context.Context, task *asynq.Task) error {
		return handlerQueues.HandOnKeywordTrendExplorerRequestedTask(ctx, task)
	}
	handOnTrendMetricRequestedTask := func(ctx context.Context, task *asynq.Task) error {
		return handlerQueues.HandOnTrendMetricRequestedTask(ctx, task)
	}

	handOnSerpExplorerRequestedTask := func(ctx context.Context, task *asynq.Task) error {
		return handlerQueues.HandOnSerpExplorerRequestedTask(ctx, task)
	}

	mux.HandleFunc(tasks_queues.TypeOnKeywordTrendExplorerRequested, handOnKeywordTrendExplorerRequestedTask)
	mux.HandleFunc(tasks_queues.TypeOnTrendMetricRequested, handOnTrendMetricRequestedTask)
	mux.HandleFunc(tasks_queues.TypeOnSerpExplorerRequested, handOnSerpExplorerRequestedTask)
}
