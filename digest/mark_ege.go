package digest

import "time"

type MarkEgeElement struct {
	Id            uint      `json:"id"`
	IdEntrant     *uint     `json:"id_entrant"`
	Surname       string    `json:"surname"`
	Name          string    `json:"name"`
	Patronymic    *string   `json:"patronymic"`
	IdDocument    *uint     `json:"id_document"`
	DocSeries     string    `json:"doc_series"`
	DocNumber     string    `json:"doc_number"`
	IdDocumentEge *uint     `json:"id_document_ege"`
	IdSubject     *uint     `json:"id_subject"`
	Subject       string    `json:"subject"`
	Mark          int64     `json:"mark"`
	Year          int64     `json:"year"`
	IdRegion      uint      `json:"id_region"`
	Region        Region    `gorm:"foreignkey:id_region"`
	Status        string    `json:"status"`
	AppPer        uint      `json:"app_per"`
	CertNumber    *string   `json:"cert_number"`
	TypNumber     *string   `json:"typ_number"`
	AppStatus     *string   `json:"app_status"`
	IdPackage     uint      `json:"id_package"`
	Checked       bool      `json:"checked"`
	Error         *string   `json:"error"`
	Created       time.Time `json:"created"`
}
type MarkEgePackages struct {
	Id             uint             `json:"id"`
	Status         PackagesStatuses `gorm:"foreignkey:id_status"`
	IdStatus       uint             `json:"id_status" schema:"id_status"`
	IdAuthor       uint             `json:"id_author"`
	IdOrganization uint             `json:"id_organization"`
	Name           string           `json:"name"`
	PathFile       string           `json:"path_file"`
	Error          *string          `json:"error"`
	Created        time.Time        `json:"created"`
	CountAll       int64            `json:"count_all"`
	CountAdd       int64            `json:"count_add"`
}

func (MarkEgeElement) TableName() string {
	return "packages.mark_ege_elements"
}
func (MarkEgePackages) TableName() string {
	return "packages.mark_ege_packages"
}
