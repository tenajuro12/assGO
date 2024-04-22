package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

var (
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
	ModuleInfo  DBModel
	Users       UserModel
	Tokens      TokenModel
	Permissions PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		ModuleInfo:  DBModel{DB: db},
		Users:       UserModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Permissions: PermissionModel{DB: db},
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

	query := `UPDATE module_info SET  module_name= $1, module_duration = $2, exam_type = $3, version = version + 1 WHERE id = $4 RETURNING version`
	args := []any{
		module.ModuleName,
		module.ModuleDuration,
		module.ExamType,
		module.ID,
	}

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

func (m DBModel) GetAll(module_name string, exam_type string, filters Filters) ([]*ModuleInfo, error) {

	query := fmt.Sprintf(`SELECT id, created_at, updated_at, module_name, module_duration, exam_type, version
	FROM module_info
	WHERE (to_tsvector('simple', module_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND  (LOWER(exam_type) = LOWER($2) OR $2 = '')
	ORDER BY  %s %s, id ASC
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{module_name, exam_type, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	module_info := []*ModuleInfo{}

	for rows.Next() {
		var module_infos ModuleInfo

		err := rows.Scan(
			&module_infos.ID,
			&module_infos.CreatedAt,
			&module_infos.UpdatedAt,
			&module_infos.ModuleName,
			&module_infos.ModuleDuration,
			&module_infos.ExamType,
			&module_infos.Version,
		)
		if err != nil {
			return nil, err
		}
		module_info = append(module_info, &module_infos)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return module_info, nil
}

func (m DBModel) GetTeachers(name string, surname string, filters Filters) ([]*teachers_info, error) {

	query := fmt.Sprintf(`SELECT id, name, surname, email, modules
	FROM teachers_info
	WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND  (LOWER(surname) = LOWER($2) OR $2 = '')
	ORDER BY  %s %s, id ASC
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{name, surname, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	teacher_info := []*teachers_info{}

	for rows.Next() {
		var teachers_info teachers_info

		err := rows.Scan(
			&teachers_info.ID,
			&teachers_info.Name,
			&teachers_info.Surname,
			&teachers_info.Email,
			&teachers_info.Module,
		)
		if err != nil {
			return nil, err
		}
		teacher_info = append(teacher_info, &teachers_info)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return teacher_info, nil
}

func (m UserModel) GetAllUser(fname string, sname string, filters Filters) ([]*User, error) {

	query := fmt.Sprintf(`
    SELECT id, created_at, fname, sname, email, activated, version
    FROM user_info
    WHERE (to_tsvector('simple', fname) @@ plainto_tsquery('simple', $1) OR $1 = '')
        AND (LOWER(sname) = LOWER($2) OR $2 = '')
	Order By %s  %s, id ASC
    LIMIT $3 OFFSET $4`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{fname, sname, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user := []*User{}

	for rows.Next() {
		var users User

		err := rows.Scan(
			&users.ID,
			&users.CreatedAt,
			&users.SName,
			&users.FName,
			&users.Email,
			&users.Activated,
			&users.Version,
		)
		if err != nil {
			return nil, err
		}
		user = append(user, &users)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return user, nil
}
