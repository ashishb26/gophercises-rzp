package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

/*
	getUserAnswer reads user answers, pushes it into the channel
	and startQuiz reads the answer from the channel
*/
func getUserAnswer(ansChan chan string) {

	var answer string
	reader := bufio.NewReader(os.Stdin)
	answer, _ = reader.ReadString('\n')
	answer = strings.Replace(answer, "\r\n", "", -1)
	ansChan <- answer

}

/*
	startQuiz starts the timer and executes the quiz, displaying one question at a time
	It evaluates the answers and returns the number of correct answers
*/
func startQuiz(records [][]string, timerFlag *time.Duration) int {

	numCorrectAns := 0
	testTimer := time.NewTimer(*timerFlag)

	for index, record := range records {
		fmt.Printf("Question %d: %s \n", index+1, record[0])

		ansChan := make(chan string)
		go getUserAnswer(ansChan)
		f := false

		select {
		case <-testTimer.C:
			fmt.Println("Time exceeded")
			testTimer.Stop()
			f = true
		case ans := <-ansChan:
			if strings.Compare(strings.Trim(strings.ToLower(ans), "\n"), record[1]) == 0 {
				numCorrectAns++
			}
		}
		if f {
			break
		}
	}
	return numCorrectAns
}

/*
	shuffleRec shuffles the order of questions displayed
*/
func shuffleRec(records [][]string) {

	rand.Seed(time.Now().UnixNano())
	for i := len(records) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		records[i], records[j] = records[j], records[i]
	}

}

func main() {

	var timerFlag = flag.Duration("timer", 30*time.Second, "Flag to set test duration. Input format : `<time>s`(without quotes)")
	var testFile = flag.String("test", "problems.csv", "File name of the test set.")
	var shuffle = flag.Bool("shuffle", false, "Boolean flag to shuffle the test")
	flag.Parse()

	csvFile, err := os.Open(*testFile)
	if err != nil {
		log.Fatalln("Error opening test file:", err)
		os.Exit(1)
	}

	csvReader := csv.NewReader(csvFile)
	stdReader := bufio.NewReader(os.Stdin)
	records, err := csvReader.ReadAll()

	if *shuffle {
		shuffleRec(records)
	}

	if err != nil {
		log.Fatalln("Cannot parse test file:", err)
		os.Exit(1)
	}

	fmt.Println("Please enter any key to start the test")
	_, err = stdReader.ReadBytes('\n')

	if err != nil {
		log.Fatalln("Error reading input:", err)
		os.Exit(1)
	}

	numCorrectAns := startQuiz(records, timerFlag)

	fmt.Printf("Score: %d/%d", numCorrectAns, len(records))
}
