package expenses

import (
	"encoding/json"

	"fmt"
	"os"
	"strconv"
	"time"
)

type Expense struct {
	Id          int     `json:"id"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
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

func readFile() (float64, error) {
	fileName := "budget.txt"
	data, err := os.ReadFile(fileName)
	if err != nil {
		return 0, fmt.Errorf("failed to read file - %w", err)
	}
	budget, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse budget - %w", err)
	}
	return budget, nil
}

func checkBudget() error {
	_, err := os.Stat("budget.txt")
	if os.IsNotExist(err) {
		return nil
	}
	budget, err := readFile()
	if err != nil {
		return fmt.Errorf("failed to read budget - %w", err)
	}

	summary, err := Summary()
	if err != nil {
		return fmt.Errorf("failed to count all amount - %w", err)
	}

	if summary > budget {
		fmt.Printf("You went over budget! Budget - %0.2f. Summary - %0.2f\n", budget, summary)
	}
	return nil
}

func AddExpense(description string, amount float64) error {
	if description == "" {
		return fmt.Errorf("task description cannot be empty")
	}
	return withTask(func(expenses *[]Expense) error {
		expense := Expense{
			Id:          generateId(*expenses),
			Date:        time.Now().Format("2006-01-02"),
			Description: description,
			Amount:      amount,
		}
		err := checkBudget()
		if err != nil {
			return fmt.Errorf("failed to check budget - %w", err)
		}
		*expenses = append(*expenses, expense)
		fmt.Printf("Expense added successfully (ID: %d)\n", expense.Id)

		return nil
	})
}

func DeleteExpense(id int) error {
	return withTask(func(expenses *[]Expense) error {
		for i := 0; i < len(*expenses); i++ {
			if (*expenses)[i].Id == id {
				*expenses = append((*expenses)[:i], (*expenses)[i+1:]...)
				fmt.Printf("Expense deleted successfully (ID: %d)\n", id)
				return nil
			}
		}
		return fmt.Errorf("task with ID %d not found", id)
	})
}

func ListExpense() error {
	return withTask(func(expenses *[]Expense) error {
		// fmt.Println("ID Date    Description  Amount")
		for i := range *expenses {
			fmt.Println((*expenses)[i])
		}
		return nil
	})
}

func Summary() (float64, error) {
	var sum float64
	err := withTask(func(expenses *[]Expense) error {
		for i := range *expenses {
			sum += (*expenses)[i].Amount
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return sum, nil
}
