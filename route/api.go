package route

import (
	"github.com/gorilla/mux"
)

func GetApiHandler(r *mux.Router) {
	AddCampaignHandler(r)
	AddClsHandler(r)
	AddAchievementsHandler(r)
}
