package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/grpc-ms/user-service/api/proto/user"
	"github.com/grpc-ms/user-service/pkg/database"
	"github.com/grpc-ms/user-service/pkg/http/rpc"
	"github.com/grpc-ms/user-service/pkg/service"
)

const (
	defaultPort = ":50053"
)

func main() {

	Port := os.Getenv("GRPC_PORT")
	if Port == "" {
		Port = defaultPort
	}
	lis, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	db, err := database.CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Database not connected: %v", err)
	} else {
		log.Println("Connected to Postgras DB!")
	}

	// auto migrate user struct to db
	db.AutoMigrate(&pb.User{})

	repo := &service.UserRepository{Database: db}
	tokenService := &service.TokenService{TRepo: repo}

	pb.RegisterUserServiceServer(s, &rpc.Handler{
		Repo:         repo,
		TokenService: tokenService,
	})

	reflection.Register(s)

	go runHttp(fmt.Sprintf("localhost%s", Port))

	log.Println("Running on port:", Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runHttp(clientAddr string) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, clientAddr, opts); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	addr := ":8080"

	s := &http.Server{
		Addr:    addr,
		Handler: allowCORS(mux),
	}

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Failed to listen and serve: %v", err)
	}
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
}
