package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

const csvName = "problems.csv"

// QuestionHandler manages quiz process
type QuestionHandler struct {
	totalQuestions int
	correctAnswers int
}

// HandleQuestion processes csv lines
func (q *QuestionHandler) HandleQuestion(s []string) error {

	question, answer := s[0], s[1]

	// Print question and get answer
	fmt.Printf("%v : ", question)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userAnswer := scanner.Text()

	// Update scores
	q.totalQuestions++
	if userAnswer == answer {
		q.correctAnswers++
	}

	return nil
}

// PrintResult to output user quiz results
func (q *QuestionHandler) PrintResult() {

	fmt.Printf("Correct answers : %v\n", q.correctAnswers)
	fmt.Printf("Total questions : %v\n", q.totalQuestions)
}

func main() {

	// Open file
	csvFile, err := os.Open(csvName)

	if err != nil {
		log.Fatalln("Error opening csvfile", err)
	}

	// Parse csv
	reader := csv.NewReader(csvFile)

	questionHandler := QuestionHandler{}

	// Iterate through records
	for {

		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln("Error reading csvfile", err)
		}

		questionHandler.HandleQuestion(record)
	}

	questionHandler.PrintResult()
}
