package db

import (
	"fmt"
	"sync"

	"github.com/guimsmendes/logbookus/config"
	"github.com/guimsmendes/logbookus/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	mutex sync.Mutex
	db    *gorm.DB
)

// Connect returns the database connection pool, creating it if no pool exist. The function is safe to call
// from multiple Go routines.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// NOTE: the err cannot be "cached" in a global var as access from multiple go routines to err
	// not safe. This means that each Connect() will re-attempt to connect (this might be useful, too).
	connection := cfg.DBConnString()
	var err error
	db, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database connection: %v", err)
	}

	err = db.AutoMigrate(model.GetModels()...)
	if err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	pg, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database connection: %v", err)
	}

	// check connection
	if err = pg.Ping(); err != nil {
		return nil, err
	}

	// set up max connections to be somewhat longer that the server max (to allow for server system module connections)
	// the postgresql server max allowed connections is by default 100
	pg.SetMaxOpenConns(90)
	return db, nil
}
