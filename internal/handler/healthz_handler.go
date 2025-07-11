package handler

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	httpx.Ok(w)
}
