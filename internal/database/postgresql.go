package database

import (
	"ai-customer-service/internal/configs"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgreSQLService struct {
	Conn *gorm.DB
}

var psqlDB *gorm.DB
var psqlOnce sync.Once

func NewPostgreSQL(config *configs.Config) *gorm.DB {
	psqlOnce.Do(func() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
			config.PostgreSQL.Host,
			config.PostgreSQL.Username,
			config.PostgreSQL.Password,
			config.PostgreSQL.Database,
			config.PostgreSQL.Port,
			config.PostgreSQL.Timezone,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			panic(err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}

		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
		psqlDB = db
	})

	return psqlDB
}

func (p *PostgreSQLService) DB() *gorm.DB {
	return p.Conn
}
