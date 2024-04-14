package data

import (
	"database/sql"
	"errors"
)

type DBModel struct {
	DB *sql.DB
}
type Models struct {
	ModuleInfo DBModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		ModuleInfo: DBModel{DB: db},
	}
}

func (m DBModel) Insert(module *ModuleInfo) error {

	query := `INSERT INTO module_info (module_name, module_duration, exam_type)VALUES ($1, $2, $3 )RETURNING id, created_at, updated_at, version`
	args := []any{module.ModuleName, module.ModuleDuration, module.ExamType}
	return m.DB.QueryRow(query, args...).Scan(&module.ID, &module.CreatedAt, &module.UpdatedAt, &module.Version)
}

func (m DBModel) Get(id int64) (*ModuleInfo, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the movie data.
	query := `SELECT id, created_at, updated_at, module_name, module_duration, exam_type, version FROM module_info WHERE id = $1`
	var module ModuleInfo

	err := m.DB.QueryRow(query, id).Scan(
		&module.ID,
		&module.CreatedAt,
		&module.UpdatedAt,
		&module.ModuleName,
		&module.ModuleDuration,
		&module.ExamType,
		&module.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &module, nil
}

func (m DBModel) Update(module *ModuleInfo) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `UPDATE module_info SET  module_name= $1, module_duration = $2, exam_type = $3, version = version + 1 WHERE id = $4 RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []any{
		module.ModuleName,
		module.ModuleDuration,
		module.ExamType,
		module.ID,
	}
	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&module.Version)
}

func (m DBModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM module_info WHERE id = $1`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
