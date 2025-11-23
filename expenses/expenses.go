package expenses

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Expense struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
}

const (
	ErrEmptyDescription = "description cannot be empty"
	ErrInvalidAmount    = "amount must be positive"
	ErrInvalidDate      = "invalid date format"
)

func generateID(expenses []Expense) int {
	maxID := 0
	for _, expense := range expenses {
		if expense.ID > maxID {
			maxID = expense.ID
		}
	}
	return maxID + 1
}

func loadExpenses() ([]Expense, error) {
	if _, err := os.Stat("expenses.json"); os.IsNotExist(err) {
		return []Expense{}, nil
	}

	file, err := os.ReadFile("expenses.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read file - %w", err)
	}

	if len(file) == 0 {
		return []Expense{}, nil
	}

	var expenses []Expense
	err = json.Unmarshal(file, &expenses)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshul - %w", err)
	}
	return expenses, err
}

func saveExpenses(expenses []Expense) error {
	data, err := json.MarshalIndent(expenses, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal - %w", err)
	}

	err = os.WriteFile("expenses.json", data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file - %w", err)
	}
	return nil
}

func withTask(operation func(expenses *[]Expense) error) error {
	expenses, err := loadExpenses()
	if err != nil {
		return fmt.Errorf("failed to load expenses - %w", err)
	}

	if err := operation(&expenses); err != nil {
		return err
	}

	return saveExpenses(expenses)
}

func AddExpense(description string, category string, amount float64) error {
	if description == "" {
		return fmt.Errorf(ErrEmptyDescription)
	}
	if amount <= 0 {
		return fmt.Errorf(ErrInvalidAmount)
	}
	return withTask(func(expenses *[]Expense) error {
		expense := Expense{
			ID:          generateID(*expenses),
			Date:        time.Now(),
			Description: description,
			Category:    category,
			Amount:      amount,
		}

		*expenses = append(*expenses, expense)
		fmt.Printf("Expense added successfully (ID: %d)\n", expense.ID)

		return nil
	})
}

func DeleteExpense(id int) error {
	return withTask(func(expenses *[]Expense) error {
		for i := 0; i < len(*expenses); i++ {
			if (*expenses)[i].ID == id {
				*expenses = append((*expenses)[:i], (*expenses)[i+1:]...)
				fmt.Printf("Expense deleted successfully (ID: %d)\n", id)
				return nil
			}
		}
		return fmt.Errorf("expense with ID %d not found", id)
	})
}

func ListExpense() error {
	return withTask(func(expenses *[]Expense) error {
		if len(*expenses) == 0 {
			fmt.Println("No expenses found")
			return nil
		}

		fmt.Printf("%-4s %-10s %-20s %-15s %s\n", "ID", "Date", "Description", "Category", "Amount")
		for _, expense := range *expenses {
			fmt.Printf("%-4d %-10s %-20s %-15s $%.2f\n",
				expense.ID,
				expense.Date.Format("2006-01-02"),
				expense.Description,
				expense.Category,
				expense.Amount)
		}
		return nil
	})
}

func Summary(month int) (float64, error) {
	var sum float64
	err := withTask(func(expenses *[]Expense) error {
		for _, expense := range *expenses {
			if month == 0 || expense.Date.Month() == time.Month(month) {
				sum += expense.Amount
			}
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to calculate summary: %w", err)
	}
	return sum, nil
}

func CategorySummary(category string) (float64, error) {
	var sum float64
	err := withTask(func(expenses *[]Expense) error {
		for _, expense := range *expenses {
			if strings.EqualFold(expense.Category, category) {
				sum += expense.Amount
			}
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to calculate category summary: %w", err)
	}
	return sum, nil
}

func GetExpensesByCategory(category string) ([]Expense, error) {
	var filtered []Expense
	err := withTask(func(expenses *[]Expense) error {
		for _, expense := range *expenses {
			if strings.EqualFold(expense.Category, category) {
				filtered = append(filtered, expense)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get expenses by category: %w", err)
	}
	return filtered, nil
}
