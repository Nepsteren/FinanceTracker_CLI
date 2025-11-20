package expenses

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Expense struct {
	Id          int     `json:"id"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
}

func expectArgs(count int, args []string) error {
	if count != len(args) {
		return fmt.Errorf("failed args. Expected - %d, have - %d", count, len(args))
	}
	return nil
}

func generateId(expenses []Expense) int {
	maxId := 0
	for _, expense := range expenses {
		if expense.Id > maxId {
			maxId = expense.Id
		}
	}
	return maxId + 1
}

func loadExpenses() ([]Expense, error) {
	file, err := os.ReadFile("expenses.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read file - %w", err)
	}
	var expenses []Expense
	err = json.Unmarshal(file, &expenses)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshul - %w", err)
	}
	return expenses, err
}

func marshalJson(expenses []Expense) error {
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
	if _, err := os.Stat("expenses.json"); os.IsNotExist(err) {
		if err := os.WriteFile("expenses.json", []byte("[]"), 0644); err != nil {
			return fmt.Errorf("failed to create tasks file: %w", err)
		}
	}
	expenses, err := loadExpenses()
	if err != nil {
		return fmt.Errorf("failed to load expenses - %w", err)
	}

	err = operation(&expenses)
	if err != nil {
		return fmt.Errorf("failed to operation - %w", err)
	}

	return marshalJson(expenses)
}

func AddExpense() error {
	// if description == "" {
	// 	return fmt.Errorf("task description cannot be empty")
	// }
	return withTask(func(expenses *[]Expense) error {
		expense := Expense{
			Id:          generateId(*expenses),
			Date:        time.Now().Format("2006-01-02 15:04:05"),
			Description: "",
			Amount:      2,
		}
		*expenses = append(*expenses, expense)
		fmt.Printf("Expense added successfully (ID: %d)\n", expense.Id)
		return nil
	})
}
