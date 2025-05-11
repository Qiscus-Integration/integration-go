package api

import (
	"context"
	"errors"
	"fmt"
	"integration-go/internal/auth"
	"integration-go/internal/client"
	"integration-go/internal/config"
	"integration-go/internal/health"
	"integration-go/internal/postgres"
	"integration-go/internal/qismo"
	"integration-go/internal/redis"
	"integration-go/internal/room"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func NewServer() *Server {
	cfg := config.Load()

	db := postgres.NewGORM(cfg.Database)
	postgres.Migrate(db)

	rdb := redis.New(cfg.Redis.URL)

	client := client.New()
	// client.DebugMode = true

	qismo := qismo.New(client, cfg.Qiscus.Omnichannel.URL, cfg.Qiscus.AppID, cfg.Qiscus.SecretKey)

	// Room
	roomRepo := room.NewRepository(db)
	roomSvc := room.NewService(roomRepo, qismo)
	roomHandler := room.NewHttpHandler(roomSvc)

	// Auth
	authMidd := auth.NewMiddleware(cfg.App.SecretKey)

	// Health
	healthRepo := health.NewRepository(db, rdb)
	healthSvc := health.NewService(healthRepo)
	healthHandler := health.NewHttpHandler(healthSvc)

	r := http.NewServeMux()
	r.Handle("GET /", http.HandlerFunc(rootHandler))
	r.Handle("GET /health", http.HandlerFunc(healthHandler.Check))
	r.Handle("POST /wh/qiscus/omnichannel/new-session", http.HandlerFunc(roomHandler.WebhookQismoNewSession))
	r.Handle("GET /api/v1/rooms/{id}", authMidd.StaticToken(http.HandlerFunc(roomHandler.GetRoomByID)))

	return &Server{router: r}
}

type Server struct {
	router *http.ServeMux
}

// Run method of the Server struct runs the HTTP server on the specified port. It initializes
// a new HTTP server instance with the specified port and the server's router.
func (s *Server) Run(port int) {
	addr := fmt.Sprintf(":%d", port)

	h := chainMiddleware(
		s.router,
		recoverHandler,
		loggerHandler(func(w http.ResponseWriter, r *http.Request) bool { return r.URL.Path == "/health" }),
		realIPHandler,
		requestIDHandler,
		corsHandler,
	)

	httpSrv := http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
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
	if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msgf("could not listen on %s", addr)
	}

	<-done
	log.Info().Msg("server stopped")

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
