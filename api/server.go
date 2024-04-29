package api

import (
	"context"
	"fmt"
	"integration-go/client"
	"integration-go/config"
	"integration-go/entity"
	"integration-go/pgsql"
	"integration-go/qismo"
	"integration-go/room"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Qiscus-Integration/chilog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"

	roomHttpHandler "integration-go/room/handler"
)

// NewServer function initializes a new instance of the Server struct, which represents
// the HTTP server for the application. The function initializes the necessary dependencies
// for the server to function properly, including the database connection, repositories,
// and use case
func NewServer() *Server {
	cfg := config.Load()
	db := pgsql.NewDatabase(cfg)

	err := db.AutoMigrate(&entity.Room{})
	if err != nil {
		log.Fatal().Msgf("unable to migrate database: %s", err.Error())
	}

	// Adapter Packages
	client := client.New()
	roomRepo := pgsql.NewRoom(db)
	qismo := qismo.New(client, cfg.Qiscus.Omnichannel.URL, cfg.Qiscus.AppID, cfg.Qiscus.SecretKey)

	// Services
	roomSvc := room.NewService(roomRepo, qismo)

	// Handlers
	roomHandler := roomHttpHandler.NewHttpHandler(roomSvc)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(chilog.Middleware(func(w http.ResponseWriter, r *http.Request) bool {
		return r.URL.Path == "/"
	}))

	r.Use(cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"POST", "GET", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		MaxAge:             60,
		AllowCredentials:   true,
		OptionsPassthrough: false,
		Debug:              false,
	}).Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/wh", func(r chi.Router) {
		r.Post("/qiscus/omnichannel/new-session", roomHandler.WebhookQismoNewSession)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.With(staticTokenAuthMiddleware(cfg.App.SecretKey)).Group(func(r chi.Router) {
			r.Get("/rooms/:id", roomHandler.GetRoomByID)
		})
	})

	return &Server{router: r}
}

type Server struct {
	router chi.Router
}

// Run method of the Server struct runs the HTTP server on the specified port. It initializes
// a new HTTP server instance with the specified port and the server's router.
func (s *Server) Run(port int) {
	addr := fmt.Sprintf(":%d", port)

	httpSrv := http.Server{
		Addr:         addr,
		Handler:      s.router,
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
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msgf("could not listen on %s", addr)
	}

	<-done
	log.Info().Msg("server stopped")

}
