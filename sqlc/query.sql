-- name: GetTodoItemById :one
SELECT * FROM todoItem
WHERE id = $1 LIMIT 1;

-- name: GetTodoItemByTitle :one
SELECT * FROM todoItem
WHERE title = $1 LIMIT 1;

-- name: GetAllTodoItem :many
SELECT * FROM todoItem ORDER BY title;

-- name: InsertTodoItem :exec
INSERT INTO todoItem (title, detail, completed, startTime, endTime, createdTime, updatedTime) 
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateTodoItem :exec
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
    id = $8;

-- name: DeleteTodoItem :exec
DELETE FROM todoItem WHERE id = $1;