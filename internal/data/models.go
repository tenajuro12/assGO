package data

import (
	"database/sql"
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

// Retrieve retrieves data from the specified table in the database.
func (m *DBModel) Retrieve(moduleInfo *ModuleInfo) (interface{}, error) {
	// Implement the retrieve logic here
	return nil, nil
}

// Update updates data in the specified table in the database.
func (m *DBModel) Update(moduleInfo *ModuleInfo) error {
	// Implement the update logic here
	return nil
}

// Delete deletes data from the specified table in the database.
func (m *DBModel) Delete(moduleInfo *ModuleInfo) error {
	// Implement the delete logic here
	return nil
}

func (m DBModel) Insert(module *ModuleInfo) error {

	query := `INSERT INTO module_info (module_name, module_duration, exam_type,version)VALUES ($1, $2, $3, $4)RETURNING id, created_at, updated_at, version`
	args := []any{module.ModuleName, module.ModuleDuration, module.ExamType}
	return m.DB.QueryRow(query, args...).Scan(&module.ID, &module.CreatedAt, &module.UpdatedAt, &module.Version)
}
