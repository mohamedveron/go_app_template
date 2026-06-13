package persistence

import (
	"context"
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohamedveron/go_app_template/internal/users/domain"
)

const defaultListLimit = 20

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
	var firstName, lastName, mobile, storeEmail pgtype.Text

	row := us.pqdriver.QueryRow(ctx, query, args...)
	err = row.Scan(
		&firstName,
		&lastName,
		&mobile,
		&storeEmail,
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

func (us *UserPostgresPersistence) List(ctx context.Context, cursor domain.Cursor) (*domain.UserPage, error) {
	limit := cursor.Limit
	if limit == 0 {
		limit = defaultListLimit
	}

	q := us.qbuilder.Select(
		"firstName",
		"lastName",
		"mobile",
		"email",
		"createdAt",
		"updatedAt",
	).From(us.tableName).
		OrderBy("createdAt ASC").
		// fetch one extra row to determine whether a next page exists
		Limit(uint64(limit) + 1)

	if cursor.After != nil {
		q = q.Where(squirrel.Gt{"createdAt": cursor.After})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.New("internal error")
	}

	rows, err := us.pqdriver.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.New("internal error")
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := new(domain.User)
		var firstName, lastName, mobile, email pgtype.Text
		if err := rows.Scan(&firstName, &lastName, &mobile, &email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, errors.New("internal error")
		}
		u.FirstName = firstName.String
		u.LastName = lastName.String
		u.Mobile = mobile.String
		u.Email = email.String
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("internal error")
	}

	page := &domain.UserPage{Users: users}

	if uint(len(users)) > limit {
		// trim the extra sentinel row and expose its createdAt as the next cursor
		page.Users = users[:limit]
		page.NextCursor = users[limit].CreatedAt
	}

	return page, nil
}

func NewUserPostgresPersistence(pqdriver *pgxpool.Pool) (*UserPostgresPersistence, error) {
	return &UserPostgresPersistence{
		pqdriver:  pqdriver,
		qbuilder:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		tableName: "Users",
	}, nil
}
