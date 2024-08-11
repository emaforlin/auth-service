package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/emaforlin/auth-service/internal/config"
	"github.com/emaforlin/auth-service/internal/handlers"
	"github.com/emaforlin/auth-service/internal/usecases"
	pb "github.com/emaforlin/auth-service/pkg/pb/protos"
	hclog "github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server interface {
	Start()
}

type rpcServer struct {
	gs  *grpc.Server
	log hclog.Logger
	cfg *config.Config
}

// Start implements Server.
func (r *rpcServer) Start() {
	r.initializeHandlers()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", r.cfg.Port))
	if err != nil {
		r.log.Error(fmt.Sprintf("Unable to listen on %s", l.Addr().String()))
		r.log.Debug(err.Error())
		os.Exit(1)
	}
	r.log.Info(fmt.Sprintf("Listening on %s", l.Addr().String()))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// go routine to serve rpc calls
	go func() {
		if err := r.gs.Serve(l); err != nil {
			r.log.Error("Error serving")
			r.log.Debug(err.Error())

		}
	}()

	// wait for interruption signal
	<-sigChan
	r.log.Info("Shutting down server...")

	r.gs.GracefulStop()
	r.log.Info("Server stopped gracefully")
}

func (r *rpcServer) initializeHandlers() {
	// create usecase
	usecase := usecases.NewAuthUsecase(r.cfg)

	// create handler
	ah := handlers.NewAuthHandler(r.log.Named("grpc"), usecase)

	// setup server reflection
	reflection.Register(r.gs)

	// register the server
	pb.RegisterAuthServer(r.gs, ah)
}

func NewRPCServer(l hclog.Logger, c *config.Config) Server {
	return &rpcServer{
		log: l,
		cfg: c,
		gs:  grpc.NewServer(),
	}
}
