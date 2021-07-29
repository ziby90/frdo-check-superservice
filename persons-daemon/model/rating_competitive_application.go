package model

import (
	"encoding/xml"
	"time"
)

// пакет рейтингов по конкурсу  база
type RatingCompetitiveApplicationPackages struct {
	Id             uint      `json:"id"`
	IdStatus       uint      `json:"id_status"`
	IdAuthor       *uint     `json:"id_author"`
	IdOrganization uint      `json:"id_organization"`
	Name           string    `json:"name"`
	PathFile       string    `json:"path_file"`
	Error          *string   `json:"error"`
	CreatedAt      time.Time `json:"created"`
	CountAll       int64     `json:"count_all"`
	CountAdd       int64     `json:"count_add"`
	Log
	Duration float64 `json:"duration"`
	TypeFile string  `json:"type_file"`
	Uid      *string `json:"uid"`
}

// элемент пакета рейтингов по конкурсу база
type RatingCompetitiveApplicationElement struct {
	Id uint `json:"id"`
	RatingCompetitiveRequest
	RatingCompetitiveApplication
	IdPackage           uint      `json:"id_package"`
	Checked             bool      `json:"checked"`
	Error               *string   `json:"error"`
	CreatedAt           time.Time `json:"created"`
	IdRatingApplication *uint     `json:"id_rating_application"`
}

// строка в базе элемента рейтингов по конкурсу , который найден
type RatingCompetitiveApplicationRow struct {
	Id uint `json:"id"`
	RatingCompetitiveRequest
	RatingCompetitiveApplication
}

// структура xml рейтингов по конкурсу при загрузке
type RatingCompetitiveApplicationXml struct {
	XMLName                          xml.Name                     `xml:"PackageData"`
	CompetitiveGroupApplicationsList CompetitiveApplicationRating `xml:"CompetitiveGroupApplicationsList"`
}

type CompetitiveApplicationRating struct {
	UIDCompetitiveGroup string    `xml:"UIDCompetitiveGroup"`
	AdmissionVolume     int64     `xml:"AdmissionVolume"`
	CountFirstStep      int64     `xml:"CountFirstStep"`
	CountSecondStep     int64     `xml:"CountSecondStep"`
	UpdatedAt           time.Time `xml:"Changed"`
	Applications        struct {
		Application []struct {
			IDApplicationChoice struct {
				UIDEpgu *int64  `xml:"UIDEpgu"`
				UID     *string `xml:"UID"`
			} `xml:"IDApplicationChoice"`
			Rating             int64    `xml:"Rating"`
			WithoutTests       bool     `xml:"WithoutTests"`
			ReasonWithoutTests *string  `xml:"ReasonWithoutTests"`
			EntranceTest1      string   `xml:"EntranceTest1"`
			Result1            float64  `xml:"Result1"`
			EntranceTest2      string   `xml:"EntranceTest2"`
			Result2            float64  `xml:"Result2"`
			EntranceTest3      string   `xml:"EntranceTest3"`
			Result3            float64  `xml:"Result3"`
			EntranceTest4      *string  `xml:"EntranceTest4"`
			Result4            *float64 `xml:"Result4"`
			EntranceTest5      *string  `xml:"EntranceTest5"`
			Result5            *float64 `xml:"Result5"`
			Benefit            bool     `xml:"Benefit"`
			ReasonBenefit      *string  `xml:"ReasonBenefit"`
			Mark               float64  `xml:"Mark"`
			SumMark            float64  `xml:"SumMark"`
			Agreed             bool     `xml:"Agreed"`
			Original           bool     `xml:"Original"`
			Addition           *string  `xml:"Addition"`
			Enlisted           *uint    `xml:"Enlisted"`
		} `xml:"Application"`
	} `xml:"Applications"`
}
type RatingCompetitiveRequest struct {
	IdCompetitiveGroup *uint     `json:"id_competitive_group" xml:"id_competitive_group"`
	CompetitiveGroup   *string   `json:"competitive_group" xml:"competitive_group"`
	IdOrganization     uint      `json:"id_organization" xml:"id_organization"`
	Organization       string    `json:"organization" xml:"organization"`
	AdmissionVolume    int64     `json:"admission_volume" xml:"admission_volume"`
	CountFirstStep     int64     `json:"count_first_step" xml:"count_first_step"`
	CountSecondStep    int64     `json:"count_second_step" xml:"count_second_step"`
	UpdatedAt          time.Time `json:"changed" xml:"changed"`
}

type RatingCompetitiveApplication struct {
	IdApplication      *uint       `json:"admission_volume" xml:"admission_volume"`
	Orderid            interface{} `json:"orderid" xml:"orderid"`
	Rating             int64       `json:"rating" xml:"rating"`
	WithoutTests       bool        `json:"without_tests" xml:"without_tests"`
	ReasonWithoutTests *string     `json:"reason_without_tests" xml:"reason_without_tests"`
	EntranceTest1      string      `json:"entrance_test1" xml:"entrance_test1"`
	Result1            float64     `json:"result1" xml:"result1"`
	EntranceTest2      string      `json:"entrance_test2" xml:"entrance_test2"`
	Result2            float64     `json:"result2" xml:"result2"`
	EntranceTest3      string      `json:"entrance_test3" xml:"entrance_test3"`
	Result3            float64     `json:"result3" xml:"result3"`
	EntranceTest4      *string     `json:"entrance_test4" xml:"entrance_test4"`
	Result4            *float64    `json:"result4" xml:"result4"`
	EntranceTest5      *string     `json:"entrance_test5" xml:"entrance_test5"`
	Result5            *float64    `json:"result5" xml:"result5"`
	Mark               float64     `json:"mark" xml:"mark"`
	Benefit            bool        `json:"benefit" xml:"benefit"`
	ReasonBenefit      *string     `json:"reason_benefit" xml:"reason_benefit"`
	SumMark            float64     `json:"sum_mark" xml:"sum_mark"`
	Agreed             bool        `json:"agreed" xml:"agreed"`
	Original           bool        `json:"original" xml:"original"`
	Addition           *string     `json:"addition" xml:"addition"`
	Enlisted           *uint       `json:"enlisted" xml:"enlisted"`
	UidEntrant         *string     `json:"uid_entrant"`
}

func (RatingCompetitiveApplicationRow) TableName() string {
	return "rating.completitive_groups_applications"
}
func (RatingCompetitiveApplicationElement) TableName() string {
	return "packages.rating_competitive_applications_element"
}
func (RatingCompetitiveApplicationPackages) TableName() string {
	return "packages.rating_competitive_applications_packages"
}
