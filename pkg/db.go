package pkg

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/serg666/card_storage/models"
)

type DBConnector struct {
	dsn    string
	logLvl logger.LogLevel
	DB     *gorm.DB
}

func (c *DBConnector) Open() error {
	db, err := gorm.Open(postgres.Open(c.dsn), &gorm.Config{})
	if err != nil {
		logrus.Errorf("Error open connection: err=%v", err)
	}
	c.DB = db.Session(&gorm.Session{
		Logger: db.Logger.LogMode(c.logLvl),
	})
	return err
}

func (c *DBConnector) Close() {
	if c.DB == nil {
		return
	}

	if db, err := c.DB.DB(); err == nil {
		logrus.Infof("DB close status: %v", db.Close())
	}
}

func (c *DBConnector) Migrate() error {
	return c.DB.AutoMigrate(
		&models.GormCard{},
		&models.GormSession{},
		&models.GormKey{},
	)
}

func (c *DBConnector) DropTables() error {
	return c.DB.Migrator().DropTable(
		&models.GormKey{},
		&models.GormCard{},
		&models.GormSession{},
	)
}

func NewDBConnector(dsn string, safeLogs bool) *DBConnector {
	var logLvl = logger.Silent
	if !safeLogs {
		logLvl = logger.Info
	}

	return &DBConnector{
		dsn:    dsn,
		logLvl: logLvl,
	}
}
