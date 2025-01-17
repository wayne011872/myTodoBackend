// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package todoItem

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteTodoItem = `-- name: DeleteTodoItem :exec
DELETE FROM todoItem WHERE id = $1
`

func (q *Queries) DeleteTodoItem(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteTodoItem, id)
	return err
}

const getAllTodoItem = `-- name: GetAllTodoItem :many
SELECT id, title, detail, completed, starttime, endtime, createdtime, updatedtime FROM todoItem ORDER BY title
`

func (q *Queries) GetAllTodoItem(ctx context.Context) ([]Todoitem, error) {
	rows, err := q.db.Query(ctx, getAllTodoItem)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Todoitem
	for rows.Next() {
		var i Todoitem
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Detail,
			&i.Completed,
			&i.Starttime,
			&i.Endtime,
			&i.Createdtime,
			&i.Updatedtime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTodoItemById = `-- name: GetTodoItemById :one
SELECT id, title, detail, completed, starttime, endtime, createdtime, updatedtime FROM todoItem
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTodoItemById(ctx context.Context, id int64) (Todoitem, error) {
	row := q.db.QueryRow(ctx, getTodoItemById, id)
	var i Todoitem
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Detail,
		&i.Completed,
		&i.Starttime,
		&i.Endtime,
		&i.Createdtime,
		&i.Updatedtime,
	)
	return i, err
}

const getTodoItemByTitle = `-- name: GetTodoItemByTitle :one
SELECT id, title, detail, completed, starttime, endtime, createdtime, updatedtime FROM todoItem
WHERE title = $1 LIMIT 1
`

func (q *Queries) GetTodoItemByTitle(ctx context.Context, title string) (Todoitem, error) {
	row := q.db.QueryRow(ctx, getTodoItemByTitle, title)
	var i Todoitem
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Detail,
		&i.Completed,
		&i.Starttime,
		&i.Endtime,
		&i.Createdtime,
		&i.Updatedtime,
	)
	return i, err
}

const insertTodoItem = `-- name: InsertTodoItem :exec
INSERT INTO todoItem (title, detail, completed, startTime, endTime, createdTime, updatedTime) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

type InsertTodoItemParams struct {
	Title       string
	Detail      pgtype.Text
	Completed   bool
	Starttime   pgtype.Timestamp
	Endtime     pgtype.Timestamp
	Createdtime pgtype.Timestamp
	Updatedtime pgtype.Timestamp
}

func (q *Queries) InsertTodoItem(ctx context.Context, arg InsertTodoItemParams) error {
	_, err := q.db.Exec(ctx, insertTodoItem,
		arg.Title,
		arg.Detail,
		arg.Completed,
		arg.Starttime,
		arg.Endtime,
		arg.Createdtime,
		arg.Updatedtime,
	)
	return err
}

const updateTodoItem = `-- name: UpdateTodoItem :exec
Update todoItem
SET
    title = $1,
    detail = $2,
    completed = $3,
    startTime = $4,
    endTime = $5,
    createdTime = $6,
    updatedTime = $7
WHERE
    id = $8
`

type UpdateTodoItemParams struct {
	Title       string
	Detail      pgtype.Text
	Completed   bool
	Starttime   pgtype.Timestamp
	Endtime     pgtype.Timestamp
	Createdtime pgtype.Timestamp
	Updatedtime pgtype.Timestamp
	ID          int64
}

func (q *Queries) UpdateTodoItem(ctx context.Context, arg UpdateTodoItemParams) error {
	_, err := q.db.Exec(ctx, updateTodoItem,
		arg.Title,
		arg.Detail,
		arg.Completed,
		arg.Starttime,
		arg.Endtime,
		arg.Createdtime,
		arg.Updatedtime,
		arg.ID,
	)
	return err
}
