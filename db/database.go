package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/pclubiitk/puppylove_tags/models"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found")
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")

	fmt.Printf("DB Config: host=%s, user=%s, port=%s, dbname=%s\n", host, user, port, dbName)

	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		host, user, password, dbName, port,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		fmt.Println("DB connection error:", err)
		panic(err)
	}

	// Auto-migrate the models
	if err := db.AutoMigrate(&models.UserTag{}); err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database!")
	DB = db
	return db
}
