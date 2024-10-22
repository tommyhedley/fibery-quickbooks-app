package webhooks

import (
	"net/http"

	"github.com/tommyhedley/fibery/fibery-tsheets-integration/internal/utils"
)

func TransformHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, struct{}{})
}
