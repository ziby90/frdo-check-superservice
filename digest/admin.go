package digest

type Organization struct {
	Id             uint   `xml:"id"`
	Ogrn           string `xml:"ogrn"`
	Kpp            string `xml:"kpp"`
	Inn            string `xml:"inn"`
	ShortTitle     string `xml:"short_title"`
	FullTitle      string `xml:"full_title"`
	IDMunicipality struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"id_municipality"`
	IDRegion string `xml:"id_region"`
	Address  string `xml:"address"`
	Phone    string `xml:"phone"`
	Email    struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"email"`
	IDOrganizationType string `xml:"id_organization_type"`
	ChiefName          struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"chief_name"`
	Created  string `xml:"created"`
	IDAuthor struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"id_author"`
	Actual      bool   `xml:"actual"`
	IDLegalForm string `xml:"id_legal_form"`
	IDEiis      string `xml:"id_eiis"`
	Site        struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"site"`
	MilitaryDepartment string `xml:"military_department"`
	Hostel             string `xml:"hostel"`
	HostelCapacity     string `xml:"hostel_capacity"`
	Filial             string `xml:"filial"`
}

// TableNames
func (Organization) TableName() string {
	return "admin.organizations"
}
