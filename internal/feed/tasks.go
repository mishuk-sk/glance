package feed

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Task struct {
	ID       int       `json:"id"`
	Text     string    `json:"text"`
	DueDate  time.Time `json:"due_date"`
	Priority int       `json:"priority"`
	Done     bool      `json:"done"`
}

type Tasks []Task

func FetchTasks(baseURL string, status string) (Tasks, error) {
	switch status {
	case "done", "not_done", "any":
	default:
		status = "any"
	}
	endpoint := fmt.Sprintf("%s/tasks?sort_by=due_date:asc,priority:desc&status=%s", baseURL, status)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", resp.Status)
	}

	var rawTasks []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawTasks); err != nil {
		return nil, err
	}

	var tasks Tasks
	for _, rawTask := range rawTasks {
		task := Task{
			ID:       int(rawTask["id"].(float64)),
			Text:     rawTask["text"].(string),
			Priority: int(rawTask["priority"].(float64)),
			Done:     rawTask["done"].(bool),
		}

		if dueDateStr, ok := rawTask["due_date"].(string); ok && dueDateStr != "" {
			dueDate, err := time.Parse(time.RFC3339, dueDateStr)
			if err != nil {
				return nil, err
			}
			task.DueDate = dueDate
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
