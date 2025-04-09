package handlers

import (
    "net/http"
    "strconv"
    "time"
    
    "github.com/gorilla/mux"
    
    "github.com/bryan/finance-tracker/internal/models"
    "github.com/bryan/finance-tracker/internal/validator"
)

// ListTransactionsHandler displays a list of all transactions
func ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters for filtering
    filter := parseTransactionFilter(r)
    
    // Get transactions based on filter
    transactions, err := models.GetTransactions(filter)
    if err != nil {
        http.Error(w, "Error fetching transactions: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Get categories for the filter form
    categories, err := models.GetAllCategories()
    if err != nil {
        http.Error(w, "Error fetching categories: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Calculate summary for the current date range
    summary, err := models.GetSummary(filter.StartDate, filter.EndDate)
    if err != nil {
        http.Error(w, "Error calculating summary: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    data := struct {
        Transactions []models.Transaction
        Categories   []models.Category
        Filter       models.TransactionFilter
        Summary      map[string]float64
    }{
        Transactions: transactions,
        Categories:   categories,
        Filter:       filter,
        Summary:      summary,
    }
    
    render(w, "transaction_list.html", data)
}

// GetTransactionFormHandler displays the form to add a new transaction
func GetTransactionFormHandler(w http.ResponseWriter, r *http.Request) {
    categories, err := models.GetAllCategories()
    if err != nil {
        http.Error(w, "Error fetching categories: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Pre-populate with today's date
    transaction := models.Transaction{
        TransactionDate: time.Now(),
    }
    
    data := struct {
        Transaction models.Transaction
        Categories  []models.Category
        Validator   *validator.Validator
    }{
        Transaction: transaction,
        Categories:  categories,
        Validator:   validator.NewValidator(),
    }
    
    render(w, "transaction_form.html", data)
}

// CreateTransactionHandler handles the submission of a new transaction
func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
    // Parse form data
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
        return
    }
    
    // Prepare form data map
    formData := make(map[string]string)
    formData["amount"] = r.FormValue("amount")
    formData["description"] = r.FormValue("description")
    formData["category_id"] = r.FormValue("category_id")
    formData["transaction_date"] = r.FormValue("transaction_date")
    
    // Parse transaction from form data
    transaction, err := models.ParseTransactionForm(formData)
    if err != nil {
        http.Error(w, "Error parsing transaction: "+err.Error(), http.StatusBadRequest)
        return
    }
    
    // Validate transaction
    v := validator.NewValidator()
    models.ValidateTransaction(v, transaction)
    
    // If validation fails, re-render the form with errors
    if !v.ValidData() {
        // Get categories for form
        categories, err := models.GetAllCategories()
        if err != nil {
            http.Error(w, "Error fetching categories: "+err.Error(), http.StatusInternalServerError)
            return
        }
        
        data := struct {
            Transaction models.Transaction
            Categories  []models.Category
            Validator   *validator.Validator
        }{
            Transaction: *transaction,
            Categories:  categories,
            Validator:   v,
        }
        
        render(w, "transaction_form.html", data)
        return
    }
    
    // Save transaction to database
    if err := transaction.Create(); err != nil {
        http.Error(w, "Error creating transaction: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Redirect to transaction list
    http.Redirect(w, r, "/transactions", http.StatusSeeOther)
}

// GetTransactionEditHandler displays the form to edit a transaction
func GetTransactionEditHandler(w http.ResponseWriter, r *http.Request) {
    // Extract transaction ID from URL
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
        return
    }
    
    // Get transaction by ID
    transaction, err := models.GetTransactionByID(id)
    if err != nil {
        http.Error(w, "Error fetching transaction: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Get categories for form
    categories, err := models.GetAllCategories()
    if err != nil {
        http.Error(w, "Error fetching categories: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    data := struct {
        Transaction models.Transaction
        Categories  []models.Category
        Validator   *validator.Validator
    }{
        Transaction: transaction,
        Categories:  categories,
        Validator:   validator.NewValidator(),
    }
    
    render(w, "transaction_edit.html", data)
}

// UpdateTransactionHandler handles the submission of an updated transaction
func UpdateTransactionHandler(w http.ResponseWriter, r *http.Request) {
    // Extract transaction ID from URL
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
        return
    }
    
    // Parse form data
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
        return
    }
    
    // Prepare form data map
    formData := make(map[string]string)
    formData["id"] = strconv.Itoa(id)
    formData["amount"] = r.FormValue("amount")
    formData["description"] = r.FormValue("description")
    formData["category_id"] = r.FormValue("category_id")
    formData["transaction_date"] = r.FormValue("transaction_date")
    
    // Parse transaction from form data
    transaction, err := models.ParseTransactionForm(formData)
    if err != nil {
        http.Error(w, "Error parsing transaction: "+err.Error(), http.StatusBadRequest)
        return
    }
    
    // Validate transaction
    v := validator.NewValidator()
    models.ValidateTransaction(v, transaction)
    
    // If validation fails, re-render the form with errors
    if !v.ValidData() {
        // Get categories for form
        categories, err := models.GetAllCategories()
        if err != nil {
            http.Error(w, "Error fetching categories: "+err.Error(), http.StatusInternalServerError)
            return
        }
        
        data := struct {
            Transaction models.Transaction
            Categories  []models.Category
            Validator   *validator.Validator
        }{
            Transaction: *transaction,
            Categories:  categories,
            Validator:   v,
        }
        
        render(w, "transaction_edit.html", data)
        return
    }
    
    // Update transaction in database
    if err := transaction.Update(); err != nil {
        http.Error(w, "Error updating transaction: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Redirect to transaction list
    http.Redirect(w, r, "/transactions", http.StatusSeeOther)
}

// DeleteTransactionHandler handles the deletion of a transaction
func DeleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
    // Extract transaction ID from URL
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
        return
    }
    
    // Create transaction object with ID
    transaction := &models.Transaction{ID: id}
    
    // Delete transaction from database
    if err := transaction.Delete(); err != nil {
        http.Error(w, "Error deleting transaction: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Redirect to transaction list
    http.Redirect(w, r, "/transactions", http.StatusSeeOther)
}

// Helper function to parse transaction filter from request
func parseTransactionFilter(r *http.Request) models.TransactionFilter {
    filter := models.TransactionFilter{}
    
    // Parse category ID filter
    if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
        id, err := strconv.Atoi(categoryID)
        if err == nil && id > 0 {
            filter.CategoryID = id
        }
    }
    
    // Parse category type filter
    if categoryType := r.URL.Query().Get("type"); categoryType == "income" || categoryType == "expense" {
        filter.CategoryType = categoryType
    }
    
    // Parse date range filters
    if startDate := r.URL.Query().Get("start_date"); startDate != "" {
        date, err := time.Parse("2006-01-02", startDate)
        if err == nil {
            filter.StartDate = date
        }
    } else {
        // Default to first day of current month
        now := time.Now()
        filter.StartDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
    }
    
    if endDate := r.URL.Query().Get("end_date"); endDate != "" {
        date, err := time.Parse("2006-01-02", endDate)
        if err == nil {
            filter.EndDate = date
        }
    } else {
        // Default to last day of current month
        now := time.Now()
        filter.EndDate = time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())
    }
    
    // Parse sorting options
    if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
        filter.SortBy = sortBy
    }
    
    if sortDir := r.URL.Query().Get("sort_dir"); sortDir == "ASC" {
        filter.SortDirection = "ASC"
    } else {
        filter.SortDirection = "DESC"
    }
    
    return filter
}