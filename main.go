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

const DefaultProblemsFile string = "problems.csv"

var problemsFile string
var timelimit int64

func GetProblemsPath() string {
	if problemsFile == "" {
		problemsFile = DefaultProblemsFile
	}
	return problemsFile
}

type Question struct {
	Q string
	A string
}

func (q *Question) Load(r []string) {
	q.Q = r[0]
	q.A = r[1]
}

type Quiz struct {
	Questions []Question
	Correct   int
	Incorrect int
	I         int
	eoq       EndOfQuiz
}

type EndOfQuiz struct {
	Q int
}

func (e *EndOfQuiz) Error() string {
	return fmt.Sprintf("Quiz Complete. Only %v question(s)", e.Q)
}

func (q *Quiz) NextQuestion() (Question, error) {
	if q.I < len(q.Questions) {
		return q.Questions[q.I], nil
	}
	return Question{}, &q.eoq
}

func (q *Quiz) CheckAnswer(a string) bool {
	question, _ := q.NextQuestion()
	r := (a == question.A)
	q.I = q.I + 1
	if r {
		q.Correct += 1
	} else {
		q.Incorrect += 1
	}
	return r
}

func (q *Quiz) TimesUp() {
	fmt.Println("You ran out of time")
	q.Incorrect = len(q.Questions) - q.Correct
	q.I = len(q.Questions)
}

func (q *Quiz) Report() {
	p := float64(q.Correct) / float64(len(q.Questions)) * 100
	fmt.Printf("\nYou got %v/%v (%.1f%%) Correct (%v Incorrect)\n\n", q.Correct, len(q.Questions), p, q.Incorrect)
}

func NewQuiz(p string) (*Quiz, error) {
	var q Quiz
	csvfile, err := os.Open(GetProblemsPath())
	if err != nil {
		return &q, err
	}
	r := csv.NewReader(csvfile)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &q, err
		}
		var question Question
		question.Load(record)
		q.Questions = append(q.Questions, question)
	}
	q.eoq.Q = len(q.Questions)

	return &q, nil
}

func init() {
	flag.StringVar(&problemsFile, "src", DefaultProblemsFile, "Path to problems CSV source")
	flag.Int64Var(&timelimit, "timelimit", 30, "Specify the time limit of the quiz.")
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	flag.Parse()
	q, err := NewQuiz(problemsFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\nPress Enter when ready to begin quiz\n")
	reader.ReadString('\n')

	timer := time.NewTimer(time.Duration(timelimit) * time.Second)
	go func() {
		<-timer.C
		q.TimesUp()
	}()

	for _, question := range q.Questions {
		if q.I >= len(q.Questions) {
			break
		}
		fmt.Printf("\n%s\n", question.Q)
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		q.CheckAnswer(text)
	}
	timer.Stop()
	q.Report()
}
