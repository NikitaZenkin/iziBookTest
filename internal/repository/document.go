package repository

import (
	"context"

	"iziBookTest/internal/util"
)

func (s *PgStorage) CreateDocument(ctx context.Context, parentSectionId string, document *util.Document) (string, error) {
	userId, _ := util.UserIdFromContext(ctx)
	newId := util.NewID()

	_, err := s.documentsDB.ExecContext(
		ctx,
		`INSERT INTO documents_db.public.documents VALUES ($1, $2, $3, $4, $5)`,
		newId, document.Name, document.Text, userId, parentSectionId,
	)

	return newId, err
}

func (s *PgStorage) GetDocument(ctx context.Context, id string) (*util.Document, error) {
	doc := &util.Document{}

	row := s.documentsDB.QueryRowxContext(ctx, `SELECT * FROM documents_db.public.documents WHERE id = $1`, id)
	err := row.Scan(&doc.ID, &doc.Name, &doc.Text, &doc.OwnerID, &doc.SectionID)

	return doc, err
}

func (s *PgStorage) UpdateDocument(ctx context.Context, id string, documentUpdates *util.Document) error {
	if documentUpdates.Name != "" {
		_, err := s.documentsDB.ExecContext(
			ctx,
			`UPDATE documents_db.public.documents SET name = $2 WHERE id = $1`,
			id, documentUpdates.Name,
		)
		if err != nil {
			return err
		}
	}

	if documentUpdates.Text != "" {
		_, err := s.documentsDB.ExecContext(
			ctx,
			`UPDATE documents_db.public.documents SET text = $2 WHERE id = $1`,
			id, documentUpdates.Text,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PgStorage) DeleteDocument(ctx context.Context, id string) error {
	_, err := s.documentsDB.ExecContext(ctx, `DELETE FROM documents_db.public.documents WHERE id = $1`, id)

	return err
}

func (s *PgStorage) CheckDocument(ctx context.Context, id string) (bool, error) {
	userId, ok := util.UserIdFromContext(ctx)
	if !ok {
		return false, nil
	}

	var permission bool
	err := s.documentsDB.GetContext(
		ctx, &permission,
		`SELECT EXISTS(SELECT * FROM documents_db.public.documents WHERE id = $1 AND owner_id = $2)`,
		id, userId,
	)

	return permission, err
}
