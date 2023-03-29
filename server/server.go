package server

import (
	"context"
	"fmt"
	"integration-go/common"
	"integration-go/domain"
	"integration-go/repository/api"
	"integration-go/repository/cache"
	"integration-go/repository/persist"
	"integration-go/usecase"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Qiscus-Integration/chilog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

// NewServer function initializes a new instance of the Server struct, which represents
// the HTTP server for the application. The function initializes the necessary dependencies
// for the server to function properly, including the database connection, repositories,
// and use case
func NewServer() *Server {
	dbConn := common.NewDatabase()
	cacheConn := common.NewCache(os.Getenv("REDIS_URL"))

	roomRepo := persist.NewPgsqlRoom(dbConn)
	roomCacheRepo := cache.NewRedisRoom(cacheConn, 10*time.Minute)
	omniRepo := api.NewApiQismo(os.Getenv("QISCUS_APP_ID"), os.Getenv("QISCUS_SECRET_KEY"))

	roomUC := usecase.NewRoom(roomRepo, omniRepo, roomCacheRepo)

	srv := &Server{
		Router: chi.NewRouter(),
		roomUC: roomUC,
	}

	srv.middlewares()
	srv.routes()

	return srv
}

type Server struct {
	Router chi.Router
	roomUC domain.RoomUsecase
}

// Run method of the Server struct runs the HTTP server on the specified port. It initializes
// a new HTTP server instance with the specified port and the server's router.
func (s *Server) Run(port int) {
	addr := fmt.Sprintf(":%d", port)

	httpSrv := http.Server{
		Addr:         addr,
		Handler:      s.Router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info().Msg("server is shuting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpSrv.SetKeepAlivesEnabled(false)
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Msg("could not gracefully shutdown the server")
		}
		close(done)
	}()

	log.Info().Msgf("server serving on port %d", port)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msgf("could not listen on %s", addr)
	}

	<-done
	log.Info().Msg("server stopped")

}

// Middleware functions to the router instance of the server.
func (s *Server) middlewares() {
	s.Router.Use(middleware.RealIP)
	s.Router.Use(chilog.Middleware(func(w http.ResponseWriter, r *http.Request) bool {
		return r.URL.Path == "/"
	}))
}
