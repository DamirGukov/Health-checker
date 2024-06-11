package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Question struct {
	ID           int    `db:"id"`
	QuestionText string `db:"question_text"`
}

type Symptom struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Disease struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type DiseaseSymptom struct {
	DiseaseID int `db:"disease_id"`
	SymptomID int `db:"symptom_id"`
	Weight    int `db:"weight"`
}

type UserAnswers struct {
	Answers map[int]bool `json:"answers" validate:"required"`
}

func GetQuestions(db *sqlx.DB) ([]Question, error) {
	var questions []Question
	err := db.Select(&questions, "SELECT id, question_text FROM questions")
	if err != nil {
		logrus.WithError(err).Error("Error selecting questions")
		return nil, errors.Wrap(err, "selecting questions")
	}
	return questions, nil
}

func GetDiseaseSymptoms(db *sqlx.DB) (map[int]map[int]int, error) {
	var diseaseSymptoms []DiseaseSymptom
	err := db.Select(&diseaseSymptoms, "SELECT disease_id, symptom_id, weight FROM disease_symptoms")
	if err != nil {
		logrus.WithError(err).Error("Error selecting disease symptoms")
		return nil, errors.Wrap(err, "selecting disease symptoms")
	}

	diseaseSymptomMap := make(map[int]map[int]int)
	for _, ds := range diseaseSymptoms {
		if _, exists := diseaseSymptomMap[ds.DiseaseID]; !exists {
			diseaseSymptomMap[ds.DiseaseID] = make(map[int]int)
		}
		diseaseSymptomMap[ds.DiseaseID][ds.SymptomID] = ds.Weight
	}

	return diseaseSymptomMap, nil
}

func GetDiseases(db *sqlx.DB) (map[int]string, error) {
	var diseases []Disease
	err := db.Select(&diseases, "SELECT id, name FROM diseases")
	if err != nil {
		logrus.WithError(err).Error("Error selecting diseases")
		return nil, errors.Wrap(err, "selecting diseases")
	}

	diseaseMap := make(map[int]string)
	for _, disease := range diseases {
		diseaseMap[disease.ID] = disease.Name
	}

	return diseaseMap, nil
}

func SaveUserAnswers(db *sqlx.DB, answers UserAnswers) error {
	for id, answer := range answers.Answers {
		_, err := db.Exec("INSERT INTO user_answers (question_id, answer) VALUES ($1, $2)", id, answer)
		if err != nil {
			logrus.WithError(err).Error("Error saving user answers")
			return errors.Wrap(err, "saving user answers")
		}
	}
	return nil
}

func AnalyzeAnswers(db *sqlx.DB, answers UserAnswers) (string, error) {
	diseaseSymptoms, err := GetDiseaseSymptoms(db)
	if err != nil {
		logrus.WithError(err).Error("Error getting disease symptoms")
		return "", err
	}

	diseaseMap, err := GetDiseases(db)
	if err != nil {
		logrus.WithError(err).Error("Error getting diseases")
		return "", err
	}

	scores := make(map[int]int)
	for diseaseID, symptomWeights := range diseaseSymptoms {
		for symptomID, weight := range symptomWeights {
			if answers.Answers[symptomID] {
				scores[diseaseID] += weight
			}
		}
	}

	var maxDiseaseID int
	var maxScore int
	for diseaseID, score := range scores {
		if score > maxScore {
			maxDiseaseID = diseaseID
			maxScore = score
		}
	}

	if maxScore == 0 {
		logrus.Warn("No diagnosis could be made based on the symptoms provided")
		return "No diagnosis could be made based on the symptoms provided.", nil
	}

	logrus.Infof("Diagnosis determined: %s", diseaseMap[maxDiseaseID])
	return diseaseMap[maxDiseaseID], nil
}
