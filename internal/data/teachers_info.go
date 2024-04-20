package data

import "database/sql"

type teachers_info struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	Module  string `json:"module"`
}

type TeacherModel struct {
	TeacherInfo DBModel
}

func NewModel1(db *sql.DB) TeacherModel {
	return TeacherModel{
		TeacherInfo: DBModel{DB: db},
	}
}
