package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const DataDir = "data"

func main() {
	fmt.Println("Quiz starting")
	f := openQuizFile()
	Play(f)
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

func openQuizFile() io.Reader {
	workingDir := getWorkingDir()
	quizFilePath := filepath.Join(workingDir, DataDir, "problems.csv")
	f, err := os.Open(quizFilePath)
	if err != nil {
		panic(err)
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
