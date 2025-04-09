package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
    "log"
    "fmt"
)

var templateCache = make(map[string]*template.Template)

// InitTemplates pre-loads and caches all templates
// InitTemplates pre-loads and caches all templates
func InitTemplates() error {
    // Define template paths
    baseLayout := filepath.Join("templates", "layout", "base.html")
    
    log.Printf("Loading base template from: %s", baseLayout)
    
    // Get all regular templates
    pages, err := filepath.Glob(filepath.Join("templates", "*.html"))
    if err != nil {
        return fmt.Errorf("template glob error: %v", err)
    }
    
    log.Printf("Found %d page templates: %v", len(pages), pages)
    
    // Parse all templates with the base layout
    for _, page := range pages {
        name := filepath.Base(page)
        
        log.Printf("Parsing template: %s", name)
        
        // Parse base layout first, then the page
        tmpl, err := template.ParseFiles(baseLayout, page)
        if err != nil {
            return fmt.Errorf("error parsing template %s: %v", name, err)
        }
        
        // Log defined templates
        log.Printf("Template %s contains definitions:", name)
        for _, t := range tmpl.Templates() {
            log.Printf("  - %s", t.Name())
        }
        
        templateCache[name] = tmpl
        log.Printf("Cached template: %s", name)
    }
    
    return nil
}

// render executes a template with provided data
// render executes a template with provided data
func render(w http.ResponseWriter, tmpl string, data interface{}) {
    // Get template from cache
    t, ok := templateCache[tmpl]
    if !ok {
        log.Printf("Template not found in cache: %s", tmpl)
        
        // Try parsing on the fly as fallback
        baseLayout := filepath.Join("templates", "layout", "base.html")
        page := filepath.Join("templates", tmpl)
        
        log.Printf("Attempting to parse %s with %s", page, baseLayout)
        
        var err error
        t, err = template.ParseFiles(baseLayout, page)
        if err != nil {
            http.Error(w, "Template not found: "+err.Error(), http.StatusInternalServerError)
            log.Printf("Error parsing template: %v", err)
            return
        }
        
        templateCache[tmpl] = t
        log.Printf("Parsed template on demand: %s", tmpl)
    }
    
    // Log template data for debugging
    log.Printf("Rendering template %s with data: %+v", tmpl, data)
    
    // Execute the template
    err := t.ExecuteTemplate(w, "base", data)
    if err != nil {
        http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
        log.Printf("Template execution error: %v", err)
    }
}