package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// QuestionHandler manages quiz process
type QuestionHandler struct {
	totalQuestions int
	correctAnswers int
}

// HandleQuestion processes csv lines
func (q *QuestionHandler) HandleQuestion(p Problem, ch <-chan time.Time) error {

	// Print question and get answer
	q.totalQuestions++
	fmt.Printf("Problem #%d , %s : ", q.totalQuestions, p.question)

	answerCh := make(chan string)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		answerCh <- scanner.Text()
	}()

	select {
	case <-ch:
		q.printResult()
	case answer := <-answerCh:
		// Update scores
		if answer == p.answer {
			q.correctAnswers++
		}
	}

	return nil
}

func (q *QuestionHandler) printResult() {

	fmt.Printf("\nCorrect answers : %v\n", q.correctAnswers)
	fmt.Printf("Total questions : %v\n", q.totalQuestions)
	os.Exit(0)

}

// Problem is a data structure that holds a question and answer
type Problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {

	// Flags
	csvName := flag.String("csv", "problems.csv", "a csv file in the format of `question, answer` ")
	duration := flag.Int("time", 5, "duration of quiz in seconds")

	flag.Parse()

	// Open file
	csvFile, err := os.Open(*csvName)
	if err != nil {
		exit(fmt.Sprintf("Error opening csvfile : %s\n", *csvName))
	}

	// Parse csv
	reader := csv.NewReader(csvFile)
	questionHandler := QuestionHandler{}

	// Set timer
	timer := time.NewTimer(time.Duration(*duration) * time.Second)

	// Iterate through records
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			exit("Error processing csv file")
		}

		problem := Problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}

		questionHandler.HandleQuestion(problem, timer.C)
	}
}
