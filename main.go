package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	csvFile := flag.String("csv", "questions.csv", "a csv file in the format of 'question,answer' ") //Setando a flag para executar o programa
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")                     //Setando a flag do time
	flag.Parse()

	file, err := os.Open(*csvFile) //Abrir o Arquivo CSV
	if err != nil {                //Tratando erro
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFile))
		os.Exit(1)
	}

	r := csv.NewReader(file)  //Ler todo o arquivo
	lines, err := r.ReadAll() //Ler arquivo linha por linha
	if err != nil {           //Tratando erro
		exit("Failed to parse the provided CSV File.")
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second) //Setar timestamp

	correct := 0 //Variavel De respostas corretas inicializada com 0 acertos

problemLoop:
	for i, p := range problems { //Range para percorrer a matriz linha por linha 
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q) // Exubindo o Enunciado das questoes
		answerCh := make(chan string) //Criando canal
		go func() {                   //Função executada em uma nova thred
			var answer string	
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer	//Para a informação do input por channel 
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemLoop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}

	fmt.Printf("You score %d out of %d.\n", correct, len(problems))
}

func parseLines(Lines [][]string) []problem { //Função que recebe as linhas [0][1], e retorna aperna uma slice
	ret := make([]problem, len(Lines)) // Declarando ret como tipo slice
	for i, line := range Lines { //Range para percorrer todas as linhas 
		ret[i] = problem{
			q: line[0], 	//q recebe o conteudo da coluna 0
			a: line[1],		//a recebe o conteudo da coluna 1
		}
	}
	return ret 
}

type problem struct { //Struct question , answer
	q string //Q = Question
	a string //A = Answer
}

func exit(msg string) { //Função para tratar erros
	fmt.Println(msg)
	os.Exit(1)
}
