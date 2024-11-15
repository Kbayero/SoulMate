package database

import (
	"log"
	"os"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbInstance *Database
	dbOnce     sync.Once
)

type Database struct {
	db     *gorm.DB
	locker sync.RWMutex
}

func (d *Database) Migrate(data interface{}) error {
	return d.db.AutoMigrate(data)
}

func (d *Database) Create(data interface{}) error {
	d.locker.Lock()
	defer d.locker.Unlock()
	return d.db.Create(data).Error
}

func (d *Database) GetAll(data interface{}) error {
	d.locker.Lock()
	defer d.locker.Unlock()
	if err := d.db.Find(data).Error; err != nil {
		return err
	}
	return nil
}

func GetDB() *Database {
	dbOnce.Do(func() {
		path := "responses.sqlite3"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			file, err := os.Create(path)
			if err != nil {
				log.Fatalf("error creating database file: %v", err)
			}
			file.Close()
		}

		conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Fatalf("error connecting with database: %v", err)
		}

		dbInstance = &Database{db: conn}

	})

	return dbInstance
}
