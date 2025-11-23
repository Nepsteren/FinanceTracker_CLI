package commands

import (
	"bufio"
	"finTrackCLI/expenses"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func greeting() {
	fmt.Println("Вас приветствует - FinanceTrackerCLI")
	fmt.Println("Если вам нужна помощь введите команду: help")
}

func Start() error {
	greeting()
	err := userInput()
	if err != nil {
		return fmt.Errorf("failed to start - %w", err)
	}
	return nil
}

func userInput() error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			err := switchCommand(scanner.Text())
			if err != nil {
				return fmt.Errorf("failed to user input - %w", err)
			}
		}
	}
}

type CommandHandler func(args []string) error

var commands = map[string]CommandHandler{
	"help":       handleHelp,
	"exit":       handleExit,
	"add":        handleAdd,
	"delete":     handleDelete,
	"list":       handleList,
	"set-budget": handleSetBudget,
	"summary":    handleSummary,
}

func parseArgs(args []string) map[string]string {
	result := make(map[string]string)

	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--") {
			flag := strings.TrimPrefix(args[i], "--")
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				result[flag] = args[i+1]
				i++
			} else {
				result[flag] = ""
			}
		}
	}

	return result
}

func handleAdd(args []string) error {
	parsedArgs := parseArgs(args)

	description, ok := parsedArgs["description"]
	if !ok || description == "" {
		return fmt.Errorf("missing or empty --description")
	}

	amountStr, ok := parsedArgs["amount"]
	if !ok {
		return fmt.Errorf("missing --amount")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		return fmt.Errorf("invalid amount: %s", amountStr)
	}

	category := parsedArgs["category"]

	if err := expenses.AddExpense(description, category, amount); err != nil {
		return fmt.Errorf("failed to add expense: %w", err)
	}

	return checkBudget()
}

func handleHelp(args []string) error {
	help()
	return nil
}

func handleExit(args []string) error {
	exit()
	return nil
}

func handleDelete(args []string) error {
	parsedArgs := parseArgs(args)

	idStr, ok := parsedArgs["id"]
	if !ok {
		return fmt.Errorf("missing --id")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id: %s", idStr)
	}

	return expenses.DeleteExpense(id)
}

func handleList(args []string) error {
	return expenses.ListExpense()
}

func handleSetBudget(args []string) error {
	parsedArgs := parseArgs(args)

	amountStr, ok := parsedArgs["amount"]
	if !ok {
		return fmt.Errorf("missing --amount")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		return fmt.Errorf("invalid amount: %s", amountStr)
	}

	return saveBudget(amount)
}

func handleSummary(args []string) error {
	parsedArgs := parseArgs(args)

	if category, hasCategory := parsedArgs["category"]; hasCategory {
		return handleCategorySummary(category)
	}

	if monthStr, hasMonth := parsedArgs["month"]; hasMonth {
		return handleMonthSummary(monthStr)
	}

	return handleTotalSummary()
}

func handleCategorySummary(category string) error {
	if category == "" {
		return fmt.Errorf("category cannot be empty")
	}
	sum, err := expenses.CategorySummary(category)
	if err != nil {
		return fmt.Errorf("failed to get category summary: %w", err)
	}

	categoryExpenses, err := expenses.GetExpensesByCategory(category)
	if err != nil {
		return fmt.Errorf("failed to get category expenses: %w", err)
	}

	if len(categoryExpenses) == 0 {
		fmt.Printf("No expenses found for category: %s\n", category)
		return nil
	}

	fmt.Printf("Expenses for category '%s':\n", category)
	fmt.Printf("%-4s %-10s %-20s %-15s %s\n", "ID", "Date", "Description", "Category", "Amount")
	fmt.Println(strings.Repeat("-", 70))

	for _, expense := range categoryExpenses {
		fmt.Printf("%-4d %-10s %-20s %-15s $%.2f\n",
			expense.ID,
			expense.Date.Format("2006-01-02"),
			expense.Description,
			expense.Category,
			expense.Amount)
	}

	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("Total for category '%s': $%.2f (%d expenses)\n",
		category, sum, len(categoryExpenses))

	return nil
}

func handleMonthSummary(monthStr string) error {
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return fmt.Errorf("invalid month: %s", monthStr)
	}

	sum, err := expenses.Summary(month)
	if err != nil {
		return err
	}

	fmt.Printf("Total expenses for month %d: $%.2f\n", month, sum)
	return nil
}

func handleTotalSummary() error {
	sum, err := expenses.Summary(0)
	if err != nil {
		return err
	}

	fmt.Printf("Total expenses: $%.2f\n", sum)
	return nil
}

func saveBudget(budget float64) error {
	fileName := "budget.txt"
	data := fmt.Sprintf("%.2f", budget)
	err := os.WriteFile(fileName, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("failed to write budget to file - %w", err)
	}
	return nil
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

	nilMonth := 0
	summary, err := expenses.Summary(nilMonth)
	if err != nil {
		return fmt.Errorf("failed to count all amount - %w", err)
	}

	if summary > budget {
		fmt.Printf("You went over budget! Budget - %0.2f. Summary - %0.2f\n", budget, summary)
	}
	return nil
}

func switchCommand(input string) error {
	if input == "" {
		return nil
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}

	cmd := parts[0]
	args := parts[1:]

	handler, exists := commands[cmd]
	if !exists {
		fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", cmd)
		return nil
	}

	if err := handler(args); err != nil {
		return fmt.Errorf("command '%s' failed: %w", cmd, err)
	}

	return nil
}

func help() {
	fmt.Println()
	fmt.Println("Список команд -")
	fmt.Println("Добавление траты: | add --description \"Lunch\" --category \"Food\" --amount 20")
	fmt.Println("Удаление траты: | delete --id 2")
	fmt.Println("Вывод всех трат: | list")
	fmt.Println("Установка бюджета: | set-budget --amount 3000")
	fmt.Println("Вывод суммы всех трат: | summary")
	fmt.Println("Вывод суммы трат за месяц: | summary --month 8")
	fmt.Println("Выход: | exit")
	fmt.Println()
}

func exit() {
	os.Exit(0)
}
