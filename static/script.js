document.addEventListener("DOMContentLoaded", function() {
    let currentQuestionIndex = 0;
    let answers = {};
    let questions = [];

    const questionElement = document.getElementById("question");
    const questionCounterElement = document.getElementById("questionCounter");
    const resultElement = document.getElementById("result");

    function showQuestion(index) {
        if (questions.length > 0) {
            questionElement.innerText = questions[index].QuestionText;
            questionCounterElement.innerText = `Питання ${index + 1} з ${questions.length}`;
        } else {
            console.error("No questions available to show.");
        }
    }

    function showResult(diagnosis) {
        questionElement.style.display = "none";
        document.querySelector(".button-container").style.display = "none";
        questionCounterElement.style.display = "none";
        resultElement.innerText = `Diagnosis: ${diagnosis}`;
        resultElement.classList.add("show");
    }

    document.getElementById("yesButton").addEventListener("click", function() {
        answers[questions[currentQuestionIndex].ID] = true;
        nextQuestion();
    });

    document.getElementById("noButton").addEventListener("click", function() {
        answers[questions[currentQuestionIndex].ID] = false;
        nextQuestion();
    });

    function nextQuestion() {
        currentQuestionIndex++;
        if (currentQuestionIndex < questions.length) {
            showQuestion(currentQuestionIndex);
        } else {
            submitAnswers();
        }
    }

    function submitAnswers() {
        fetch("http://localhost:8080/submit", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ answers: answers })
        })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    console.error("Diagnosis error:", data.error);
                } else {
                    console.log("Diagnosis data:", data);
                    showResult(data.diagnosis);
                }
            })
            .catch(error => {
                console.error("Error submitting answers:", error);
            });
    }

    fetch("http://localhost:8080/questions")
        .then(response => response.json())
        .then(data => {
            questions = data;
            if (questions.length > 0) {
                showQuestion(currentQuestionIndex);
            } else {
                console.error("No questions received from the server.");
            }
        })
        .catch(error => {
            console.error("Error fetching questions:", error);
        });
});
