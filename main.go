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
	english string
	french  string
}

type UserInput struct {
	Filename string
}

func synthesizeSpeech(svc *polly.Polly, word Word, i int, c chan string) {
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

	outFile, err := os.Create(fmt.Sprintf("mp3s/%d-%s.mp3", i+1, word.english))
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

	c <- fmt.Sprintf("Success: %s", word.english)
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
		word := Word{line[0], line[1]}
		go synthesizeSpeech(svc, word, i, c)
	}

	for i := 0; i < len(csvLines); i++ {
		fmt.Printf("Completed: %s\n", <-c)
	}

	fmt.Println("DONE")
}
