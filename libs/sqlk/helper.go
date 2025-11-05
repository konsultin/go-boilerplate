package sqlk

import (
	"database/sql"
)

func IsUpdated(result sql.Result) error {
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return RowNotUpdatedError
	}

	return nil
}
