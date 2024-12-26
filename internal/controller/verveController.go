package controller

import (
	appcontext "Verve/internal/configs/appContext"
	e "Verve/internal/configs/errorResponse"
	"Verve/internal/model/request"
	"errors"
	"net/http"
)

func GetApi(w http.ResponseWriter, r *http.Request) {
	appCtx := appcontext.GetAppContext()

	request, err := request.SanitizeUrlParams(r)

	if err != nil {
		e.SendError(w, http.StatusBadRequest, errors.New("Failed"))
	}

	err = appCtx.VerveService.SaveAndPost(r.Context(), *request)
	if err != nil {
		e.SendError(w, http.StatusInternalServerError, errors.New("Failed"))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}
