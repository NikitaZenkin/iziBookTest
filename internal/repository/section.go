package repository

import (
	"context"
	"github.com/jmoiron/sqlx"

	"iziBookTest/internal/util"
)

func (s *PgStorage) CreateSection(ctx context.Context, parentSectionId, name string) (string, error) {
	userId, _ := util.UserIdFromContext(ctx)
	newId := util.NewID()

	section := &util.Section{
		ID:      newId,
		Name:    name,
		OwnerID: userId,
	}

	if parentSectionId == "" {
		section.Path = section.ID
	} else {
		var parentPath string

		err := s.documentsDB.GetContext(
			ctx, &parentPath,
			`SELECT path FROM documents_db.public.sections WHERE id = $1`,
			parentSectionId,
		)
		if err != nil {
			return "", err
		}

		section.Path = parentPath + "." + section.ID
	}

	_, err := s.documentsDB.ExecContext(
		ctx,
		`INSERT INTO documents_db.public.sections VALUES ($1, $2, $3, $4)`,
		section.ID, section.Path, name, section.OwnerID,
	)

	return section.ID, err
}

func (s *PgStorage) GetSectionTree(ctx context.Context, id string) (*util.SectionTree, error) {
	tx, err := s.documentsDB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	tree, err := getSectionTree(ctx, id, tx)
	if err != nil {
		return nil, err
	}

	return tree, tx.Commit()
}

func (s *PgStorage) UpdateSection(ctx context.Context, id, newName string) error {
	_, err := s.documentsDB.ExecContext(
		ctx,
		`UPDATE documents_db.public.sections SET name = $2 WHERE id = $1`,
		id, newName,
	)

	return err
}

func (s *PgStorage) DeleteSection(ctx context.Context, id string) error {
	_, err := s.documentsDB.ExecContext(
		ctx,
		`DELETE FROM documents_db.public.sections WHERE path <@ (SELECT path FROM sections WHERE id = $1)`,
		id,
	)

	return err
}

func (s *PgStorage) CheckSection(ctx context.Context, id string) (bool, error) {
	userId, ok := util.UserIdFromContext(ctx)
	if !ok {
		return false, nil
	}

	var permission bool
	err := s.documentsDB.GetContext(
		ctx, &permission,
		`SELECT EXISTS(SELECT * FROM documents_db.public.sections WHERE id = $1 AND owner_id = $2)`,
		id, userId,
	)

	return permission, err
}

func getSectionTree(ctx context.Context, id string, tx *sqlx.Tx) (*util.SectionTree, error) {
	section := &util.Section{}
	row := tx.QueryRowxContext(ctx, `SELECT * FROM documents_db.public.sections WHERE id = $1`, id)
	err := row.Scan(&section.ID, &section.Path, &section.Name, &section.OwnerID)
	if err != nil {
		return nil, err
	}

	sectionTree := section.ToSectionTree()

	documents, err := getSectionDocuments(ctx, id, tx)
	if err != nil {
		return nil, err
	}

	sectionTree.Documents = documents

	subSectionIds := make([]string, 0)

	rows, err := tx.QueryxContext(
		ctx,
		`SELECT id FROM documents_db.public.sections WHERE path <@ $1 AND (nlevel(path) - 1 = nlevel($1))`,
		section.Path,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var nextId string
		err = rows.Scan(&nextId)
		if err != nil {
			return nil, err
		}

		subSectionIds = append(subSectionIds, nextId)
	}

	for _, nextId := range subSectionIds {
		nextSubSection, err := getSectionTree(ctx, nextId, tx)
		if err != nil {
			return nil, err
		}

		sectionTree.SubSections = append(sectionTree.SubSections, nextSubSection)
	}

	return sectionTree, nil
}

func getSectionDocuments(ctx context.Context, sectionId string, tx *sqlx.Tx) ([]*util.Document, error) {
	documents := make([]*util.Document, 0)

	rows, err := tx.QueryxContext(ctx, `SELECT * FROM documents_db.public.documents WHERE section_id = $1`, sectionId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		doc := &util.Document{}
		err = rows.Scan(&doc.ID, &doc.Name, &doc.Text, &doc.OwnerID, &doc.SectionID)
		if err != nil {
			return nil, err
		}

		documents = append(documents, doc)
	}

	return documents, nil
}
