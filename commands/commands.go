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

func parseCommand(commands string) (string, []string) {
	if commands == "" {
		return "", nil
	}
	input := strings.Split(commands, " ")
	cmd := input[0]
	args := input[1:]

	return cmd, args
}

func expectArgs(count int, args []string) error {
	if count != len(args) {
		return fmt.Errorf("failed args. Expected - %d, have - %d", count, len(args))
	}
	return nil
}

func addArgs(args []string) (string, float64, error) {
	err := expectArgs(4, args)
	if err != nil {
		return "", 0, fmt.Errorf("incorrect args. Expected - %d, have - %d", 4, len(args))
	}

	var description string
	var amount float64

	if args[0] == "--description" && args[2] == "--amount" {
		description = args[1]
		amount, err = strconv.ParseFloat(args[3], 64)
		if err != nil {
			return "", 0, fmt.Errorf("incorrect amount - %w", err)
		}
	}

	return description, amount, nil
}

func delArgs(args []string) (int, error) {
	err := expectArgs(2, args)
	if err != nil {
		return 0, fmt.Errorf("incorrect input. Expected - %d, have - %d args", 2, len(args))
	}
	if args[0] != "--id" {
		return 0, fmt.Errorf("incorrect flag")
	}
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return 0, fmt.Errorf("failed id - %w", err)
	}
	return id, nil
}

func switchCommand(commands string) error {
	cmd, args := parseCommand(commands)

	switch cmd {
	case "help":
		help()
	case "exit":
		exit()
	case "add":
		description, amount, err := addArgs(args)
		if err != nil {
			return fmt.Errorf("failed to addArgs - %w", err)
		}
		err = expenses.AddExpense(description, amount)
		if err != nil {
			return fmt.Errorf("failed to add -%w", err)
		}
	case "delete":
		id, err := delArgs(args)
		if err != nil {
			return fmt.Errorf("incorrect id - %w", err)
		}
		err = expenses.DeleteExpense(id)
		if err != nil {
			return fmt.Errorf("failed to delete expense - %w", err)
		}
	case "list":
		{
		}
		expenses.ListExpense()
	}
	return nil
}

func help() {
	fmt.Println()
	fmt.Println("Список команд -")
	fmt.Println("Добавление траты: | add --description \"Lunch\" --amount 20 ")
	fmt.Println("Удаление траты: | delete --id 2")
	fmt.Println("Вывод всех трат: | list")
	fmt.Println("Установка бюджета: | set-budget --amount 3000")
	fmt.Println("Вывод суммы всех трат: | summary")
	fmt.Println("Вывод суммы трат за месяц: | summary --month 8")
	fmt.Println("Вывод суммы трат по категориям: | summary --category \"some category\"")
	fmt.Println("Выход: | exit")
	fmt.Println()
}

func exit() {
	os.Exit(0)
}
