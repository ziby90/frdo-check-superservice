package digest

import "time"

type News struct {
	Id        uint      `json:"id" schema:"id"`
	Title     *string   `json:"title" schema:"title"`
	Content   *string   `json:"content" schema:"content"`
	DateNews  time.Time `json:"date_news" schema:"date_news"`
	Created   time.Time `json:"created" schema:"created"`
	Published bool      `json:"published" schema:"published"`
	Deleted   bool      `json:"deleted" schema:"deleted"`
	IdAuthor  uint      `json:"id_author" schema:"id_author"`
}

type FileNew struct {
	Id       uint      `json:"id"`
	IdNews   uint      `json:"id_news"`
	Title    string    `json:"title"`
	Key      string    `json:"key"`
	Size     int64     `json:"size"`
	Mime     string    `json:"-"`
	PathFile string    `json:"path_file"`
	IdAuthor uint      `json:"id_author" schema:"id_author"`
	Created  time.Time `json:"created"`
}

func (News) TableName() string {
	return `info.news`
}
func (FileNew) TableName() string {
	return `info.files`
}
