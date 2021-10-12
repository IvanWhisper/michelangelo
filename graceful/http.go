package graceful

import (
	"context"
	"fmt"
	mlog "github.com/IvanWhisper/michelangelo/infrastructure/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Name   string
	Ip     string
	Port   int
	Engine *gin.Engine
}

func New() *Server {
	return &Server{
		Engine: gin.New(),
	}
}

func (s *Server) Start() {
	// bind ip&port
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Ip, s.Port),
		Handler: s.Engine,
	}

	errCh := make(chan error, 1)
	defer close(errCh)
	// listen
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-ctx.Done():
		stop()
		mlog.Info(fmt.Sprintf("Shutdown: Receive Sign(%s)", s))
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(timeoutCtx); err != nil {
			mlog.Error(fmt.Sprintf("Shutdown: %s", err.Error()))
		}
		mlog.Info("Shutdown: exit")
		break
	case e := <-errCh:
		mlog.Error(fmt.Sprintf("Listen: Receive Error %s", e.Error()))
		break
	}
}
