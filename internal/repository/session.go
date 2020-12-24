package repository

import (
	"context"
	"net/http"
	"time"

	"iziBookTest/internal/util"
)

func (s *PgStorage) CreateSession(ctx context.Context, w http.ResponseWriter, userID string) error {
	sessionId := util.NewID()

	_, err := s.usersDB.ExecContext(ctx,
		`DELETE FROM users_db.public.sessions WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return err
	}

	_, err = s.usersDB.ExecContext(
		ctx,
		`INSERT INTO users_db.public.sessions (id, user_id) VALUES ($1, $2)`,
		sessionId,
		userID,
	)
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sessionId,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	return nil
}

func (s *PgStorage) GetSession(r *http.Request) (*util.Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	sessionId := sessionCookie.Value
	session := &util.Session{}
	err = s.usersDB.GetContext(
		r.Context(), session,
		`SELECT * FROM users_db.public.sessions WHERE id = $1`,
		sessionId,
	)

	return session, err
}

func (s *PgStorage) LogOut(ctx context.Context, w http.ResponseWriter) error {
	userId, _ := util.UserIdFromContext(ctx)
	_, err := s.usersDB.ExecContext(ctx, `DELETE FROM users_db.public.sessions WHERE user_id = $1`, userId)
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	return nil
}
