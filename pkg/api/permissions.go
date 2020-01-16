package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/wexel-nath/wexel-auth/pkg/jwt"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
	"github.com/wexel-nath/wexel-auth/pkg/permission"
)

func getServicePermissions(w http.ResponseWriter, r *http.Request) {
	_, err := jwt.Authorize(r, "user.create")
	if err != nil && err != jwt.ErrExpiredToken {
		if err != jwt.ErrExpiredToken {
			logger.Error(err)
		}
		jsonResponse(w, http.StatusUnauthorized, nil, err.Error())
		return
	}

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
