// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: user_get.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const userGet = `-- name: UserGet :one
SELECT id, name, bio, birth_date, last_location_long, last_location_lat, mobile, last_active, email, sex, interested_in, created_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) UserGet(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, userGet, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Bio,
		&i.BirthDate,
		&i.LastLocationLong,
		&i.LastLocationLat,
		&i.Mobile,
		&i.LastActive,
		&i.Email,
		&i.Sex,
		&i.InterestedIn,
		&i.CreatedAt,
	)
	return i, err
}