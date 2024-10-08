package models

type Task struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreateDate  string `json:"createdate"`
	Status      string `json:"status"`
}

func (t *Task) Validate() bool {
	return !(t.Title == "" || t.Description == "" || t.CreateDate != "" || t.Status == "")
}

func (t *Task) FullValidate() bool {
	return t.Validate() && t.Id != 0
}
