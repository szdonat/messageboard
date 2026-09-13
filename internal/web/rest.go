package web

import (
	"encoding/json"
	"log"
	"net/http"

	"message-board/internal/repository"
)

type Rest struct {
	repository repository.Repository
	mux        *http.ServeMux
}

func NewRest(repo repository.Repository) *Rest {
	r := &Rest{
		repository: repo,
		mux:        http.NewServeMux(),
	}

	r.mux.HandleFunc("GET /healthz", r.healthzHandler)
	r.mux.HandleFunc("POST /messages", r.createMessageHandler)
	r.mux.HandleFunc("GET /messages", r.listMessagesHandler)
	r.mux.HandleFunc("POST /messages/{messageID}/replies", r.createReplyHandler)
	r.mux.HandleFunc("GET /messages/{messageID}/replies", r.listRepliesHandler)

	return r
}

func (r *Rest) Router() http.Handler {
	return r.mux
}

func (r *Rest) healthzHandler(w http.ResponseWriter, _ *http.Request) {
	response, err := r.repository.Ping()
	if err != nil {
		log.Printf("health check failed: %v", err)
		http.Error(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	if response != "PONG" {
		log.Printf("health check failed: unexpected Redis response %q", response)
		http.Error(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (r *Rest) createMessageHandler(w http.ResponseWriter, req *http.Request) {
	var input struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	messageID, err := r.repository.CreateMessage(input.Title, input.Body)
	if err != nil {
		log.Printf("create message failed: %v", err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, struct {
		ID string `json:"id"`
	}{ID: messageID})
}

func (r *Rest) listMessagesHandler(w http.ResponseWriter, _ *http.Request) {
	messages, err := r.repository.ListMessages()
	if err != nil {
		log.Printf("list messages failed: %v", err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeRawJSON(w, http.StatusOK, messages)
}

func (r *Rest) createReplyHandler(w http.ResponseWriter, req *http.Request) {
	var input struct {
		Reply string `json:"reply"`
	}

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := r.repository.CreateReply(req.PathValue("messageID"), input.Reply); err != nil {
		log.Printf("create reply failed: %v", err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "created"})
}

func (r *Rest) listRepliesHandler(w http.ResponseWriter, req *http.Request) {
	replies, err := r.repository.ListReplies(req.PathValue("messageID"))
	if err != nil {
		log.Printf("list replies failed: %v", err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeRawJSON(w, http.StatusOK, replies)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeRawJSON(w http.ResponseWriter, status int, value string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(value))
}
