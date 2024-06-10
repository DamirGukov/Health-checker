package handlers

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"healthchecker-app/models"
	"net/http"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}

func (h *Handler) GetQuestions(c echo.Context) error {
	questions, err := models.GetQuestions(h.db)
	if err != nil {
		logrus.WithError(err).Error("Error getting questions")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, questions)
}

func (h *Handler) SubmitAnswers(c echo.Context) error {
	var req models.UserAnswers
	if err := c.Bind(&req); err != nil {
		logrus.WithError(err).Error("Error binding request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	logrus.Infof("Received answers: %+v", req) // Отладочная информация

	if err := models.SaveUserAnswers(h.db, req); err != nil {
		logrus.WithError(err).Error("Error saving user answers")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	diagnosis, err := models.AnalyzeAnswers(h.db, req)
	if err != nil {
		logrus.WithError(err).Error("Error analyzing answers")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"diagnosis": diagnosis})
}
