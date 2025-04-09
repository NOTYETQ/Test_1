CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('income', 'expense')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert default categories
INSERT INTO categories (name, type) VALUES
('Salary', 'income'),
('Freelance', 'income'),
('Investments', 'income'),
('Other Income', 'income'),
('Food & Dining', 'expense'),
('Rent/Mortgage', 'expense'),
('Utilities', 'expense'),
('Transportation', 'expense'),
('Entertainment', 'expense'),
('Healthcare', 'expense'),
('Shopping', 'expense'),
('Education', 'expense'),
('Other Expense', 'expense');