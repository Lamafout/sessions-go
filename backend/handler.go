package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret"))

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
	http.HandleFunc("/nextQuestion", answerHandler)
	http.ListenAndServe(":8080", nil)
}

func questionHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "quiz-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentQuestionID, ok := session.Values["currentQuestion"].(int)
	if !ok {
		currentQuestionID = 0
	}

	if currentQuestionID >= len(questions) {
		calculateResult(w, r)
		return
	}

	json.NewEncoder(w).Encode(questions[currentQuestionID])

	session.Values["currentQuestion"] = currentQuestionID + 1
	session.Save(r, w)
}

func answerHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "quiz-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var question Question
	err = json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["answers"] = append(session.Values["answers"].([]int), question.Answer)
	session.Save(r, w)
}

func calculateResult(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "quiz-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resultValue := 0
	for i, answer := range session.Values["answers"].([]int) {
		if answer == questions[i].Answer {
			resultValue++
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"result": resultValue})
}
