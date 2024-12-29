package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Name      string
	Score     int
	Questions []Question
}

func (gameState *GameState) Init() {
	fmt.Println("Wellcome to the Quiz")
	fmt.Println("Write your name: ")
	reader := bufio.NewReader(os.Stdin)

	name, err := reader.ReadString('\n')

	if err != nil {
		panic("Error when reading name")
	}

	gameState.Name = name

	fmt.Printf("Lets to the game %s", gameState.Name)
}

func (gameState *GameState) ProcessCSV() {
	file, err := os.Open("quizgo.csv")
	if err != nil {
		panic("Error when reading file")
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		panic("Error when reading csv")
	}

	for index, record := range records {
		if index > 0 {
			correctAnswer, _ := toInt(record[5])
			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer:  correctAnswer,
			}

			gameState.Questions = append(gameState.Questions, question)
		}
	}
}

func (gamaState *GameState) Run() {
	const yellow = "\033[33m"
	const reset = "\033[0m"
	for indexQuestion, question := range gamaState.Questions {
		fmt.Printf("%s%d. %s%s\n", yellow, indexQuestion+1, question.Text, reset)

		for indexOption, option := range question.Options {
			fmt.Printf("[%d] %s\n", indexOption+1, option)
		}

		fmt.Println("Write your answer")

		var answer int
		var err error

		for {
			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\n')

			answer, err = toInt(read[:len(read)-1]) // remove /n
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			break
		}

		if answer == question.Answer {
			fmt.Println("Congratulations, you're right!")
			gamaState.Score += 10
		} else {
			fmt.Println("This is bad, you are wrong!")
		}
	}
}

func main() {
	game := &GameState{}
	go game.ProcessCSV()
	game.Init()
	game.Run()

	fmt.Printf("End of the game, you scored %d points\n", game.Score)
}

func toInt(s string) (int, error) {
	i, err := strconv.Atoi(s)

	if err != nil {
		return 0, errors.New("characters other than numbers are not allowed")
	}

	return i, nil
}
