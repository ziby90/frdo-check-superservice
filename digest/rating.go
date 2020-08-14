package digest

import "time"

type RatingCompetitiveGroupsApplications struct {
	Id                 uint      `gorm:"primary_key"`
	IdOrganization     uint      `json:"id_organization"`
	IdCompetitiveGroup uint      `json:"id_competitive_group"`
	IdApplication      uint      `json:"id_application"`
	AdmissionVolume    int64     `json:"admission_volume"`
	CommonRating       int64     `json:"common_rating"`
	FirstRating        int64     `json:"first_rating"`
	AgreedRating       int64     `json:"agreed_rating"`
	ChangeRating       int64     `json:"change_rating"`
	CountFirstStep     int64     `json:"count_first_step"`
	CountSecondStep    int64     `json:"count_second_step"`
	CountApplication   int64     `json:"count_application"`
	CountAgreed        int64     `json:"count_agreed"`
	Changed            time.Time `json:"changed"`
}

func (RatingCompetitiveGroupsApplications) TableName() string {
	return "rating.completitive_groups_applications"
}
