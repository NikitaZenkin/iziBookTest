package repository

import (
	"context"

	"iziBookTest/internal/util"
)

func (s *PgStorage) CreateUser(ctx context.Context, user *util.User) (string, error) {
	userId := util.NewID()

	_, err := s.usersDB.ExecContext(
		ctx,
		`INSERT INTO users_db.public.users VALUES ($1, $2, $3, $4, $5)`,
		userId, user.Login, util.Md5String(user.PassWord), user.Name, user.DateOfBirth,
	)

	return userId, err
}
func (s *PgStorage) UpdateUser(ctx context.Context, userUpdates *util.User) error {
	userId, _ := util.UserIdFromContext(ctx)

	if userUpdates.PassWord != "" {
		_, err := s.usersDB.ExecContext(
			ctx,
			`UPDATE users_db.public.users SET password = $2 WHERE id = $1`,
			userId, util.Md5String(userUpdates.PassWord),
		)
		if err != nil {
			return err
		}
	}

	if userUpdates.Name != "" {
		_, err := s.usersDB.ExecContext(
			ctx,
			`UPDATE users_db.public.users SET name = $2 WHERE id = $1`,
			userId, userUpdates.Name,
		)
		if err != nil {
			return err
		}
	}

	if userUpdates.DateOfBirth != nil {
		_, err := s.usersDB.ExecContext(
			ctx,
			`UPDATE users_db.public.users SET date_of_birth = $2 WHERE id = $1`,
			userId, userUpdates.DateOfBirth,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PgStorage) DeleteUser(ctx context.Context) error {
	userId, _ := util.UserIdFromContext(ctx)

	_, err := s.documentsDB.ExecContext(ctx, `DELETE FROM documents_db.public.sections WHERE owner_id = $1`, userId)
	if err != nil {
		return err
	}

	_, err = s.usersDB.ExecContext(ctx, `DELETE FROM users_db.public.users WHERE id = $1`, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *PgStorage) LoginExist(ctx context.Context, login string) (bool, error) {
	var loginExist bool
	err := s.usersDB.GetContext(
		ctx, &loginExist,
		`SELECT EXISTS(SELECT * FROM users_db.public.users WHERE login = $1 )`,
		login,
	)

	return loginExist, err
}

func (s *PgStorage) CheckPassWord(ctx context.Context, login, password string) (string, bool, error) {
	info := struct {
		ID       string `db:"id"`
		Password string `db:"password"`
	}{}
	err := s.usersDB.GetContext(
		ctx, &info,
		`SELECT id, password FROM users_db.public.users WHERE login = $1`,
		login,
	)
	if err != nil {
		return "", false, err
	}

	return info.ID, info.Password == util.Md5String(password), nil
}
