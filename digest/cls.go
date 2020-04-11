package digest

import "time"

type AchievementCategory struct {
	Id       uint `xml:"AchievemntCategoryId"`
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type AppealStatuses struct {
	Id       uint `xml:"AppealStatusId"`
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
type CompositionThemes struct {
	Id          uint `xml:"CompositionThemeID"`
	Name        string
	YearTheme   int64
	NumberTheme int64
	Created     time.Time
	IdAuthor    uint // Идентификатор автора
	Actual      bool `json:"actual"` // Актуальность
}
type CompatriotCategory struct {
	Id      uint `xml:"CompositionThemeID"`
	Name    string
	Created time.Time
	Actual  bool `json:"actual"` // Актуальность
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
type DisabilityTypes struct {
	Id      uint `xml:"DisabilityTypeID"`
	Name    string
	Created time.Time
	Actual  bool `json:"actual"` // Актуальность
}
type DocumentCategories struct {
	Id       uint `xml:"DocumentCategorieID"`
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type DocumentSysCategories struct {
	Id       uint `xml:"DocumentSysCategorieID"`
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type DocumentTypes struct {
	Id                    uint `xml:"DocumentTypeID"`
	Name                  string
	DocumentCategories    DocumentCategories    `json:"document_categorie" gorm:"foreignkey:id_category"`
	IdCategory            uint                  `json:"id_category"`
	DocumentSysCategories DocumentSysCategories `json:"document_sys_categorie" gorm:"foreignkey:id_sys_category"`
	IdSysCategory         uint                  `json:"id_sys_category"`
	Created               time.Time
	Actual                bool `json:"actual"` // Актуальность
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
type Gender struct {
	Id      uint `xml:"GenderID"`
	Name    string
	Created time.Time
	Actual  bool `json:"actual"` // Актуальность
}
type LevelBudget struct {
	Id       uint `xml:"LevelBudgetID"`
	Code     string
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type MilitaryCategories struct {
	Id      uint `xml:"MilitaryCategoryID"`
	Name    string
	Created time.Time
	Actual  bool `json:"actual"` // Актуальность
}
type Okcm struct {
	Id        uint `xml:"OkcmId"`
	ShortName string
	FullName  string
	Af        string
	Afg       string
	Created   time.Time
	Actual    bool `json:"actual"` // Актуальность
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
type OrphanCategories struct {
	Id       uint `xml:"OrphanCategoryID"`
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type ParentsLostCategories struct {
	Id       uint `xml:"ParentsLostCategoryID"`
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type RadiationWorkCategories struct {
	Id       uint `xml:"RadiationWorkCategoryID"`
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type VeteranCategories struct {
	Id       uint `xml:"VeteranCategoryID"`
	Name     string
	Created  time.Time
	IdAuthor uint // Идентификатор автора
	Actual   bool `json:"actual"` // Актуальность
}
type Region struct {
	Id      uint `xml:"RegionID"`
	Code    string
	Name    string
	Created time.Time
	Actual  bool `json:"actual"` // Актуальность
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
func (AppealStatuses) TableName() string {
	return "cls.appeal_statuses"
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
func (CompositionThemes) TableName() string {
	return "cls.composition_themes"
}
func (CompatriotCategory) TableName() string {
	return "cls.compatriot_categories"
}
func (Direction) TableName() string {
	return "cls.directions"
}
func (DisabilityTypes) TableName() string {
	return "cls.disability_types"
}
func (DocumentCategories) TableName() string {
	return "cls.document_categories"
}
func (DocumentSysCategories) TableName() string {
	return "cls.document_sys_categories"
}
func (DocumentTypes) TableName() string {
	return "cls.document_types"
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
func (Gender) TableName() string {
	return "cls.gender"
}
func (LevelBudget) TableName() string {
	return "cls.level_budget"
}
func (MilitaryCategories) TableName() string {
	return "cls.military_categories"
}
func (OlympicDiplomType) TableName() string {
	return "cls.olympic_diploma_types"
}
func (OlympicLevel) TableName() string {
	return "cls.olympic_levels"
}
func (Okcm) TableName() string {
	return "cls.okcm"
}
func (OrphanCategories) TableName() string {
	return "cls.orphan_categories"
}
func (ParentsLostCategories) TableName() string {
	return "cls.parents_lost_categories"
}
func (RadiationWorkCategories) TableName() string {
	return "cls.radiation_work_categories"
}
func (Region) TableName() string {
	return "cls.regions"
}
func (Subject) TableName() string {
	return "cls.subjects"
}
func (VeteranCategories) TableName() string {
	return "cls.veteran_categories"
}
