package api

import (
	"errors"
	"net/http"
	"strings"

	"goFinalProject/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := decodeTask(r)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(task.ID) == "" {
		writeError(w, "не указан идентификатор задачи", http.StatusBadRequest)
		return
	}
	if err := validateTask(task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(task); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, "задача не найдена", http.StatusNotFound)
		} else {
			writeInternalError(w, "не удалось обновить задачу", err)
		}
		return
	}

	writeJSON(w, struct{}{}, http.StatusOK)
}
