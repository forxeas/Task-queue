package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func checkErr(cmd pgconn.CommandTag, err error) error {
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("job not found")
	}

	return nil
}
