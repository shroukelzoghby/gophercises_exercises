package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const csvFile = "problems.csv"

func main() {
	timeLimit, shouldShuffle := parseFlags()
	records, err := loadProblems(csvFile)
	if err != nil {
		fmt.Println("Error loading quiz:", err)
		return
	}

	if shouldShuffle {
		shuffleProblems(records)
	}

	reader := bufio.NewReader(os.Stdin)
	if err := waitForStart(reader); err != nil {
		fmt.Println("Error reading start input:", err)
		return
	}

	deadline := time.After(time.Duration(timeLimit) * time.Second)
	correctAnswers := 0

	for i, record := range records {
		if len(record) < 2 {
			continue
		}

		correct, timedOut := askQuestion(reader, i+1, record[0], record[1], deadline)
		if timedOut {
			fmt.Println("Time is up!")
			break
		}
		if correct {
			correctAnswers++
		}
	}

	printScore(correctAnswers, len(records))
}

func parseFlags() (int, bool) {
	timeLimit := flag.Int("limit", 30, "time limit for the quiz in seconds")
	shuffle := flag.Bool("shuffle", false, "shuffle the quiz questions")
	flag.Parse()
	return *timeLimit, *shuffle
}

func loadProblems(filename string) ([][]string, error) {
	openFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer openFile.Close()

	return csv.NewReader(openFile).ReadAll()
}

func shuffleProblems(records [][]string) {
	rand.Shuffle(len(records), func(i, j int) {
		records[i], records[j] = records[j], records[i]
	})
}

func waitForStart(reader *bufio.Reader) error {
	fmt.Println("Press Enter to start the quiz.")
	_, err := reader.ReadString('\n')
	return err
}

func askQuestion(reader *bufio.Reader, number int, question, correctAnswer string, deadline <-chan time.Time) (bool, bool) {
	fmt.Printf("Question %d: %s\n", number, question)
	answerCh := make(chan string)
	go func() {
		answer, err := reader.ReadString('\n')
		if err == nil {
			answerCh <- answer
		}
	}()

	select {
	case answer := <-answerCh:
		return strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(correctAnswer)), false
	case <-deadline:
		return false, true
	}
}

func printScore(correctAnswers, totalQuestions int) {
	fmt.Printf("Your score is: %d/%d\n", correctAnswers, totalQuestions)
}
