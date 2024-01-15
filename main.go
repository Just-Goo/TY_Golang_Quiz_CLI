package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type quiz struct {
	Question, Answer string
}

func main() {
	fileName := flag.String("f", "quiz.csv", "csv file path")
	timer := flag.Int("t", 30, "quiz timer")
	flag.Parse()

	quizQuestions, err := getQuestions(*fileName)
	if err != nil {
		log.Fatalf("an error occurred: %v", err.Error())
	}

	var correctAnswers int
	quizTimer := time.NewTimer(time.Duration(*timer) * time.Second)

	answersChannel := make(chan string) // create a channel for the answers

quizLoop:
	for index, quiz := range quizQuestions {
		var answer string 
		fmt.Printf("Question %d: %s = ", index+1, quiz.Question)

		go func() {  
			fmt.Scan(&answer)
			answersChannel <- answer
		}()

		select {
		case <-quizTimer.C: // If timer is done 
			close(answersChannel)
			break quizLoop
		case userAnswer := <-answersChannel:
			if userAnswer == quiz.Answer { // check if user's answer is correct
				correctAnswers++
			}
			if index == (len(quizQuestions) - 1) { // If all questions have been asked
				close(answersChannel) // close the channel
			}
		}
	}

	fmt.Println()
	fmt.Printf("You scored %d out of %d questions", correctAnswers, len(quizQuestions))
	<-answersChannel

}

func getQuestions(filename string) ([]quiz, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf(" %v, while reading data from csv file: %v", err.Error(), filename)
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	cLines, err := csvReader.ReadAll()

	if err != nil {
		return nil, fmt.Errorf("%v, when opening file: %v", err.Error(), filename)
	}

	return parseQuiz(cLines), nil
}

func parseQuiz(lines [][]string) []quiz {
	quizSlice := make([]quiz, len(lines)) 
	for index, value := range lines { 
		quizSlice[index] = quiz{Question: value[0], Answer: value[1]}
	} 
	return quizSlice
}
