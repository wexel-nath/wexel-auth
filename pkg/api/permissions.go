package api

import (
	"net/http"

	"wexel-auth/pkg/jwt"
	"wexel-auth/pkg/logger"
	"wexel-auth/pkg/permission"

	"github.com/go-chi/chi"
)

func getServicePermissions(w http.ResponseWriter, r *http.Request) {
	serviceName := chi.URLParam(r, "serviceName")

	permissions, err := permission.GetAllForService(serviceName)
	if err != nil {
		if err != jwt.ErrExpiredToken {
			logger.Error(err)
		}
		jsonResponse(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, permissions, "")
}
