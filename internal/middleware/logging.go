package middleware

import (
    "log"
    "net/http"
    "time"
)

// LoggingMiddleware logs information about each HTTP request
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Record the start time
        start := time.Now()
        
        // Log the request details
        log.Printf("REQUEST: %s %s FROM %s", r.Method, r.URL.Path, r.RemoteAddr)

        // Call the next handler
        next.ServeHTTP(w, r)

        // Calculate and log the response time
        duration := time.Since(start)
        log.Printf("RESPONSE: %s %s - completed in %v", r.Method, r.URL.Path, duration)
    })
}