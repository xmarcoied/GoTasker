zpackage tasker

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type TaskStore struct {
	DB *sql.DB
}

func (ts *TaskStore) Save(t *Task) (*Task, error) {
	var cmd string
	if t.ID == 0 {
		// Create a new task.
		cmd = fmt.Sprintf("INSERT INTO tasks (name, action, time) VALUES ('%s', '%s', '%s')",
			t.Name, t.Action, t.ScheduledTime)
	} else {
		cmd = fmt.Sprintf("UPDATE tasks t SET t.name='%s', t.action='%s', t.time='%s' WHERE id = %d",
			t.Name, t.Action, t.ScheduledTime, t.ID)
	}

	insert, err := ts.DB.Exec(cmd)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	id, err := insert.LastInsertId()
	if t.ID == 0 && err != nil {
		log.Println(err)
	} else if t.ID == 0 {
		t.ID = id
	}

	return t, nil
}
