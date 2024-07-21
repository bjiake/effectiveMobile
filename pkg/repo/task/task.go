package task

import (
	"context"
	"database/sql"
	"effectiveMobile/pkg/db"
	"effectiveMobile/pkg/domain/task"
	interfaces "effectiveMobile/pkg/repo/task/interface"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"log"
	"strconv"
	"strings"
	"time"
)

type taskDataBase struct {
	db *sql.DB
}

func NewTaskDataBase(db *sql.DB) interfaces.TaskRepository {
	return &taskDataBase{
		db: db,
	}
}

func (r *taskDataBase) Migrate(ctx context.Context) error {
	accQuery := `
    CREATE TABLE IF NOT EXISTS task (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
 		startTime TIMESTAMP NOT NULL,
		endTime TIMESTAMP,
		totalTime INTERVAL
	);
    `
	_, err := r.db.ExecContext(ctx, accQuery)
	if err != nil {
		message := db.ErrMigrate.Error() + " book"
		log.Printf("%q: %s\n", message, err.Error())
		return db.ErrMigrate
	}

	return err
}

func (r *taskDataBase) Post(ctx context.Context, newTask task.Task) (*task.Task, error) {
	var id int64

	err := r.db.QueryRowContext(ctx, "INSERT INTO task(name, description, startTime) values($1, $2, $3) RETURNING id", newTask.Name, newTask.Description, newTask.StartTime).Scan(&id)
	// Check if a task with the same books already exists
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, db.ErrDuplicate
			}
		}
		return nil, err
	}

	// Add the new task
	requestTask := &task.Task{
		ID:          id,
		Name:        newTask.Name,
		Description: newTask.Description,
		StartTime:   newTask.StartTime,
	}

	return requestTask, nil
}

func (r *taskDataBase) Put(ctx context.Context, id int64, updateTask task.Task) (*task.Task, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE task SET name = $1, description = $2, startTime = $3, endTime = $4, totalTime = $5 WHERE task.id = $6",
		updateTask.Name, updateTask.Description, updateTask.StartTime, updateTask.EndTime, updateTask.TotalTime, id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, db.ErrDuplicate
			}
		}
		return nil, err
	}

	result := &task.Task{
		ID:          id,
		Name:        updateTask.Name,
		Description: updateTask.Description,
		StartTime:   updateTask.StartTime,
		EndTime:     updateTask.EndTime,
		TotalTime:   updateTask.TotalTime,
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, db.ErrUpdateFailed
	}

	return result, nil
}

func (r *taskDataBase) Get(ctx context.Context, id int64) (*task.Task, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM task WHERE id = $1", id)
	var result task.Task
	if err := row.Scan(&result.ID, &result.Name, &result.Description, &result.StartTime, &result.EndTime, &result.TotalTime); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, db.ErrNotExist
		}
		return nil, err
	}

	return &result, nil
}

func (r *taskDataBase) GetLaborCost(ctx context.Context, startTime time.Time, endTime time.Time) (task.Slice, error) {
	query := `
        SELECT id, name, description, startTime, endTime, totalTime
        FROM task
        WHERE startTime >= $1 AND startTime <= $2
        AND (endTime IS NOT NULL AND endTime <= $2)
    `

	rows, err := r.db.QueryContext(ctx, query, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks task.Slice
	for rows.Next() {
		var t task.Task
		var endTime sql.NullTime
		var totalTimeStr sql.NullString

		err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.StartTime, &endTime, &totalTimeStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task row: %w", err)
		}

		if endTime.Valid {
			t.EndTime = &endTime.Time
		}

		if totalTimeStr.Valid {
			duration, err := parseDuration(totalTimeStr.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse duration: %w", err)
			}
			t.TotalTime = &duration
		}

		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tasks, nil
}

func parseDuration(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid duration format")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hours: %w", err)
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %w", err)
	}

	seconds, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid seconds: %w", err)
	}

	duration := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds*float64(time.Second))

	return duration, nil
}
func (r *taskDataBase) GetAll(ctx context.Context) ([]task.Task, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM task")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []task.Task
	for rows.Next() {
		var t task.Task
		var endTime sql.NullTime
		var totalTime sql.NullString

		if err = rows.Scan(&t.ID, &t.Name, &t.Description, &t.StartTime, &endTime, &totalTime); err != nil {
			return nil, err
		}

		if endTime.Valid {
			t.EndTime = &endTime.Time
		}

		if totalTime.Valid {
			duration, err := parseDuration(totalTime.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse duration: %w", err)
			}
			t.TotalTime = &duration
		}

		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
func (r *taskDataBase) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM task WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return db.ErrDeleteFailed
	}

	return err
}
