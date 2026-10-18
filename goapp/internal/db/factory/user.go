package factory

import (
	"context"
	"database/sql"
	"time"
	"vincehpicton/click/internal/db/sqlc"
)

type UserOption func(*sqlc.CreateUserParams)

func FakeUser(
	ctx context.Context,
	q *sqlc.Queries,
	options ...UserOption,
) (sqlc.AppUser, error) {
	params := sqlc.CreateUserParams{
		Name:      randNullString(10),
		Bio:       randNullString(200),
		BirthDate: sql.NullTime{},
		Mobile:    randString(20),
		Email:     randNullString(20),
		Sex: sql.NullInt16{
			Valid: false,
		},
		InterestedIn: sql.NullInt16{
			Valid: false,
		},
	}

	for _, option := range options {
		option(&params)
	}

	return q.CreateUser(ctx, params)
}

func WithName(name string) UserOption {
	return func(p *sqlc.CreateUserParams) {
		p.Name = sql.NullString{
			String: name,
			Valid:  true,
		}
	}
}

func WithBio(bio string) UserOption {
	return func(p *sqlc.CreateUserParams) {
		p.Bio = sql.NullString{
			String: bio,
			Valid:  true,
		}
	}
}

func WithBirthDate(birthDate time.Time) UserOption {
	return func(p *sqlc.CreateUserParams) {
		p.BirthDate = sql.NullTime{
			Time:  birthDate,
			Valid: true,
		}
	}
}

func WithMobile(mobile string) UserOption {
	return func(p *sqlc.CreateUserParams) {
		p.Mobile = mobile
	}
}

func WithEmail(email string) UserOption {
	return func(p *sqlc.CreateUserParams) {
		p.Email = sql.NullString{
			String: email,
			Valid:  true,
		}
	}
}

func WithSex(sex int16) UserOption {
	return func(p *sqlc.CreateUserParams) {
		p.Sex = sql.NullInt16{
			Int16: sex,
			Valid: true,
		}
	}
}

func WithInterestedIn(interestedIn int16) UserOption {
	return func(p *sqlc.CreateUserParams) {
		p.InterestedIn = sql.NullInt16{
			Int16: interestedIn,
			Valid: true,
		}
	}
}