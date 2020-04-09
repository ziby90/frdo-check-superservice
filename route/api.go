package route

import (
	"github.com/gorilla/mux"
)

func GetApiHandlerAuth(r *mux.Router) {
	AddCampaignHandler(r)
	AddClsHandler(r)
	AddAchievementsHandler(r)
}

func GetApiHandlerNoOrg(r *mux.Router) {
	AddUserHandler(r)
	AddOrgsHandler(r)
	AddEntrantHandler(r)
	AddDocsHandler(r)
}

func GetApiHandlerNoAuth(r *mux.Router) {
	AddAuthHandler(r)
}
