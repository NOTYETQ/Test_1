package database

import (
    "database/sql"
    "fmt"
    "log"
    "regexp"
    "strings"
    "time"

    _ "github.com/lib/pq"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "finance"  
    password = "finance"
    dbname   = "finance"
)

var DB *sql.DB

// Connect establishes a connection to the PostgreSQL database
func Connect() error {
    // Connection string
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    
    // Open connection
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("error opening database: %v", err)
    }
    
    // Test the connection
    if err = db.Ping(); err != nil {
        return fmt.Errorf("error connecting to the database: %v", err)
    }
    
    // Set connection parameters
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    DB = db
    log.Println("Successfully connected to the database")
    return nil
}

// Close closes the database connection
func Close() {
    if DB != nil {
        DB.Close()
    }
}

// Validation helpers
func ValidateRequired(value string) bool {
    return strings.TrimSpace(value) != ""
}

func ValidateMaxLength(value string, maxLength int) bool {
    return len(value) <= maxLength
}

func ValidateEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9.%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return emailRegex.MatchString(email)
}

// Truncate helper to limit string length
func Truncate(value string, maxLength int) string {
    if len(value) > maxLength {
        return value[:maxLength]
    }
    return value
}