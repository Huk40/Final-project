package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"goFinalProject/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "метод не поддерживается")
		return
	}

	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeError(w, "не указан идентификатор задачи")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeDoneError(w, err)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
	} else {
		var next string
		next, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err == nil {
			err = db.UpdateDate(next, id)
		}
	}
	if err != nil {
		writeDoneError(w, err)
		return
	}

	writeJSON(w, struct{}{})
}

func writeDoneError(w http.ResponseWriter, err error) {
	if errors.Is(err, db.ErrTaskNotFound) {
		writeError(w, "задача не найдена")
	} else {
		writeError(w, fmt.Sprintf("не удалось выполнить задачу: %v", err))
	}
}
