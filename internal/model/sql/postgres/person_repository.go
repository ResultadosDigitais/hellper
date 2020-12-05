package postgres

import (
	"fmt"
	"hellper/internal/log"
	"hellper/internal/model"
	"hellper/internal/model/sql"
)

type personRepository struct {
	logger log.Logger
	db     sql.DB
}

// NewPersonRepository creates a new instance of a repository for accessing person data
func NewPersonRepository(logger log.Logger, db sql.DB) model.PersonRepository {
	return &personRepository{logger: logger, db: db}
}

func (r *personRepository) UpdatePersonSlackID(user *model.User) error {
	query := `
		UPDATE public.person SET slack_member_id = $1 WHERE email = $2
	`

	result, err := r.db.Exec(query, user.SlackID, user.Email)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return fmt.Errorf("Person \"%s\" didn't have its slack id updated", user.Email)
	}

	return nil
}
