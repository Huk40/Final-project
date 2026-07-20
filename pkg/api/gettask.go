package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"goFinalProject/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeError(w, "не указан идентификатор задачи")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, "задача не найдена")
		} else {
			writeError(w, fmt.Sprintf("не удалось получить задачу: %v", err))
		}
		return
	}

	writeJSON(w, task)
}
