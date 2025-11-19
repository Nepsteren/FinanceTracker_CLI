package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func greeting() {
	fmt.Println("Вас приветствует - FinanceTrackerCLI")
	fmt.Println("Если вам нужна помощь введите команду: help")
}

func Start() {
	greeting()
	userInput()
}

func userInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			switchCommand(scanner.Text())
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

func switchCommand(commands string) {
	cmd, _ := parseCommand(commands)

	switch cmd {
	case "help":
		help()
	case "exit":
		exit()
	case "add":

	}
}

func help() {
	fmt.Println()
	fmt.Println("Список команд -")
	fmt.Println("Добавление траты: | add --description \"Lunch\" --amount 20 ")
	fmt.Println("Удаление траты: | delete --id 2")
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
