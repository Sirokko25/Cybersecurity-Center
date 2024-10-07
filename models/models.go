package models

type Task struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreateDate  string `json:"createdate"`
	Status      string `json:"status"`
}

func (t *Task) PostCheckingFields() bool {
	return !(t.Title == "" || t.Description == "" || t.CreateDate != "" || t.Status == "")
}

func (t *Task) PutCheckingFields() bool {
	return t.PostCheckingFields() && t.Id != 0
}
