
-- name: InsertPerson :one
-- add person into database
INSERT INTO persons (fname, lname, email)
VALUES (
    ?, ?, ?
) RETURNING *;

-- name: ReadPerson :one
SELECT *
FROM persons
WHERE id = ?;

-- name: UpdatePerson :one
UPDATE persons
SET 
    fname = ?,
    lname = ?,
    email = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeletePerson :execresult
DELETE FROM persons WHERE id = ?;