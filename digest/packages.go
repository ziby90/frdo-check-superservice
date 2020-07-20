package digest

import (
	"time"
)

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

type RatingApplicationsPackages struct {
	Id             uint      `json:"id"`
	IdStatus       uint      `json:"id_status"`
	IdAuthor       uint      `json:"id_author"`
	IdOrganization uint      `json:"id_organization"`
	Name           string    `json:"name"`
	PathFile       string    `json:"path_file"`
	Error          *string   `json:"error"`
	Created        time.Time `json:"created"`
	CountAll       int64     `json:"count_all"`
	CountAdd       int64     `json:"count_add"`
}

type RatingApplicationsElement struct {
	Id uint `json:"id"`
	RatingRequest
	RatingApplication
	IdPackage           uint      `json:"id_package"`
	Checked             bool      `json:"checked"`
	Error               *string   `json:"error"`
	Created             time.Time `json:"created"`
	IdRatingApplication *uint     `json:"id_rating_application"`
}

type RatingRequest struct {
	IdCompetitiveGroup uint   `json:"id_competitive_group" xml:"id_competitive_group"`
	CompetitiveGroup   string `json:"competitive_group" xml:"competitive_group"`
	IdOrganization     uint   `json:"id_organization" xml:"id_organization"`
	Organization       string `json:"organization" xml:"organization"`
	AdmissionVolume    int64  `json:"admission_volume" xml:"admission_volume"`
	CountFirstStep     int64  `json:"count_first_step" xml:"count_first_step"`
	CountSecondStep    int64  `json:"count_second_step" xml:"count_second_step"`
	Changed            string `json:"changed" xml:"changed"`
}

type RatingApplication struct {
	IdApplication      uint     `json:"admission_volume" xml:"admission_volume"`
	Orderid            *int64   `json:"orderid" xml:"orderid"`
	Fio                string   `json:"fio" xml:"fio"`
	Rating             int64    `json:"rating" xml:"rating"`
	WithoutTests       bool     `json:"without_tests" xml:"without_tests"`
	ReasonWithoutTests *string  `json:"reason_without_tests" xml:"reason_without_tests"`
	EntranceTest1      string   `json:"entrance_test1" xml:"entrance_test1"`
	Result1            float64  `json:"result1" xml:"result1"`
	EntranceTest2      string   `json:"entrance_test2" xml:"entrance_test2"`
	Result2            float64  `json:"result2" xml:"result2"`
	EntranceTest3      string   `json:"entrance_test3" xml:"entrance_test3"`
	Result3            float64  `json:"result3" xml:"result3"`
	EntranceTest4      *string  `json:"entrance_test4" xml:"entrance_test4"`
	Result4            *float64 `json:"result4" xml:"result4"`
	EntranceTest5      *string  `json:"entrance_test5" xml:"entrance_test5"`
	Result5            *float64 `json:"result5" xml:"result5"`
	Mark               float64  `json:"mark" xml:"mark"`
	Benefit            bool     `json:"benefit" xml:"benefit"`
	ReasonBenefit      *string  `json:"reason_benefit" xml:"reason_benefit"`
	SumMark            float64  `json:"sum_mark" xml:"sum_mark"`
	Agreed             bool     `json:"agreed" xml:"agreed"`
	Original           bool     `json:"original" xml:"original"`
	Addition           *string  `json:"addition" xml:"addition"`
	Enlisted           *uint    `json:"enlisted" xml:"enlisted"`
}

func (RatingApplicationsElement) TableName() string {
	return "packages.rating_applications_elements"
}
func (RatingApplicationsPackages) TableName() string {
	return "packages.rating_applications_packages"
}
func (MarkEgeElement) TableName() string {
	return "packages.mark_ege_elements"
}
func (MarkEgePackages) TableName() string {
	return "packages.mark_ege_packages"
}
