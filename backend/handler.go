package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

var currentQuestion = 0
var answersReq = make([]int, 0)

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Answer  int      `json:"answer"`
}

var questions = []Question{
	Question{0, "Как расшифровывается \"JSON\"?", []string{"Jedi's Son On Neptun", "JavaScript Oblect Notation", "Java System On Network", "Just Some Ordinary Nonsense"}, 1},
	Question{1, "Что такое HTML?", []string{"HyperText Markup Language", "Hot Tamale and Meatball Lasagna", "Hyper Technical Machine Language", "Home Tool Markup Language"}, 0},
	Question{2, "Что такое JavaScript?", []string{"Сценарий на Ява-йоге", "Язык, на котором разговаривают кофе-машины", "Язык программирования для веба", "Жареный скрипт"}, 2},
	Question{3, "Что такое 'cookies' в контексте веб-разработки?", []string{"Кружки с печеньем, которые высыпаются из вашего компьютера", "Файлы, хранящиеся на вашем компьютере для отслеживания информации", "Популярный десерт для серверов", "Кодовое слово для хакерской атаки"}, 1},
	Question{4, "Какой символ используется в CSS для обозначения класса?", []string{"exclamation mark (!)", "hashtag (#)", "question mark (?)", "dot (.)"}, 3},
}

func main() {
	http.HandleFunc("/question", questionHandler)
	http.HandleFunc("/answer", answerHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
	})

	handler := c.Handler(http.DefaultServeMux)

	http.ListenAndServe(":8080", handler)
}

func questionHandler(w http.ResponseWriter, r *http.Request) {
	if currentQuestion >= len(questions) {
		calculateResult(w, r)
		return
	}

	json.NewEncoder(w).Encode(questions[currentQuestion])

	currentQuestion++
}

func answerHandler(w http.ResponseWriter, r *http.Request) {

	var question Question
	json.NewDecoder(r.Body).Decode(&question)

	json.NewEncoder(w).Encode(map[string]interface{}{"result": -1})

	answersReq = append(answersReq, question.Answer)
}

func calculateResult(w http.ResponseWriter, r *http.Request) {
	resultValue := 0
	for i, answer := range answersReq {
		if answer == questions[i].Answer {
			resultValue++
		}
	}
	fmt.Println(resultValue)

	json.NewEncoder(w).Encode(struct {
		Result int `json:"result"`
	}{
		Result: resultValue,
	})
}
