package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	tasks := make([]*Task, 0)

	var (
		query string
		args  []any
	)

	switch {
	case search == "":
		query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			ORDER BY date ASC
			LIMIT ?
		`
		args = append(args, limit)

	default:
		searchDate, err := time.Parse("02.01.2006", search)

		if err == nil {
			query = `
				SELECT id, date, title, comment, repeat
				FROM scheduler
				WHERE date = ?
				ORDER BY date ASC
				LIMIT ?
			`

			args = append(args, searchDate.Format("20060102"), limit)
		} else {
			query = `
				SELECT id, date, title, comment, repeat
				FROM scheduler
				WHERE title LIKE ? OR comment LIKE ?
				ORDER BY date ASC
				LIMIT ?
			`

			pattern := "%" + search + "%"
			args = append(args, pattern, pattern, limit)
		}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

		if err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task

	err := DB.QueryRow(
		`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?`,
		id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Задача не найдена")
		}
		return nil, err
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`

	res, err := DB.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}
