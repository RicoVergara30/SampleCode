package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

// connection to env
func ConnectDB() error {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	data, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	Database = data

	return nil

}

// package database

// import (
// 	"log"
// 	"os"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// 	"gorm.io/gorm/logger"
// 	// "gorm.io/gorm/logger"
// )

// type Database struct {
// 	Databasedb *gorm.DB
// }

// var Data Database

// func ConnectDb() {
// 	databasedns := "host=localhost user=fdsap-v.rico password=3030rico3030 dbname=tryandtry port=5432 sslmode=disable"
// 	databasedb, err := gorm.Open(postgres.Open(databasedns), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("Failed to connect in bakend \n", err.Error())
// 		os.Exit(2)
// 	}

// 	databasedb.Logger = logger.Default.LogMode(logger.Info)
// 	log.Println("Running migration for backend database")
// 	// if err := databasedb.AutoMigrate(&models.AddCart{}); err != nil {
// 	// 	log.Fatalf("Error running migration for backend database: %v", err)
// 	// }

// 	// log.Println("Running Migration")
// 	// Data.BackendDb.AutoMigrate(&models.AddCart{})
// 	// Data.OasisDb.AutoMigrate()

// 	Data = Database{
// 		Databasedb: databasedb,
// 	}
// }
