package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

type Word struct {
	id      string
	english string
	french  string
}

type UserInput struct {
	Filename string
}

func synthesizeSpeech(svc *polly.Polly, word Word, i int, c chan string) {
	exists := "exists"
	fileName := fmt.Sprintf("mp3s/%s.mp3", word.id)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		input := &polly.SynthesizeSpeechInput{
			OutputFormat: aws.String("mp3"),
			Text:         aws.String(word.french),
			VoiceId:      aws.String("Celine"),
		}
		output, err := svc.SynthesizeSpeech(input)
		if err != nil {
			c <- fmt.Sprintf("Failed %s. Err: %s", word.english, err)
			panic(err)
		}

		outFile, err := os.Create(fileName)
		if err != nil {
			c <- fmt.Sprintf("Failed %s. Err: %s", word.english, err)
			panic(err)
		}

		defer outFile.Close()

		_, err = io.Copy(outFile, output.AudioStream)
		if err != nil {
			c <- fmt.Sprintf("Failed %s. Err: %s", word.english, err)
			panic(err)
		}
		exists = "created"
	}
	c <- fmt.Sprintf("Success\nAction: %s\nWord: %s", exists, word.english)
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		panic("Filename is required")
	}

	ui := UserInput{
		args[0],
	}

	csvFile, err := os.Open(ui.Filename)

	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		panic(err)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := polly.New(sess)

	c := make(chan string)

	for i, line := range csvLines {
		word := Word{line[0], line[1], line[2]}
		go synthesizeSpeech(svc, word, i, c)
	}

	for i := 0; i < len(csvLines); i++ {
		fmt.Printf("Completed: %s\n", <-c)
	}

	fmt.Println("DONE")
}
