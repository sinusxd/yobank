package bootstrap

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"yobank/domain"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func newPostgresDB(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to database: %v\n", err)
		return nil, err
	}

	return db, nil
}

func runMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(&domain.User{}, &domain.EmailLoginCode{}, &domain.Wallet{}, &domain.Rate{}, &domain.Transfer{})
	if err != nil {
		log.Printf("Error running migrations: %v\n", err)
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
