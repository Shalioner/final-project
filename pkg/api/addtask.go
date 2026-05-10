package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"final-project/pkg/db"
)

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(dateFormat)

	if task.Date == "" || task.Date == "today" {
		task.Date = today
	}

	if _, err := time.Parse(dateFormat, task.Date); err != nil {
		return err
	}

	if task.Date < today {
		if task.Repeat == "" {
			task.Date = today
			return nil
		}
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = next
		return nil
	}

	if task.Repeat != "" {
		_, err := NextDate(now, task.Date, task.Repeat)
		return err
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}
