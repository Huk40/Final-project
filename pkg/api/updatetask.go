package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"goFinalProject/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := decodeTask(r)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	if strings.TrimSpace(task.ID) == "" {
		writeError(w, "не указан идентификатор задачи")
		return
	}
	if err := validateTask(task); err != nil {
		writeError(w, err.Error())
		return
	}

	if err := db.UpdateTask(task); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, "задача не найдена")
		} else {
			writeError(w, fmt.Sprintf("не удалось обновить задачу: %v", err))
		}
		return
	}

	writeJSON(w, struct{}{})
}
