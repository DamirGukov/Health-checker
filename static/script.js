document.addEventListener("DOMContentLoaded", function() {
    const questionsDiv = document.getElementById("questions");
    const resultDiv = document.getElementById("result");

    function showError(message) {
        resultDiv.innerText = `Error: ${message}`;
        resultDiv.style.display = "block";
        resultDiv.style.color = "red";
    }

    function loadQuestions() {
        fetch("http://localhost:8080/questions")
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                console.log("Questions data:", data); // Логування отриманих даних
                questionsDiv.innerHTML = "";
                data.forEach(question => {
                    console.log("Processing question:", question); // Логування кожного питання
                    const div = document.createElement("div");
                    div.classList.add("question");
                    div.innerHTML = `
                        <label>
                            <input type="checkbox" name="question${question.ID}" value="${question.ID}">
                            ${question.QuestionText}
                        </label>
                    `;
                    questionsDiv.appendChild(div);
                });
            })
            .catch(error => {
                console.error("Error fetching questions:", error);
                showError("Could not load questions. Please try again later.");
            });
    }

    function submitForm(event) {
        event.preventDefault();

        const answers = {};
        const formData = new FormData(form);
        formData.forEach((value, key) => {
            console.log(`Processing: key=${key}, value=${value}`);
            const questionId = key.replace('question', '');
            answers[questionId] = form[key].checked;
        });

        console.log("Submitting answers:", answers);

        fetch("http://localhost:8080/submit", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ answers })
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                console.log("Diagnosis data:", data);
                resultDiv.innerText = `Diagnosis: ${data.diagnosis}`;
                resultDiv.style.display = "block";
                resultDiv.style.color = "black";
            })
            .catch(error => {
                console.error("Error submitting answers:", error);
                showError("Could not submit answers. Please try again later.");
            });
    }

    const form = document.getElementById("healthForm");
    form.addEventListener("submit", submitForm);

    // Load questions on page load
    loadQuestions();
});
