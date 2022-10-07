package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	DataDir         = "data"
	DefaultQuizFile = "problems.csv"
	LimitSeconds    = 5
)

func main() {
	gameSeconds, quizFile := parseInArgs()
	f := openQuizFile(quizFile)
	PlayWithTimer(f, gameSeconds)
}

func PlayWithTimer(f io.Reader, seconds int) {
	done := make(chan int)
	gameDuration := getGameDuration(seconds)
	gameTally := NewGameTally()

	startGame()
	go PlayAsync(f, done, gameTally)
	select {
	case <-time.After(gameDuration):
		fmt.Println("\nTime is up, you did not finish all the quiz!")
	case <-done:
		fmt.Println("\nYou have finished all the quiz!")
		close(done)
	}
	gameTally.printScore()
}

func startGame() {
	fmt.Println("Press enter to start the game")
	r := bufio.NewReader(os.Stdin)
	parseUserAnswer(r)
	fmt.Println("Game starting ...")
}

func getGameDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

func PlayAsync(f io.Reader, done chan int, tally *GameTally) {
	stdInReader := bufio.NewReader(os.Stdin)
	csvReader := csv.NewReader(f)
	for {
		r, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		tally.increaseTotal()
		quiz, expectedAnswer := r[0], r[1]
		fmt.Printf("quiz: %s = ", quiz)

		userAnswer := parseUserAnswer(stdInReader)
		if answerIsCorrect(expectedAnswer, userAnswer) {
			tally.increaseCorrect()
		}
	}
	done <- 0
}

func Play(f io.Reader) {
	var totalQuizCnt, correctAnswerCnt int32
	stdInReader := bufio.NewReader(os.Stdin)
	csvReader := csv.NewReader(f)
	for {
		r, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		totalQuizCnt += 1
		quiz, expectedAnswer := r[0], r[1]
		fmt.Printf("quiz: %s = ", quiz)

		userAnswer := parseUserAnswer(stdInReader)
		if answerIsCorrect(expectedAnswer, userAnswer) {
			correctAnswerCnt += 1
		}
	}
	fmt.Printf("You have got %d out of %d quiz right", correctAnswerCnt, totalQuizCnt)
}

func answerIsCorrect(expected string, got string) bool {
	fmt.Println("your answer is", got, "expected", expected)
	return expected == got
}

func parseUserAnswer(reader *bufio.Reader) string {
	userAnswer, _ := reader.ReadString('\n')
	userAnswer = strings.ReplaceAll(userAnswer, "\n", "")
	return userAnswer
}

func openQuizFile(fileName string) io.Reader {
	workingDir := getWorkingDir()
	quizFilePath := filepath.Join(workingDir, DataDir, fileName)
	f, err := os.Open(quizFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	return f
}

func getWorkingDir() string {
	_, fileName, _, ok := runtime.Caller(0)
	if !ok {
		panic(fileName)
	}
	return filepath.Dir(fileName)
}

func parseInArgs() (int, string) {
	var t int
	var quizFileName string
	flag.IntVar(&t, "seconds", LimitSeconds, "time limit")
	flag.StringVar(&quizFileName, "quiz-file", DefaultQuizFile, "quiz file name")
	flag.Parse()
	return t, quizFileName
}
