package people

import (
	"context"
	"database/sql"
	"effectiveMobile/pkg/db"
	"effectiveMobile/pkg/domain/people"
	"effectiveMobile/pkg/domain/task"
	interfaces "effectiveMobile/pkg/repo/people/interface"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgconn"
)

type accountDataBase struct {
	db *sql.DB
}

func NewPeopleDataBase(db *sql.DB) interfaces.PeopleRepository {
	return &accountDataBase{
		db: db,
	}
}

func (r *accountDataBase) Migrate(ctx context.Context) error {
	accQuery := `
    CREATE TABLE IF NOT EXISTS people (
		id SERIAL PRIMARY KEY,
		name TEXT ,
		surName TEXT ,
		patronymic TEXT ,
		address TEXT ,
		tasks INTEGER[],
		passportNumber TEXT NOT NULL,
		password TEXT NOT NULL
	);
    `
	_, err := r.db.ExecContext(ctx, accQuery)
	if err != nil {
		message := db.ErrMigrate.Error() + " people"
		log.Printf("%q: %s\n", message, err.Error())
		return db.ErrMigrate
	}

	return err
}

func (r *accountDataBase) Info(ctx context.Context, passportNumber string) (*people.Info, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM people WHERE passportNumber = $1", passportNumber)

	var result people.Info
	var nameNull, surnameNull, patronymicNull, addressNull sql.NullString
	var tasksArray sql.NullString
	var password *string

	if err := row.Scan(
		&result.ID,
		&nameNull,
		&surnameNull,
		&patronymicNull,
		&addressNull,
		&tasksArray,
		&passportNumber,
		&password,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, db.ErrNotExist
		}
		return nil, err
	}

	// Set values or defaults for potentially nil fields
	if nameNull.Valid {
		result.Name = nameNull.String
	} else {
		result.Name = "Unknown"
	}

	if surnameNull.Valid {
		result.Surname = surnameNull.String
	} else {
		result.Surname = "Unknown"
	}

	if patronymicNull.Valid {
		result.Patronymic = patronymicNull.String
	} else {
		result.Patronymic = "Unknown"
	}

	if addressNull.Valid {
		result.Address = addressNull.String
	} else {
		result.Address = "Unknown"
	}

	return &result, nil
}

func (r *accountDataBase) Registration(ctx context.Context, newPeople people.Registration) (*int64, error) {
	var existingCount int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM people WHERE passportNumber = $1", newPeople.PassportNumber).Scan(&existingCount)
	if err != nil {
		return nil, err
	}

	if existingCount > 0 {
		return nil, db.ErrDuplicate
	}

	var id int64

	err = r.db.QueryRowContext(ctx,
		"INSERT INTO people(passportNumber, password) values($1, $2) RETURNING id",
		newPeople.PassportNumber, newPeople.Password).Scan(&id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, db.ErrDuplicate
			}
		}
		return nil, err
	}

	return &id, nil
}

func (r *accountDataBase) Login(ctx context.Context, acc people.Registration) (int64, error) {
	var id int64
	row := r.db.QueryRowContext(ctx, "SELECT id FROM people WHERE passportNumber = $1 and password = $2", acc.PassportNumber, acc.Password)

	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, db.ErrNotExist
		}
		return 0, err
	}
	return id, nil
}

func (r *accountDataBase) Get(ctx context.Context, filter *people.Filter, pagination *people.Pagination) ([]people.Request, error) {
	query := "SELECT id, name, surname, patronymic, address, tasks, passportNumber FROM people"
	var args []interface{}

	if filter != nil {
		var whereClauses []string
		if filter.ID != nil {
			args = append(args, *filter.ID)
			whereClauses = append(whereClauses, fmt.Sprintf("id = $%d", len(args)))
		}
		if filter.Name != nil {
			args = append(args, *filter.Name)
			whereClauses = append(whereClauses, fmt.Sprintf("name = $%d", len(args)))
		}
		if filter.Surname != nil {
			args = append(args, *filter.Surname)
			whereClauses = append(whereClauses, fmt.Sprintf("surname = $%d", len(args)))
		}
		if filter.Patronymic != nil {
			args = append(args, *filter.Patronymic)
			whereClauses = append(whereClauses, fmt.Sprintf("patronymic = $%d", len(args)))
		}
		if filter.Address != nil {
			args = append(args, *filter.Address)
			whereClauses = append(whereClauses, fmt.Sprintf("address = $%d", len(args)))
		}
		if filter.Tasks != nil {
			args = append(args, pq.Array(*filter.Tasks))
			whereClauses = append(whereClauses, fmt.Sprintf("tasks @> $%d", len(args)))
		}

		if len(whereClauses) > 0 {
			query += " WHERE " + strings.Join(whereClauses, " AND ")
		}
	}

	// Apply pagination
	if pagination != nil {
		if pagination.Limit > 0 {
			args = append(args, pagination.Limit)
			query += fmt.Sprintf(" LIMIT $%d", len(args))
		}
		if pagination.Offset > 0 {
			args = append(args, pagination.Offset)
			query += fmt.Sprintf(" OFFSET $%d", len(args))
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []people.Request
	for rows.Next() {
		var (
			id             sql.NullInt64
			name           sql.NullString
			surname        sql.NullString
			patronymic     sql.NullString
			address        sql.NullString
			tasksArray     sql.NullString
			passportNumber string
		)

		if err := rows.Scan(&id, &name, &surname, &patronymic, &address, &tasksArray, &passportNumber); err != nil {
			return nil, err
		}

		requestPeople := people.Request{
			PassportNumber: passportNumber,
		}

		if id.Valid {
			requestPeople.ID = id.Int64
		}
		if name.Valid {
			requestPeople.Name = name.String
		}
		if surname.Valid {
			requestPeople.Surname = surname.String
		}
		if patronymic.Valid {
			requestPeople.Patronymic = patronymic.String
		}
		if address.Valid {
			requestPeople.Address = address.String
		}

		if tasksArray.Valid {
			taskIDs, err := parsePostgresArray(tasksArray.String)
			if err != nil {
				return nil, err
			}
			requestPeople.Tasks = make([]task.Task, 0, len(taskIDs))
			for _, taskID := range taskIDs {
				currTask, err := r.getTasks(ctx, taskID)
				if err != nil {
					// You might want to log this error or handle it differently
					continue
				}
				requestPeople.Tasks = append(requestPeople.Tasks, currTask)
			}
		}

		result = append(result, requestPeople)
	}
	return result, nil
}

func (r *accountDataBase) getTasks(ctx context.Context, taskID int64) (task.Task, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, description, startTime, endTime, totalTime FROM task WHERE id = $1", taskID)

	var (
		currTask    task.Task
		name        sql.NullString
		description sql.NullString
		startTime   sql.NullTime
		endTime     sql.NullTime
		totalTime   sql.NullInt64
	)

	err := row.Scan(&currTask.ID, &name, &description, &startTime, &endTime, &totalTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return task.Task{}, db.ErrTasks
		}
		return task.Task{}, err
	}

	if name.Valid {
		currTask.Name = name.String
	}
	if description.Valid {
		currTask.Description = description.String
	}
	if startTime.Valid {
		currTask.StartTime = startTime.Time
	}
	if endTime.Valid {
		currTask.EndTime = &endTime.Time
	}
	if totalTime.Valid {
		duration := time.Duration(totalTime.Int64)
		currTask.TotalTime = &duration
	}

	return currTask, nil
}

func parsePostgresArray(arrayStr string) ([]int64, error) {
	// Remove the curly braces
	arrayStr = strings.Trim(arrayStr, "{}")

	// Split the string by comma
	strNumbers := strings.Split(arrayStr, ",")

	result := make([]int64, 0, len(strNumbers))

	for _, strNum := range strNumbers {
		num, err := strconv.ParseInt(strings.TrimSpace(strNum), 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}

	return result, nil
}

func (r *accountDataBase) Put(ctx context.Context, id int64, updatePeople people.Info) (*people.Info, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE people SET name = $1, surname = $2, patronymic = $3, address = $4 WHERE people.id = $5",
		updatePeople.Name, updatePeople.Surname, updatePeople.Patronymic, updatePeople.Address, id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, db.ErrDuplicate
			}
		}
		return nil, err
	}

	result := &people.Info{
		ID:         id,
		Name:       updatePeople.Name,
		Surname:    updatePeople.Surname,
		Patronymic: updatePeople.Patronymic,
		Address:    updatePeople.Address,
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

func (r *accountDataBase) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM people WHERE id = $1", id)
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

func (r *accountDataBase) AppendTask(ctx context.Context, id int64, task task.Task) error {
	_, err := r.db.ExecContext(ctx, "UPDATE people SET tasks = array_append(tasks, $1) WHERE people.id = $2", task.ID, id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return db.ErrDuplicate
			}
		}
		return err
	}
	return nil
}
