package database

import (
	"ai-customer-service/domain"
	"ai-customer-service/internal/configs"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgreSQLService struct {
	Config *configs.Config
	Conn   *gorm.DB
}

var psqlDB *gorm.DB
var psqlOnce sync.Once

func NewPostgreSQL(config *configs.Config) domain.PostgreSQLInterface {
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

		psqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}

		psqlDB.SetMaxOpenConns(25)
		psqlDB.SetMaxIdleConns(5)
		psqlDB.SetConnMaxLifetime(5 * time.Minute)
	})

	return &PostgreSQLService{
		Config: config,
		Conn:   psqlDB,
	}
}
