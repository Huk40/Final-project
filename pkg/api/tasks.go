package api

import (
	"fmt"
	"net/http"

	"goFinalProject/pkg/db"
)

const tasksLimit = 50

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "метод не поддерживается")
		return
	}

	tasks, err := db.Tasks(tasksLimit, r.URL.Query().Get("search"))
	if err != nil {
		writeError(w, fmt.Sprintf("не удалось получить список задач: %v", err))
		return
	}

	writeJSON(w, TasksResponse{Tasks: tasks})
}
