package sqlite

import (
	"fmt"
	"time"

	"github.com/choirulanwar/textify/backend/config"
	"github.com/choirulanwar/textify/backend/models"
	"github.com/choirulanwar/textify/backend/pkg/helper"
	"github.com/choirulanwar/textify/backend/pkg/log"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	isAutoMigrate bool
	databasePath  string
)

func WithConnect(logger *logrus.Logger, conf *config.Conf) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)
	runtimeDir := helper.GetRuntimeUserHomeDir()
	databasePath = runtimeDir + "/" + conf.App.AppName + "/" + conf.App.DbName
	fmt.Println("dbPath", databasePath)
	if !helper.IsPathExist(databasePath) {
		_, err := helper.MakeFileOrPath(databasePath)
		if err != nil {
			return nil, err
		}
		db, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		log.PrintInfo("======== migrate start ========")
		logger.Info("======== migrate start ========")
		autoMigrate(db)
		log.PrintInfo("======== migrate end ========")
		logger.Info("======== migrate end ========")
		isAutoMigrate = true
	} else {
		db, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)

	sqlDB.SetMaxOpenConns(100)

	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, err
}

func GetDatabasePath() string {
	return databasePath
}

func GetIsAutoMigrate() bool {
	return isAutoMigrate
}

func autoMigrate(db *gorm.DB) {
	db.Migrator().CreateTable(
		&models.Setting{},
		&models.KeywordTrendExplorer{},
		&models.Keyword{},
		// &models.SERPResult{},
		// &models.Link{},
		// &models.Trend{},
		// &models.RelatedQuestion{},
	)

	MockSetting(db)
	MockKeyword(db)
	MockTrend(db)
	MockTagList(db)
}
