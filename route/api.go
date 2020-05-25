package route

import (
	"github.com/gorilla/mux"
)

func GetApiHandlerAuth(r *mux.Router) {
	AddCampaignHandler(r)
	AddClsHandler(r)
	AddAchievementsHandler(r)
	AddAdmissionHandler(r)
	AddCompetitiveGroupsHandler(r)
}

func GetApiHandlerNoOrg(r *mux.Router) {
	AddUserHandler(r)
	AddOrgsHandler(r)
	AddOlympHandler(r)
	AddEntrantHandler(r)
	AddApplicationHandler(r)
	AddDocsHandler(r)
}

func GetApiHandlerNoAuth(r *mux.Router) {
	AddAuthHandler(r)
	AddTestHandler(r)
}
