package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"goFinalProject/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeError(w, "не указан идентификатор задачи")
		return
	}

	if err := db.DeleteTask(id); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, "задача не найдена")
		} else {
			writeError(w, fmt.Sprintf("не удалось удалить задачу: %v", err))
		}
		return
	}

	writeJSON(w, struct{}{})
}
