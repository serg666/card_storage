package internal

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/serg666/card_storage/pkg"
)

type server struct {
	r           *gin.Engine
	port        string
	cardRepo    CardRepository
	sessionRepo SessionRepository
	keyRepo     KeyRepository
}

func (s *server) Run() error {
	srv := &http.Server{
		Addr:    s.port,
		Handler: s.r,
	}
	s.registerRoutes()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logrus.Infof("Run server on port: %s", s.port)
	go func() {
		<-done
		logrus.Infof("Wait. Server closing...")

		s.stop()

		logrus.Infof("Shutdown status: %v", srv.Shutdown(context.Background()))
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.Errorf("Error on run server: %v", err)
	}

	return nil
}

func (s *server) stop() {
	logrus.Info("Server stopped")
}

func (s *server) registerRoutes() {
	rV1 := s.r.Group("/v1")

	cardCtrl := NewCardController(s.cardRepo, s.keyRepo)
	rV1.GET("/cards", cardCtrl.Index)
	rV1.GET("/cards/token/:token", cardCtrl.GetByToken)
	rV1.POST("/cards", cardCtrl.Create)
	rV1.DELETE("/cards/:id", cardCtrl.Delete)

	sessionCtrl := NewSessionController(s.sessionRepo, s.keyRepo)
	rV1.POST("/sessions", sessionCtrl.Create)
	rV1.GET("/sessions/:key", sessionCtrl.Get)
}

func NewServer(cfg *pkg.AppConfigs, cardRepo CardRepository, sessionRepo SessionRepository, keyRepo KeyRepository) *server {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.LoggerWithFormatter(pkg.GinLogFormatter))
	r.Use(gin.Recovery())

	return &server{
		r:           r,
		port:        cfg.ServerPort,
		cardRepo:    cardRepo,
		sessionRepo: sessionRepo,
		keyRepo:     keyRepo,
	}
}
