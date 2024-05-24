let question = document.querySelector('.question')
let answers = document.querySelectorAll('.answer')
let questionData = {
    id: 0,
    text: '',
    options: Array(answers).map(answer => answer.innerHTML),
    answer: 0
}

async function buildQuestion() {
    const response = await fetch('http://localhost:8080/question')
    const data = await response.json()
    questionData = data;

    if (data.result == null){
        question.innerHTML = questionData.text
        answers.forEach((answer, index) => {
            answer.innerHTML = questionData.options[index]
        })
    }
    else{
        alert('Вы ответили на все вопросы! Результат: ' + data.result + '/5')   
    }
}

async function sendAnswer(index) {
    await fetch('http://localhost:8080/answer', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            id: questionData.id,
            text: questionData.text,
            options: questionData.options,
            answer: index
        })
    })
}

buildQuestion();

answers.forEach((answer, index) => {
    answer.addEventListener('click', async () => {
        await sendAnswer(index)
        await buildQuestion()
    })
})