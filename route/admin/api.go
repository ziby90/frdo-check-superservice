package route_admin

import (
	"github.com/gorilla/mux"
)

func GetApiAdminHandlerAuth(r *mux.Router) {
	AddCampaignHandler(r)
	//AddClsHandler(r)
	//AddAchievementsHandler(r)
	//AddAdmissionHandler(r)
	AddCompetitiveGroupsHandler(r)
	AddUserHandler(r)
	//AddCertificatesHandler(r)
	//AddChecksHandler(r)
	//AddApplicationHandler(r)
	//AddEntityHandler(r)

}

func GetApiHandlerNoOrg(r *mux.Router) {
	//AddUserHandler(r)
	//AddOrgsHandler(r)
	//AddOlympHandler(r)
	//AddEntrantHandler(r)
	//AddDocsHandler(r)

}

func GetApiHandlerNoAuth(r *mux.Router) {
	//AddAuthHandler(r)
	//AddNewHandler(r)
	//AddTestHandler(r)
	//AddAnalyticsHandler(r)
}
