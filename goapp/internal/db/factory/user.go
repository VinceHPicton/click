package factory

import (
	"context"
	"database/sql"
	"vincehpicton/click/internal/db/sqlc"
)

func FakeUser(ctx context.Context, q *sqlc.Queries) (sqlc.AppUser, error) {
	params := sqlc.CreateUserParams{
		Name:      randNullString(10),
		Bio:       randNullString(200),
		BirthDate: sql.NullTime{},
		Mobile:    randString(20),
		Email:     randNullString(20),
		Sex: sql.NullInt16{
			Int16: 0,
			Valid: false,
		},
		InterestedIn: sql.NullInt16{
			Int16: 0,
			Valid: false,
		},
	}

	return q.CreateUser(ctx, params)
}
