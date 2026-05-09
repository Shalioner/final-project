package db

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID      int64  `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

func (t Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"id":      strconv.FormatInt(t.ID, 10),
		"date":    t.Date,
		"title":   t.Title,
		"comment": t.Comment,
		"repeat":  t.Repeat,
	})
}

func AddTask(task *Task) (int64, error) {
	res, err := DB.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int, search string) ([]*Task, error) {
	tasks := make([]*Task, 0)
	var err error

	if search != "" {
		if t, parseErr := time.Parse("02.01.2006", search); parseErr == nil {
			dateStr := t.Format("20060102")
			err = DB.Select(&tasks, `SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`, dateStr, limit)
		} else {
			like := "%" + search + "%"
			err = DB.Select(&tasks, `SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`, like, like, limit)
		}
	} else {
		err = DB.Select(&tasks, `SELECT * FROM scheduler ORDER BY date LIMIT ?`, limit)
	}

	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task
	err := DB.Get(&task, `SELECT * FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	res, err := DB.Exec(
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	res, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
