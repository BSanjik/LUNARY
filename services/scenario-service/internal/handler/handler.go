// HTTP handlers
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"scenario-service/internal/model"
	"scenario-service/internal/service"
	"strings"

	"github.com/BSanjik/LUNARY/services/common/shared"
)

type Handler struct {
	Service service.ScenarioServiceInterface
}

type scenarioRequest struct {
	Query string `json:"query"`
}

// POST Scenario
func (h *Handler) GetScenario(w http.ResponseWriter, r *http.Request) {
	var req scenarioRequest

	// Декодирование JSON
	if err := h.decodeRequest(r, &req); err != nil {
		shared.WriteErrorResponse(w, shared.NewValidationError("Некорректный формат запроса"))
		return
	}

	// Очистка входных данных
	h.trimStringFields(&req)

	// Поиск сценариев через сервис (валидация внутри)
	scenario, err := h.Service.SearchScenarios(context.Background(), req.Query)
	if err != nil {
		shared.WriteErrorResponse(w, err)
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scenario)
}

// POST Feedback
func (h *Handler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	var req model.FeedbackRequest

	// Декодирование JSON
	if err := h.decodeRequest(r, &req); err != nil {
		shared.WriteErrorResponse(w, shared.NewValidationError("Некорректный формат данных"))
		return
	}

	// Сохранение фидбека через сервис (валидация внутри)
	if err := h.Service.SubmitFeedback(context.Background(), &req); err != nil {
		shared.WriteErrorResponse(w, err)
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Отзыв успешно сохранён",
	})
}

// Вспомогательные функции
func (h *Handler) decodeRequest(r *http.Request, target interface{}) error {
	return json.NewDecoder(r.Body).Decode(target)
}

func (h *Handler) trimStringFields(req *scenarioRequest) {
	req.Query = strings.TrimSpace(req.Query)
}
