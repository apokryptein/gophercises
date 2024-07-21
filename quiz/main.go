package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Flags struct {
	fileName  string
	timeLimit int
	shuffle   bool
}

type Question struct {
	question string
	answer   string
}

type Quiz struct {
	questions   []Question
	questionNum int
	correctNum  int
}

func main() {
	flags := getFlags()

	csvFile := flags.fileName
	timeLimit := flags.timeLimit
	shuffle := flags.shuffle

	printBanner(csvFile, timeLimit)
	ConfirmStart()

	quiz := newQuiz(csvFile)
	if shuffle {
		quiz.shuffle()
	}

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	go func() {
		for i, q := range quiz.questions {
			var answer string

			fmt.Printf("Problem #%d: %s = ", i+1, q.question)
			fmt.Scanln(&answer)

			if answer == q.answer {
				quiz.correctNum++
			}
		}
	}()

	<-timer.C
	fmt.Println("\n======== TIME IS UP ========")
	quiz.reportScore()
}

func getFlags() Flags {
	csvFile := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	shuffle := flag.Bool("s", false, "shuffle quiz questions")

	flag.Parse()

	if len(os.Args) < 1 {
		flag.Usage()
		os.Exit(0)
	}

	flags := Flags{
		fileName:  *csvFile,
		timeLimit: *timeLimit,
		shuffle:   *shuffle,
	}

	return flags
}

func loadFile(fileName string) []Question {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	fileScan := bufio.NewScanner(f)

	questions := []Question{}

	for fileScan.Scan() {
		line := strings.Split(fileScan.Text(), ",")
		questions = append(questions, Question{question: line[0], answer: strings.TrimSpace(line[1])})
	}
	return questions
}

func newQuiz(f string) Quiz {
	q := Quiz{}

	q.questions = loadFile(f)
	q.questionNum = len(q.questions)
	q.correctNum = 0

	return q
}

func (q Quiz) reportScore() {
	fmt.Printf("You scored %d out of %d\n", q.correctNum, q.questionNum)
}

func (q Quiz) shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range q.questions {
		newPosition := r.Intn(len(q.questions) - 1)

		q.questions[i], q.questions[newPosition] = q.questions[newPosition], q.questions[i]
	}
}

func printBanner(csvFile string, timeLimit int) {
	logo := `
  
   ____       _       __ _                   
  /___ \_   _(_)____ / _\ |__   _____      __
 //  / / | | | |_  / \ \| '_ \ / _ \ \ /\ / /
/ \_/ /| |_| | |/ /  _\ \ | | | (_) \ V  V / 
\___,_\ \__,_|_/___| \__/_| |_|\___/ \_/\_/  
                                             
  `
	fmt.Println(logo)

	fmt.Println("-------------------------------------")
	fmt.Printf("Quiz File: %s\nTime Limit: %d\n", csvFile, timeLimit)
	fmt.Println("-------------------------------------")
}

func ConfirmStart() {
	fmt.Printf("\nPress ENTER to Begin!\n")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
