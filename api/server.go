package api

import (
	"context"
	"fmt"
	"integration-go/client"
	"integration-go/config"
	"integration-go/entity"
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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewServer() *Server {
	cfg := config.Load()

	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Name,
		cfg.Database.Password,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal().Err(err).Msg("unable to open db connection")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get sql db")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.AutoMigrate(&entity.Room{})
	if err != nil {
		log.Fatal().Msgf("unable to migrate database: %s", err.Error())
	}

	client := client.New()
	qismo := qismo.New(client, cfg.Qiscus.Omnichannel.URL, cfg.Qiscus.AppID, cfg.Qiscus.SecretKey)

	// Room
	roomRepo := room.NewRepository(db)
	roomSvc := room.NewService(roomRepo, qismo)
	roomHandler := room.NewHttpHandler(roomSvc)

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
