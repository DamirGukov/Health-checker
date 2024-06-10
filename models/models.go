package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
)

type Question struct {
	ID           int    `db:"id" json:"id"`
	QuestionText string `db:"question_text" json:"question_text"`
}

type UserAnswer struct {
	UserID     int  `db:"user_id"`
	QuestionID int  `db:"question_id"`
	Answer     bool `db:"answer"`
}

type UserAnswers struct {
	Answers map[string]bool `json:"answers"`
}

func GetQuestions(db *sqlx.DB) ([]Question, error) {
	var questions []Question
	err := db.Select(&questions, "SELECT id, question_text FROM questions")
	if err != nil {
		logrus.WithError(err).Error("Error getting questions")
		return nil, err
	}
	return questions, nil
}

func SaveUserAnswers(db *sqlx.DB, answers UserAnswers) error {
	tx := db.MustBegin()
	for questionID, answer := range answers.Answers {
		_, err := tx.Exec("INSERT INTO user_answers (question_id, answer) VALUES ($1, $2)", questionID, answer)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func AnalyzeAnswers(db *sqlx.DB, answers UserAnswers) (string, error) {
	type Disease struct {
		ID       int    `db:"id"`
		Name     string `db:"name"`
		Symptoms string `db:"symptoms"`
	}

	var diseases []Disease

	query := `
        SELECT d.id, d.name, string_agg(s.name, ',') AS symptoms
        FROM diseases d
        JOIN disease_symptoms ds ON ds.disease_id = d.id
        JOIN symptoms s ON s.id = ds.symptom_id
        GROUP BY d.id, d.name
    `

	err := db.Select(&diseases, query)
	if err != nil {
		return "", err
	}

	matches := make(map[string]int)
	for _, disease := range diseases {
		symptoms := strings.Split(disease.Symptoms, ",")
		for _, symptom := range symptoms {
			for questionID, answer := range answers.Answers {
				// Сравнение как строки
				if (symptom == "headache" && questionID == "1" && answer) ||
					(symptom == "fever" && questionID == "2" && answer) ||
					(symptom == "nausea" && questionID == "3" && answer) {
					matches[disease.Name]++
				}
			}
		}
	}

	if len(matches) == 0 {
		return "No diagnosis found", nil
	}

	var bestMatch string
	var maxCount int
	for disease, count := range matches {
		if count > maxCount {
			bestMatch = disease
			maxCount = count
		}
	}

	return bestMatch, nil
}
