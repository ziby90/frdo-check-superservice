package digest

import "time"

type AchievementCategory struct {
	Id       uint `xml:"AchievemntCategoryId"`
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type Benefit struct {
	Id        uint `xml:"BenefitID"`
	Name      string
	Created   time.Time
	IdAuthor  uint // Идентификатор автора
	Actual    bool `json:"actual"` // Актуальность
	ShortName string
}
type CampaignStatus struct {
	Id       uint `xml:"CampaignStatusID"`
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type CampaignType struct {
	Id       uint `xml:"CampaignTypeID"`
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type Direction struct {
	Id           uint `xml:"DirectionID"`
	Section      int64
	Code         string
	Name         string
	CodeMcko2011 string
	CodeMcko2013 string
	Created      time.Time
	IdAuthor     uint // Идентификатор автора
	Actual       bool `json:"actual"` // Актуальность
	IdParent     uint
}
type EducationForm struct {
	Id       uint `xml:"EducationFormID"`
	Code     string
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type EducationLevel struct {
	Id       uint `xml:"EducationLevelID"`
	Code     string
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type EducationSource struct {
	Id       uint `xml:"EducationSourceID"`
	Code     string
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type EntranceTestType struct {
	Id       uint `xml:"EntranceTestTypeID"`
	Code     string
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type LevelBudget struct {
	Id       uint `xml:"LevelBudgetID"`
	Code     string
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type OlympicDiplomType struct {
	Id       uint `xml:"OlympicDiplomTypeID"`
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type OlympicLevel struct {
	Id       uint `xml:"OlympicLevelID"`
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type Subject struct {
	Id       uint `xml:"SubjectID"`
	Code     string
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
	Olympic  bool
}

// TableNames

func (AchievementCategory) TableName() string {
	return "cls.achievement_categories"
}
func (Benefit) TableName() string {
	return "cls.benefits"
}
func (CampaignStatus) TableName() string {
	return "cls.campaign_statuses"
}
func (CampaignType) TableName() string {
	return "cls.campaign_types"
}
func (Direction) TableName() string {
	return "cls.directions"
}
func (EducationForm) TableName() string {
	return "cls.education_forms"
}
func (EducationLevel) TableName() string {
	return "cls.education_levels"
}
func (EducationSource) TableName() string {
	return "cls.education_sources"
}
func (EntranceTestType) TableName() string {
	return "cls.entrance_test_types"
}
func (LevelBudget) TableName() string {
	return "cls.level_budget"
}
func (OlympicDiplomType) TableName() string {
	return "cls.olympic_diploma_types"
}
func (OlympicLevel) TableName() string {
	return "cls.olympic_levels"
}
func (Subject) TableName() string {
	return "cls.subjects"
}
