//the categories will not be implimpent since i decided its best to just have default vaules
// i still will leave it incase i decide to change my mind but for now my focus was on the transactions

package handlers

import (
    "net/http"
    
    "github.com/bryan/finance-tracker/internal/models"
    "github.com/bryan/finance-tracker/internal/validator"
)

// GetCategoryFormHandler displays the form to add a new category
func GetCategoryFormHandler(w http.ResponseWriter, r *http.Request) {
    data := struct {
        Category  models.Category
        Validator *validator.Validator
    }{
        Validator: validator.NewValidator(),
    }
    
    render(w, "category_form.html", data)
}

// CreateCategoryHandler handles the submission of a new category
func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
    // Parse form data
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
        return
    }
    
    // Create category from form data
    category := &models.Category{
        Name: r.FormValue("name"),
        Type: r.FormValue("type"),
    }
    
    // Validate category
    v := validator.NewValidator()
    models.ValidateCategory(v, category)
    
    // If validation fails, re-render the form with errors
    if !v.ValidData() {
        data := struct {
            Category  models.Category
            Validator *validator.Validator
        }{
            Category:  *category,
            Validator: v,
        }
        
        render(w, "category_form.html", data)
        return
    }
    
    // Save category to database
    if err := category.Create(); err != nil {
        http.Error(w, "Error creating category: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Redirect to transaction form
    http.Redirect(w, r, "/transactions/new", http.StatusSeeOther)
}