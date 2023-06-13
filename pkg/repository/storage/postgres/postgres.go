package postgres

import (
	"brief/internal/config"
	"brief/internal/model"
	"brief/utility"
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

var (
	db                  *gorm.DB
	generalQueryTimeout = 60 * time.Second
)

func GetDB() *Postgres {
	return &Postgres{db}
}

func ConnectToDB() *gorm.DB {
	logger := utility.NewLogger()

	database, err := gorm.Open(postgres.Open(dsn()), &gorm.Config{})
	if err != nil {
		log.Fatalf("could not connect to postgres, got error: %s", err)
	}
	db = database

	if err := migrateDB(logger); err != nil {
		log.Fatalf("could not run db migrations, got error: %s", err)
	}

	// IF EVERYTHING IS OKAY, THEN CONNECTION IS ESTABLISHED
	fmt.Println("POSTGRES CONNECTION ESTABLISHED")
	logger.Info("POSTGRES CONNECTION ESTABLISHED")

	return db
}

func dsn() string {
	pgHost := config.GetConfig().PGHost
	pgPort := config.GetConfig().PGPort
	pgUser := config.GetConfig().PGUser
	pgDB := config.GetConfig().PGDatabase
	pgPassword := config.GetConfig().PGPassword

	dsn := "host=" + pgHost + " user=" + pgUser +
		" password=" + pgPassword + " dbname=" + pgDB + " port=" + pgPort + " sslmode=disable"

	return dsn
}

// migrateDB creates db schemas
func migrateDB(logger *utility.Logger) error {
	err := db.AutoMigrate(
		&model.User{},
	)
	if err != nil {
		return err
	}

	logger.Info("DATABASE MIGRATION SUCCESSFUL")
	return nil
}

// DBWithTimeout returns a database with timeout, and the context's cancel func
func (p *Postgres) DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, generalQueryTimeout)
	return p.db.WithContext(ctx), cancel
}
