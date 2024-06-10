document.addEventListener("DOMContentLoaded", function() {
    fetch("http://localhost:8080/questions")
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log("Questions data:", data);
            const questionsDiv = document.getElementById("questions");
            questionsDiv.innerHTML = "";
            data.forEach(question => {
                const div = document.createElement("div");
                div.classList.add("question");
                div.innerHTML = `
                    <label>
                        <input type="checkbox" name="question${question.id}" value="${question.id}">
                        ${question.question_text}
                    </label>
                `;
                questionsDiv.appendChild(div);
            });
        })
        .catch(error => {
            console.error("Error fetching questions:", error);
        });

    const form = document.getElementById("healthForm");
    form.addEventListener("submit", function(event) {
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
                const resultDiv = document.getElementById("result");
                resultDiv.innerText = `Diagnosis: ${data.diagnosis}`;
                resultDiv.style.display = "block";
            })
            .catch(error => {
                console.error("Error submitting answers:", error);
            });
    });
});

