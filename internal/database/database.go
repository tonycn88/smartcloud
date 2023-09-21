package database

import (
	"smartcloud/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Init_database(dbconf config.Database) {
	if dbconf.Datatype == "sqlite3" {
		var err error
		Db, err = gorm.Open(sqlite.Open(dbconf.Url), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
	}

}
