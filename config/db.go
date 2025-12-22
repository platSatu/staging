// package config

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv" // <--- harus ditambahkan
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// func InitDB() *gorm.DB {
// 	// Load .env
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("No .env file found, using system environment variables")
// 	}

// 	// Gunakan nama variabel yang sesuai dengan .env
// 	user := os.Getenv("DB_USER")
// 	password := os.Getenv("DB_PASS") // sebelumnya salah tulis DB_PASSWORD
// 	host := os.Getenv("DB_HOST")
// 	dbName := os.Getenv("DB_NAME")

// 	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		user,
// 		password,
// 		host,
// 		dbName,
// 	)

// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("Failed to connect to database: ", err)
// 	}

//		return db
//	}
package config

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s",
		user,
		password,
		host,
		dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true, // 🔥 penting untuk traffic tinggi
	})
	if err != nil {
		return nil, err
	}

	// 🔥 DB connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(50) // sesuaikan CPU & DB
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
