package models

import (
    "fmt"
    "strconv"
    "time"
    
    "github.com/bryan/finance-tracker/internal/database"
    "github.com/bryan/finance-tracker/internal/validator"
)

type Transaction struct {
    ID              int       `json:"id"`
    Amount          float64   `json:"amount"`
    Description     string    `json:"description"`
    CategoryID      int       `json:"category_id"`
    CategoryName    string    `json:"category_name,omitempty"` // Used in joins
    CategoryType    string    `json:"category_type,omitempty"` // Used in joins
    TransactionDate time.Time `json:"transaction_date"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// TransactionFilter represents options for filtering transactions
type TransactionFilter struct {
    CategoryID      int
    CategoryType    string
    StartDate       time.Time
    EndDate         time.Time
    SortBy          string
    SortDirection   string
}

// Create adds a new transaction to the database
func (t *Transaction) Create() error {
    stmt := `
        INSERT INTO transactions (amount, description, category_id, transaction_date) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at`

    return database.DB.QueryRow(
        stmt, t.Amount, t.Description, t.CategoryID, t.TransactionDate,
    ).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

// Update updates an existing transaction in the database
func (t *Transaction) Update() error {
    stmt := `
        UPDATE transactions 
        SET amount = $1, description = $2, category_id = $3, transaction_date = $4, updated_at = CURRENT_TIMESTAMP
        WHERE id = $5
        RETURNING updated_at`

    return database.DB.QueryRow(
        stmt, t.Amount, t.Description, t.CategoryID, t.TransactionDate, t.ID,
    ).Scan(&t.UpdatedAt)
}

// Delete removes a transaction from the database
func (t *Transaction) Delete() error {
    stmt := `DELETE FROM transactions WHERE id = $1`
    _, err := database.DB.Exec(stmt, t.ID)
    return err
}

// GetTransactionByID retrieves a transaction by its ID
func GetTransactionByID(id int) (Transaction, error) {
    var transaction Transaction
    
    stmt := `
        SELECT t.id, t.amount, t.description, t.category_id, c.name, c.type, t.transaction_date, t.created_at, t.updated_at
        FROM transactions t
        JOIN categories c ON t.category_id = c.id
        WHERE t.id = $1`

    err := database.DB.QueryRow(stmt, id).Scan(
        &transaction.ID, 
        &transaction.Amount, 
        &transaction.Description, 
        &transaction.CategoryID,
        &transaction.CategoryName,
        &transaction.CategoryType,
        &transaction.TransactionDate, 
        &transaction.CreatedAt, 
        &transaction.UpdatedAt,
    )
    
    return transaction, err
}

// GetTransactions retrieves transactions with optional filtering
func GetTransactions(filter TransactionFilter) ([]Transaction, error) {
    // Start with the base query
    query := `
        SELECT t.id, t.amount, t.description, t.category_id, c.name, c.type, t.transaction_date, t.created_at, t.updated_at
        FROM transactions t
        JOIN categories c ON t.category_id = c.id
        WHERE 1=1`
    
    // Build args array for the query parameters
    args := []interface{}{}
    paramCount := 0

    // Add filter conditions if provided
    if filter.CategoryID > 0 {
        paramCount++
        query += fmt.Sprintf(" AND t.category_id = $%d", paramCount)
        args = append(args, filter.CategoryID)
    }

    if filter.CategoryType != "" {
        paramCount++
        query += fmt.Sprintf(" AND c.type = $%d", paramCount)
        args = append(args, filter.CategoryType)
    }

    if !filter.StartDate.IsZero() {
        paramCount++
        query += fmt.Sprintf(" AND t.transaction_date >= $%d", paramCount)
        args = append(args, filter.StartDate)
    }

    if !filter.EndDate.IsZero() {
        paramCount++
        query += fmt.Sprintf(" AND t.transaction_date <= $%d", paramCount)
        args = append(args, filter.EndDate)
    }

    // Add sorting
    sortBy := "t.transaction_date"
    if filter.SortBy != "" {
        switch filter.SortBy {
        case "amount":
            sortBy = "t.amount"
        case "category":
            sortBy = "c.name"
        case "date":
            sortBy = "t.transaction_date"
        }
    }
    
    sortDirection := "DESC"
    if filter.SortDirection == "ASC" {
        sortDirection = "ASC"
    }
    
    query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortDirection)

    // Execute the query
    rows, err := database.DB.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []Transaction

    for rows.Next() {
        var transaction Transaction
        if err := rows.Scan(
            &transaction.ID, 
            &transaction.Amount, 
            &transaction.Description, 
            &transaction.CategoryID,
            &transaction.CategoryName,
            &transaction.CategoryType,
            &transaction.TransactionDate, 
            &transaction.CreatedAt, 
            &transaction.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        transactions = append(transactions, transaction)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return transactions, nil
}

// GetSummary retrieves summary statistics for the transactions
func GetSummary(startDate, endDate time.Time) (map[string]float64, error) {
    summary := map[string]float64{
        "totalIncome":  0,
        "totalExpense": 0,
        "balance":      0,
    }

    // Query for total income and expenses
    stmt := `
        SELECT c.type, SUM(t.amount) as total
        FROM transactions t
        JOIN categories c ON t.category_id = c.id
        WHERE t.transaction_date BETWEEN $1 AND $2
        GROUP BY c.type`

    rows, err := database.DB.Query(stmt, startDate, endDate)
    if err != nil {
        return summary, err
    }
    defer rows.Close()

    for rows.Next() {
        var categoryType string
        var total float64

        err := rows.Scan(&categoryType, &total)
        if err != nil {
            return summary, err
        }

        if categoryType == "income" {
            summary["totalIncome"] = total
        } else if categoryType == "expense" {
            summary["totalExpense"] = total
        }
    }

    if err = rows.Err(); err != nil {
        return summary, err
    }

    summary["balance"] = summary["totalIncome"] - summary["totalExpense"]
    return summary, nil
}

// ValidateTransaction validates transaction data
func ValidateTransaction(v *validator.Validator, transaction *Transaction) {
    // Check amount is greater than zero
    v.Check(transaction.Amount > 0, "amount", "Amount must be greater than zero")
    
    // Check description length
    v.Check(validator.MaxLength(transaction.Description, 500), "description", "Description cannot exceed 500 characters")
    
    // Check category ID is valid
    v.Check(transaction.CategoryID > 0, "category_id", "Please select a valid category")
    
    // Check transaction date is not empty
    v.Check(!transaction.TransactionDate.IsZero(), "transaction_date", "Transaction date is required")
    
    // Check transaction date is not in the future
    v.Check(transaction.TransactionDate.Before(time.Now().AddDate(0, 0, 1)), "transaction_date", "Transaction date cannot be in the future")
}

// ParseTransactionForm parses the form data to create a Transaction object
func ParseTransactionForm(form map[string]string) (*Transaction, error) {
    transaction := &Transaction{}
    
    // Parse amount
    if form["amount"] != "" {
        amount, err := strconv.ParseFloat(form["amount"], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid amount format")
        }
        transaction.Amount = amount
    }
    
    // Parse description
    transaction.Description = form["description"]
    
    // Parse category ID
    if form["category_id"] != "" {
        categoryID, err := strconv.Atoi(form["category_id"])
        if err != nil {
            return nil, fmt.Errorf("invalid category ID format")
        }
        transaction.CategoryID = categoryID
    }
    
    // Parse transaction date
    if form["transaction_date"] != "" {
        date, err := time.Parse("2006-01-02", form["transaction_date"])
        if err != nil {
            return nil, fmt.Errorf("invalid date format. Use YYYY-MM-DD")
        }
        transaction.TransactionDate = date
    }
    
    // Parse ID for updates
    if form["id"] != "" {
        id, err := strconv.Atoi(form["id"])
        if err != nil {
            return nil, fmt.Errorf("invalid ID format")
        }
        transaction.ID = id
    }
    
    return transaction, nil
}