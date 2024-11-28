package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {

	load_env := godotenv.Load()
	if load_env != nil {
		log.Fatal("Error loading .env file")
	}

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/time_api", mysqlUser, mysqlPassword, mysqlHost)

	// Connect to MySQL
	var err error
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// Ensure connection is valid
	if err := db.Ping(); err != nil {
		log.Fatalf("MySQL connection error: %v", err)
	}
	fmt.Println("Connected to MySQL!")

	// Set up Gin router
	r := gin.Default()

	// Define the /current-time endpoint
	r.GET("/current-time", func(c *gin.Context) {
		// Get Toronto time
		loc, err := time.LoadLocation("America/Toronto")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load timezone"})
			return
		}
		torontoTime := time.Now().In(loc)

		// Log time to database
		_, err = db.Exec("INSERT INTO time_log (timestamp) VALUES (?)", torontoTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log time to database"})
			return
		}

		// Respond with current Toronto time
		c.JSON(http.StatusOK, gin.H{"current_time": torontoTime.Format(time.RFC3339)})
	})

	// Start the server
	r.Run(":8080")
}
