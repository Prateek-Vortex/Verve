package controller

import (
	appcontext "Verve/internal/configs/appContext"
	e "Verve/internal/configs/errorResponse"
	"Verve/internal/model/request"
	"net/http"
)

func GetApi(w http.ResponseWriter, r *http.Request) {
	appCtx := appcontext.GetAppContext()

	request, err := request.SanitizeUrlParams(r)

	if err != nil {
		e.SendResponse(w, http.StatusBadRequest, "failed")
		return
	}

	err = appCtx.VerveService.SaveAndPost(r.Context(), *request)
	if err != nil {
		e.SendResponse(w, http.StatusInternalServerError, "failed")
		return
	}

	e.SendResponse(w, http.StatusOK, "ok")
}
