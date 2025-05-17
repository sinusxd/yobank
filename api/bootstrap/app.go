package bootstrap

import (
	"log"

	"gorm.io/gorm"
)

type Application struct {
	Env *Env
	DB  *gorm.DB
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()

	// Initialize PostgreSQL
	db, err := NewPostgresDatabase(app.Env)
	if err != nil {
		log.Fatal("PostgreSQL init error: ", err)
	}
	app.DB = db

	return *app
}

func NewPostgresDatabase(env *Env) (*gorm.DB, error) {
	config := Config{
		Host:     env.DBHost,
		Port:     env.DBPort,
		User:     env.DBUser,
		Password: env.DBPass,
		DBName:   env.DBName,
	}

	db, err := newPostgresDB(config)
	if err != nil {
		return nil, err
	}

	err = runMigrations(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *Application) CloseDBConnection() {
	sqlDB, err := app.DB.DB()
	if err != nil {
		log.Printf("Error getting underlying SQL DB: %v\n", err)
		return
	}
	err = sqlDB.Close()
	if err != nil {
		log.Printf("Error closing database connection: %v\n", err)
	}
}
