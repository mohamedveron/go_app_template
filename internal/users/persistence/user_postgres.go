package persistence

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohamedveron/go_app_template/internal/users/domain"
)

type UserPostgresPersistence struct {
	qbuilder  squirrel.StatementBuilderType
	pqdriver  *pgxpool.Pool
	tableName string
}

func (us *UserPostgresPersistence) Create(ctx context.Context, u *domain.User) error {
	query, args, err := us.qbuilder.Insert(us.tableName).SetMap(map[string]interface{}{
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"mobile":    u.Mobile,
		"email":     u.Email,
		"createdAt": u.CreatedAt,
		"updatedAt": u.UpdatedAt,
	}).ToSql()
	if err != nil {
		return errors.New("internal error")
	}

	_, err = us.pqdriver.Exec(ctx, query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			return errors.New("user with email '%s' already exists")
		}
		return errors.New("internal error")
	}

	return nil
}

func (us *UserPostgresPersistence) ReadByEmail(ctx context.Context, email string) (*domain.User, error) {
	query, args, err := us.qbuilder.Select(
		"firstName",
		"lastName",
		"mobile",
		"email",
		"createdAt",
		"updatedAt",
	).From(
		us.tableName,
	).Where(
		squirrel.Eq{
			"email": email,
		},
	).ToSql()
	if err != nil {
		return nil, errors.New("internal error")
	}

	user := new(domain.User)
	firstName := new(sql.NullString)
	lastName := new(sql.NullString)
	mobile := new(sql.NullString)
	storeEmail := new(sql.NullString)

	row := us.pqdriver.QueryRow(ctx, query, args...)
	err = row.Scan(
		firstName,
		lastName,
		mobile,
		storeEmail,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("email not found")
		}

		return nil, errors.New("internal error")
	}

	user.FirstName = firstName.String
	user.LastName = lastName.String
	user.Mobile = mobile.String
	user.Email = storeEmail.String

	return user, nil
}

func NewUserPostgresPersistence(pqdriver *pgxpool.Pool) (*UserPostgresPersistence, error) {
	return &UserPostgresPersistence{
		pqdriver:  pqdriver,
		qbuilder:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		tableName: "Users",
	}, nil
}
