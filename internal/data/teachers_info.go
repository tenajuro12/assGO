package data

type teachers_info struct {
	ID      int        `json:"id"`
	Name    string     `json:"name"`
	Surname string     `json:"surname"`
	Email   string     `json:"email"`
	Module  ModuleInfo `json:"module"`
}
