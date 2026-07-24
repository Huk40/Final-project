package db

import (
	"database/sql"
	"errors"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	if database == nil {
		return 0, errors.New("database is not initialized")
	}

	result, err := database.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	if database == nil {
		return nil, errors.New("database is not initialized")
	}

	task := new(Task)
	err := database.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
		id,
	).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTaskNotFound
	}
	if err != nil {
		return nil, err
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	if database == nil {
		return errors.New("database is not initialized")
	}

	result, err := database.Exec(
		`UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?`,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)
	if err != nil {
		return err
	}

	return checkAffected(result)
}

func UpdateDate(next string, id string) error {
	if database == nil {
		return errors.New("database is not initialized")
	}

	result, err := database.Exec(
		`UPDATE scheduler SET date = ? WHERE id = ?`,
		next,
		id,
	)
	if err != nil {
		return err
	}

	return checkAffected(result)
}

func DeleteTask(id string) error {
	if database == nil {
		return errors.New("database is not initialized")
	}

	result, err := database.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}

	return checkAffected(result)
}

func checkAffected(result sql.Result) error {
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	tasks := make([]*Task, 0)
	if database == nil {
		return tasks, errors.New("database is not initialized")
	}
	if limit <= 0 {
		return tasks, nil
	}

	query := `SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date, id
		LIMIT ?`
	args := []any{limit}

	if search != "" {
		if date, err := time.Parse("02.01.2006", search); err == nil {
			query = `SELECT id, date, title, comment, repeat
				FROM scheduler
				WHERE date = ?
				ORDER BY date, id
				LIMIT ?`
			args = []any{date.Format("20060102"), limit}
		} else {
			query = `SELECT id, date, title, comment, repeat
				FROM scheduler
				WHERE title LIKE ? OR comment LIKE ?
				ORDER BY date, id
				LIMIT ?`
			pattern := "%" + search + "%"
			args = []any{pattern, pattern, limit}
		}
	}

	rows, err := database.Query(query, args...)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		task := new(Task)
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil
}
