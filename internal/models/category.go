package models

import (
    "time"
    
    "github.com/bryan/finance-tracker/internal/database"
    "github.com/bryan/finance-tracker/internal/validator"
)

type Category struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Type      string    `json:"type"` // 'income' or 'expense'
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Create adds a new category to the database
func (c *Category) Create() error {
    stmt := `
        INSERT INTO categories (name, type) 
        VALUES ($1, $2)
        RETURNING id, created_at, updated_at`

    return database.DB.QueryRow(stmt, c.Name, c.Type).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

// GetAllCategories retrieves all categories from the database
func GetAllCategories() ([]Category, error) {
    stmt := `
        SELECT id, name, type, created_at, updated_at 
        FROM categories 
        ORDER BY name`

    rows, err := database.DB.Query(stmt)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []Category

    for rows.Next() {
        var category Category
        if err := rows.Scan(&category.ID, &category.Name, &category.Type, &category.CreatedAt, &category.UpdatedAt); err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return categories, nil
}

// GetCategoriesByType retrieves categories filtered by type
func GetCategoriesByType(categoryType string) ([]Category, error) {
    stmt := `
        SELECT id, name, type, created_at, updated_at 
        FROM categories 
        WHERE type = $1
        ORDER BY name`

    rows, err := database.DB.Query(stmt, categoryType)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []Category

    for rows.Next() {
        var category Category
        if err := rows.Scan(&category.ID, &category.Name, &category.Type, &category.CreatedAt, &category.UpdatedAt); err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return categories, nil
}

// GetCategoryByID retrieves a category by its ID
func GetCategoryByID(id int) (Category, error) {
    var category Category
    
    stmt := `
        SELECT id, name, type, created_at, updated_at 
        FROM categories 
        WHERE id = $1`

    err := database.DB.QueryRow(stmt, id).Scan(
        &category.ID, &category.Name, &category.Type, &category.CreatedAt, &category.UpdatedAt)
    
    return category, err
}

// ValidateCategory validates category data
func ValidateCategory(v *validator.Validator, category *Category) {
    v.Check(validator.NotBlank(category.Name), "name", "Category name is required")
    v.Check(validator.MaxLength(category.Name, 100), "name", "Category name cannot exceed 100 characters")
    
    v.Check(validator.NotBlank(category.Type), "type", "Category type is required")
    v.Check(category.Type == "income" || category.Type == "expense", "type", "Category type must be either 'income' or 'expense'")
}