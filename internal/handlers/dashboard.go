package handlers

import (
    "net/http"
    "time"
    
    "github.com/bryan/finance-tracker/internal/models"
)

// DashboardHandler displays the dashboard page with summary information
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
    // Get current month's date range
    now := time.Now()
    startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
    endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())
    
    // Get transactions for current month
    filter := models.TransactionFilter{
        StartDate: startOfMonth,
        EndDate:   endOfMonth,
    }
    
    transactions, err := models.GetTransactions(filter)
    if err != nil {
        http.Error(w, "Error fetching transactions: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Calculate summary for current month
    summary, err := models.GetSummary(startOfMonth, endOfMonth)
    if err != nil {
        http.Error(w, "Error calculating summary: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Get recent transactions (limited to 5)
    recentFilter := models.TransactionFilter{
        SortBy:        "date",
        SortDirection: "DESC",
    }
    
    recentTransactions, err := models.GetTransactions(recentFilter)
    if err != nil {
        http.Error(w, "Error fetching recent transactions: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Limit to 5 transactions
    if len(recentTransactions) > 5 {
        recentTransactions = recentTransactions[:5]
    }
    
    data := struct {
        Summary           map[string]float64
        Transactions      []models.Transaction
        RecentTransactions []models.Transaction
        CurrentMonth      string
    }{
        Summary:           summary,
        Transactions:      transactions,
        RecentTransactions: recentTransactions,
        CurrentMonth:      now.Format("January 2006"),
    }
    
    render(w, "dashboard.html", data)
}