package task

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresTaskStore struct {
	db *sql.DB
}

func NewPostgresTaskStore(dsn string) (*PostgresTaskStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database, %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping database, %w", err)
	}

	return &PostgresTaskStore{db: db}, nil
}

func (p *PostgresTaskStore) GetAllTasks() (List, error) {

	rows, err := p.db.Query(`
		SELECT id, title, description
		FROM tasks
		ORDER BY id ASC
	`)

	if err != nil {
		return nil, fmt.Errorf("GetAllTasks query failed, %w", err)
	}

	defer rows.Close()

	// var tasks List
	// changed to below because we want to avoid nil slices and have an empty slice instead when there are no tasks
	tasks := make(List, 0)

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description); err != nil {
			return nil, fmt.Errorf("GetAllTasks scan failed, %w", err)
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllTasks rows error, %w", err)
	}

	return tasks, nil

}

func (p *PostgresTaskStore) AddTask(task Task) (int, error) {
	var newId int

	err := p.db.QueryRow(
		`INSERT INTO tasks (title, description)
		VALUES ($1, $2)
		RETURNING id`,
		task.Title, task.Description).Scan(&newId)

	if err != nil {
		return 0, fmt.Errorf("AddTask query failed, %w", err)
	}

	return newId, nil
}

func (p *PostgresTaskStore) DeleteTask(id int) (bool, error) {
	result, err := p.db.Exec(`DELETE FROM tasks where ID = $1`, id)
	if err != nil {
		return false, fmt.Errorf("DeleteTask query failed, %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("DeleteTask RowsAffected check failed, %w", err)
	}

	// rowsAffected == 0 means no rows matched
	return rowsAffected > 0, nil
}

func (p *PostgresTaskStore) UpdateTask(id int, task Task) (bool, error) {
	result, err := p.db.Exec(`
		UPDATE tasks
		SET title = $1, 
		description = $2, 
		updated_at = NOW()
		WHERE id = $3`,
		task.Title, task.Description, id)

	if err != nil {
		return false, fmt.Errorf("UpdateTask query failed, %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("UpdateTask rows affected check failed %w", err)
	}

	return rowsAffected > 0, nil

}
