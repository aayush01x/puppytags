package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pclubiitk/puppylove_tags/services"
)

type Router struct {
	UserSimilarityService *services.UserSimilarityService
}

func NewRouter(uss *services.UserSimilarityService) *Router {
	return &Router{UserSimilarityService: uss}
}

func (r *Router) RegisterRoutes() {
	http.HandleFunc("/user", r.updateUserHandler)
	http.HandleFunc("/similar", r.querySimilarHandler)
}


func (r *Router) updateUserHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	type updateRequest struct {
		UserID string `json:"user_id"`
		Tags   []int  `json:"tags"`
	}

	var data updateRequest
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if data.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	if err := r.UserSimilarityService.UpdateUser(data.UserID, data.Tags); err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s updated successfully\n", data.UserID)
}


func (r *Router) querySimilarHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	userID := req.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}
	offsetStr := req.URL.Query().Get("offset")
	limitStr := req.URL.Query().Get("limit")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	similarUsers, err := r.UserSimilarityService.QuerySimilar(userID, offset, limit)
	if err != nil {
		http.Error(w, "Query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	resp := map[string][]string{"similar_users": similarUsers}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
