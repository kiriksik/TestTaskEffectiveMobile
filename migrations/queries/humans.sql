-- name: CreateHuman :one
INSERT INTO humans (id, name, surname, patronymic, age, gender, country, created_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetHumanByID :one
SELECT * FROM humans 
WHERE id = $1;


-- name: GetHumans :many
SELECT * FROM humans;

-- name: UpdateHuman :one
UPDATE humans
SET name = $2, surname = $3, patronymic = $4, age = $5, gender = $6, country = $7
WHERE id = $1
RETURNING *;


-- name: DeleteHuman :one
DELETE FROM humans
WHERE id = $1
RETURNING *;