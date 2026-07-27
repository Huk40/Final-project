package api

import (
	"errors"
	"net/http"
	"strings"

	"goFinalProject/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeError(w, "не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, "задача не найдена", http.StatusNotFound)
		} else {
			writeInternalError(w, "не удалось удалить задачу", err)
		}
		return
	}

	writeJSON(w, struct{}{}, http.StatusOK)
}
