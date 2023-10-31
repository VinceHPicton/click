// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: delete_user.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}
