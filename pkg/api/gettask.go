package api

import (
	"errors"
	"net/http"
	"strings"

	"goFinalProject/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeError(w, "не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, "задача не найдена", http.StatusNotFound)
		} else {
			writeInternalError(w, "не удалось получить задачу", err)
		}
		return
	}

	writeJSON(w, task, http.StatusOK)
}
