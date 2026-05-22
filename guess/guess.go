package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Logs struct {
	Date         string `json:"date"`
	Win          bool   `json:"win"`
	Attempt      int    `json:"attempt"`
	Level        string `json:"level"`
	SecretNumber int    `json:"secret_number"`
}

var (
	reset   = "\033[0m"
	success = "\033[32;1m"
	failure = "\033[31;1m"
	warning = "\033[33;1m"
	info    = "\033[36;1m"
)

func guessNumber(randomNumber, interval, maxAttempt int, level string) []int {
	var (
		attempt int = 0
	)
	remember := make([]int, maxAttempt)
	scanner := bufio.NewScanner(os.Stdin)

	for attempt < maxAttempt {
		fmt.Printf(info+"Попытка #%d - Введите число: "+reset, attempt+1)
		guess, err := correctInput(scanner, interval)

		if err != nil {
			fmt.Printf(failure+"❌ Ошибка: %v\n"+reset, err)
			continue
		}

		remember[attempt] = guess

		if guess == randomNumber {
			fmt.Println(success + "Вы угадали!🙌" + reset)
			fmt.Println(success + "Игра закончена!" + reset)
			err := saveJson(true, attempt+1, level, randomNumber)
			if err != nil {
				return nil
			}
			return remember[:attempt+1]
		}

		if guess > randomNumber {
			fmt.Println(info + "Секретное число меньше👇" + reset)
		} else {
			fmt.Println(info + "Секретное число больше👆" + reset)
		}

		giveHint(randomNumber, guess)
		attempt++
	}

	failureAnswer(randomNumber)

	err := saveJson(false, attempt, level, randomNumber)
	if err != nil {
		return nil
	}

	return remember
}

func saveJson(win bool, attempt int, level string, secretNumber int) error {
	var results []Logs

	file, err := os.ReadFile("results.json")
	if err == nil {
		json.Unmarshal(file, &results)
	}

	results = append(results, Logs{
		Win:          win,
		Attempt:      attempt,
		Level:        level,
		Date:         time.Now().Format("2006-01-02 15:04:05"),
		SecretNumber: secretNumber,
	})

	file, err = json.MarshalIndent(results, "", "  ")
	err = os.WriteFile("results.json", file, 0644)

	if err != nil {
		return err
	}

	return nil
}

func failureAnswer(randomNumber int) {
	fmt.Println(failure + "Вы проиграли!😢" + reset)
	fmt.Printf(failure+"Секретное число было: %d\n"+reset, randomNumber)
	fmt.Println(failure + "Игра закончена!" + reset)
}

func correctInput(scanner *bufio.Scanner, interval int) (int, error) {
	fmt.Printf(info+"Введите число (от 1 до %d): "+reset, interval)
	scanner.Scan()

	input := strings.TrimSpace(scanner.Text())
	guess, err := strconv.Atoi(input)

	if err != nil {
		return 0, fmt.Errorf(warning + "пожалуйста, введите только целое число" + reset)
	}

	if guess < 1 || guess > interval {
		return 0, fmt.Errorf(warning+"число должно быть в диапазоне от 1 до %d"+reset, interval)
	}

	return guess, nil
}

func giveHint(randomNumber, guess int) {
	difference := int(math.Abs(float64(randomNumber - guess)))

	switch {
	case difference <= 5:
		fmt.Println(warning + "🔥 Горячо" + reset)
	case difference <= 15:
		fmt.Println(info + "🙂 Тепло" + reset)
	default:
		fmt.Println(info + "❄️ Холодно" + reset)
	}
}

func chooseLevel(level string) (interval, attempts int) {
	level = strings.ToLower(level)

	switch level {
	case "easy":
		return 50, 15
	case "medium":
		return 100, 10
	case "hard":
		return 200, 5
	default:
		fmt.Println(failure + "❌ Неверный выбор уровня сложности!" + reset)
		fmt.Println(info + "Доступные уровни: Easy, Medium, Hard" + reset)
		return 0, 0
	}
}

func main() {
	fmt.Println(info + "Игра 'Угадай число' - от 1 до 100 началась!" + reset)
	fmt.Println(info + "Угадайте число за 10 попыток!" + reset)

	var level string
	fmt.Printf(info + "Выберите уровень сложности: " + reset)
	fmt.Scan(&level)
	interval, attempt := chooseLevel(level)

	secretNumber := rand.Intn(interval) + 1
	remember := guessNumber(secretNumber, interval, attempt, level)

	fmt.Println(info+"Ваши попытки:"+reset, remember)
}
