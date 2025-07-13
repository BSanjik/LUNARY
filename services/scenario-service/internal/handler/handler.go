// HTTP handlers
package handler

import (
	"context"
	"encoding/json"
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
	var req scenarioRequest

	//read json
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	//check request
	if req.Query == "" {
		http.Error(w, "`query` is required", http.StatusBadRequest)
		return
	}

	//get scenario
	scenario, err := h.Service.GetScenario(context.Background(), req.Query)
	if err != nil {
		http.Error(w, "scenario not found", http.StatusNotFound)
		return
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scenario)
}

// POST Feedback
func (h *Handler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	var req model.FeedbackRequest

	//read json
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
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
