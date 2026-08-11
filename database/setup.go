package database

import (
	"database/sql"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB_POSTGRES     *gorm.DB
	DB_SQL_POSTGRES *sql.DB
)

func ConnnectDatabasePostgres() {
	err := godotenv.Load()
	if err != nil {
		sentry.CaptureException(err)
	}

	DB_POSTGRES_HOST := os.Getenv("DB_POSTGRES_HOST")
	DB_POSTGRES_PORT := os.Getenv("DB_POSTGRES_PORT")
	DB_POSTGRES_USER := os.Getenv("DB_POSTGRES_USER")
	DB_POSTGRES_PASS := os.Getenv("DB_POSTGRES_PASS")
	DB_POSTGRES_NAME := os.Getenv("DB_POSTGRES_NAME")

	dsn := "host=" + DB_POSTGRES_HOST + " user=" + DB_POSTGRES_USER +
		" password=" + DB_POSTGRES_PASS + " dbname=" + DB_POSTGRES_NAME +
		" port=" + DB_POSTGRES_PORT + " sslmode=disable TimeZone=Asia/Jakarta"

	// GORM connection
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		sentry.CaptureException(err)
	}
	DB_POSTGRES = db

	// Raw database/sql connection (for insight raw queries)
	sqlDB, sqlErr := sql.Open("pgx", dsn)
	if sqlErr != nil {
		sentry.CaptureException(sqlErr)
	}
	if sqlDB != nil {
		if err := sqlDB.Ping(); err != nil {
			sentry.CaptureException(err)
		}
	}
	DB_SQL_POSTGRES = sqlDB
}
