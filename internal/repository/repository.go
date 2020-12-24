package repository

import (
	"context"
	"net/http"

	"iziBookTest/internal/util"
)

type Repository interface {
	UserRepository
	SectionRepository
	DocumentRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *util.User) (string, error)
	UpdateUser(ctx context.Context, user *util.User) error
	DeleteUser(ctx context.Context) error
	LoginExist(ctx context.Context, login string) (bool, error)
	CheckPassWord(ctx context.Context, login, password string) (string, bool, error)
}

type SectionRepository interface {
	CreateSection(ctx context.Context, parentSectionId, Name string) (string, error)
	GetSectionTree(ctx context.Context, id string) (*util.SectionTree, error)
	UpdateSection(ctx context.Context, id, newName string) error
	DeleteSection(ctx context.Context, id string) error
	CheckSection(ctx context.Context, id string) (bool, error)
}

type DocumentRepository interface {
	CreateDocument(ctx context.Context, parentSectionId string, document *util.Document) (string, error)
	GetDocument(ctx context.Context, id string) (*util.Document, error)
	UpdateDocument(ctx context.Context, id string, document *util.Document) error
	DeleteDocument(ctx context.Context, id string) error
	CheckDocument(ctx context.Context, id string) (bool, error)
}

type SessionManager interface {
	CreateSession(ctx context.Context, w http.ResponseWriter, userID string) error
	GetSession(r *http.Request) (*util.Session, error)
	LogOut(ctx context.Context, w http.ResponseWriter) error
}
