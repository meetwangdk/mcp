package server

import (
	"github.com/gin-gonic/gin"
	apiV1 "hcs-agent/api/v1"
	"hcs-agent/internal/middleware"
	"hcs-agent/internal/router"
	"hcs-agent/pkg/server/http"
)

func NewHTTPServer(
	deps router.RouterDeps,
) *http.Server {
	if deps.Config.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	s := http.NewServer(
		gin.Default(),
		deps.Logger,
		http.WithServerHost(deps.Config.HTTP.Host),
		http.WithServerPort(deps.Config.HTTP.Port),
	)

	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(deps.Logger),
		middleware.RequestLogMiddleware(deps.Logger),
		//middleware.SignMiddleware(log),
	)
	s.GET("/", func(ctx *gin.Context) {
		deps.Logger.WithContext(ctx).Info("hello")
		apiV1.HandleSuccess(ctx, map[string]interface{}{
			":)": "Thank you for using nunu!",
		})
	})

	v1 := s.Group("/v1")
	router.InitUserRouter(deps, v1)

	return s
}
