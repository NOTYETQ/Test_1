package main

import (
    "fmt"
    "log"
    "net/http"
    
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/gorilla/mux"
    
    "github.com/bryan/finance-tracker/internal/database"
    "github.com/bryan/finance-tracker/internal/handlers"
    "github.com/bryan/finance-tracker/internal/middleware"
)

const (
    migrationURL = "file://migrations"
    port         = ":4000"
)

func main() {
    // Connect to the database
    if err := database.Connect(); err != nil {
        log.Fatalf("Database connection failed: %v", err)
    }
    defer database.Close()
    
    // Initialize templates
    if err := handlers.InitTemplates(); err != nil {
        log.Fatalf("Error initializing templates: %v", err)
    }
    
    // Run migrations
    if err := runMigrations(); err != nil {
        log.Fatalf("Error running migrations: %v", err)
    }
    
    // Create router
    r := mux.NewRouter()
    r.Use(middleware.LoggingMiddleware)
    
    // Static files
    staticDir := http.Dir("./static")
    staticHandler := http.StripPrefix("/static/", http.FileServer(staticDir))
    r.PathPrefix("/static/").Handler(staticHandler)
    
    // Dashboard routes
    r.HandleFunc("/", handlers.DashboardHandler).Methods("GET")
    
    // Transaction routes
    r.HandleFunc("/transactions", handlers.ListTransactionsHandler).Methods("GET")
    r.HandleFunc("/transactions/new", handlers.GetTransactionFormHandler).Methods("GET")
    r.HandleFunc("/transactions", handlers.CreateTransactionHandler).Methods("POST")
    r.HandleFunc("/transactions/{id:[0-9]+}/edit", handlers.GetTransactionEditHandler).Methods("GET")
    r.HandleFunc("/transactions/{id:[0-9]+}", handlers.UpdateTransactionHandler).Methods("POST")
    r.HandleFunc("/transactions/{id:[0-9]+}/delete", handlers.DeleteTransactionHandler).Methods("POST")

    
    // Start server
    fmt.Printf("Server starting on port %s\n", port)
    log.Fatal(http.ListenAndServe(port, r))
}

func runMigrations() error {
    db := database.DB
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("could not create migration driver: %v", err)
    }
    
    m, err := migrate.NewWithDatabaseInstance(
        migrationURL,
        "postgres", driver)
    if err != nil {
        return fmt.Errorf("could not create migrate instance: %v", err)
    }
    
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("could not run migrations: %v", err)
    }
    
    return nil
}