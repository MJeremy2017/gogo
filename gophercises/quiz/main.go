package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
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

var (
	gameSeconds int
	quizFile    string
	shuffle     bool
)

func main() {
	parseInArgs()
	PlayWithTimer()
}

func PlayWithTimer() {
	done := make(chan int)
	gameDuration := getGameDuration(gameSeconds)
	gameTally := NewGameTally()

	problems := parseAllProblems()
	gameTally.SetTotal(int32(len(problems.quiz)))

	go PlayAsync(done, gameTally, problems)
	select {
	case <-time.After(gameDuration):
		fmt.Println("\nTime is up, you did not finish all the quiz!")
	case <-done:
		fmt.Println("\nYou have finished all the quiz!")
		close(done)
	}
	gameTally.printScore()
}

func parseAllProblems() Problems {
	var qs []string
	var ans []string
	f := openQuizFile(quizFile)
	csvReader := csv.NewReader(f)
	for {
		r, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		quiz, expectedAnswer := r[0], r[1]
		qs = append(qs, quiz)
		ans = append(ans, expectedAnswer)
	}
	p := Problems{qs, ans}
	if shuffle {
		shuffleProblems(&p)
	}
	return p
}

func shuffleProblems(p *Problems) {
	var ind []int
	q, ans := p.quiz, p.answers
	n := len(q)
	for i := 0; i < n; i++ {
		ind = append(ind, i)
	}
	rand.Seed(time.Now().Unix())
	rand.Shuffle(n, func(i int, j int) {
		q[i], q[j] = q[j], q[i]
		ans[i], ans[j] = ans[j], ans[i]
	})
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

func PlayAsync(done chan int, tally *GameTally, problems Problems) {
	startGame()
	stdInReader := bufio.NewReader(os.Stdin)
	for i := 0; i < len(problems.quiz); i++ {
		quiz, expectedAnswer := problems.quiz[i], problems.answers[i]
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
	expected = formatString(expected)
	got = formatString(got)
	fmt.Println("your answer is", got, "expected", expected)
	return expected == got
}

func formatString(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
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

func parseInArgs() {
	flag.IntVar(&gameSeconds, "seconds", LimitSeconds, "time limit")
	flag.StringVar(&quizFile, "quiz-file", DefaultQuizFile, "quiz file name")
	flag.BoolVar(&shuffle, "shuffle", false, "shuffle problems")
	flag.Parse()
}
