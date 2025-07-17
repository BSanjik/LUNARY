// HTTP handlers
package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"scenario-service/internal/model"
	"scenario-service/internal/service"
)

type Handler struct {
	Service *service.ScenarioService
}

type scenarioRequest struct {
	Query string `json:"query"`
}

// POST Scenario
func (h *Handler) GetScenario(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 Scenario-service: получен запрос %s %s", r.Method, r.URL.Path)
	log.Printf("🔍 Scenario-service: заголовки: %v", r.Header)

	var req scenarioRequest

	//read json
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Scenario-service: ошибка декодирования JSON: %v", err)
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 Scenario-service: запрос query = '%s'", req.Query)

	//check request
	if req.Query == "" {
		log.Printf("❌ Scenario-service: пустой query")
		http.Error(w, "Поле 'query' обязательно для заполнения", http.StatusBadRequest)
		return
	}

	//get scenario
	scenario, err := h.Service.GetScenario(context.Background(), req.Query)
	if err != nil {
		log.Printf("❌ Scenario-service: сценарий не найден для query '%s': %v", req.Query, err)
		http.Error(w, "Сценарий не найден. Попробуйте изменить запрос", http.StatusNotFound)
		return
	}

	log.Printf("✅ Scenario-service: найден сценарий '%s' для query '%s'", scenario.Title, req.Query)

	//response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scenario)

	log.Printf("✅ Scenario-service: ответ отправлен успешно")
}

// POST Feedback
func (h *Handler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	var req model.FeedbackRequest

	//read json
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный формат данных", http.StatusBadRequest)
		return
	}

	// Сохраняем фидбек
	if err := h.Service.SubmitFeedback(context.Background(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Feedback submitted successfully",
	})
}
