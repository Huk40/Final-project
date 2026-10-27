package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"goFinalProject/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := decodeTask(r)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validateTask(task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(task)
	if err != nil {
		writeInternalError(w, "не удалось добавить задачу", err)
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)}, http.StatusCreated)
}

func decodeTask(r *http.Request) (*db.Task, error) {
	task := new(db.Task)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(task); err != nil {
		return nil, fmt.Errorf("не удалось прочитать JSON: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("тело запроса должно содержать один JSON-объект")
	}
	return task, nil
}

func validateTask(task *db.Task) error {
	if strings.TrimSpace(task.Title) == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}
	return checkDate(task)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	date, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("неверный формат даты: %w", err)
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("неверное правило повторения: %w", err)
		}
	}

	if after(now, date) {
		if task.Repeat == "" {
			task.Date = now.Format(DateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}
