package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	// Import the generated protobuf code
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/ragul28/grpc-ms/consignment-service/proto/consignment"
	userProto "github.com/ragul28/grpc-ms/user-service/proto/user"
	vesselProto "github.com/ragul28/grpc-ms/vessel-service/proto/vessel"
)

const (
	Port          = ":50051"
	DBHost        = "mongodb://localhost:27017"
	VesselAddress = "localhost:50052"
	UserAddress   = "localhost:50053"
)

func main() {

	Port := getEnv("GRPC_PORT", Port)
	DBHost := getEnv("DB_HOST", DBHost)
	VesselAddress := getEnv("VESSEL_HOST", "localhost:50052")

	lis, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		// GRPC server with auth Interceptor
		grpc.UnaryInterceptor(AuthInterceptor),
	)

	//create mongodb client session
	session, err := createClient(DBHost)
	if err != nil {
		log.Panic(err)
	}
	log.Println("DB connected at", DBHost)
	defer session.Disconnect(context.TODO())

	consignmentCollection := session.Database("grpc-ms").Collection("consignments")
	repository := &MongoRepository{consignmentCollection}

	// connection to vessel client via grpc
	conn, err := grpc.Dial(VesselAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	vesselClient := vesselProto.NewVesselServiceClient(conn)

	// Register handler
	pb.RegisterShippingServiceServer(s, &handler{repository, vesselClient})

	reflection.Register(s)

	go runHttp(fmt.Sprintf("localhost%s", Port))

	//Run grpc server
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
	if err := pb.RegisterShippingServiceHandlerFromEndpoint(ctx, mux, clientAddr, opts); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	addr := ":8081"

	s := &http.Server{
		Addr:    addr,
		Handler: allowCORS(mux),
	}

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Failed to listen and serve: %v", err)
	}
}

//Auth middleware to validate token in consignment svc api
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	DisableAuth := getEnv("DISABLE_AUTH", "false")
	UserAddress := getEnv("USER_HOST", "localhost:50053")

	//To skip auth for dev
	if DisableAuth == "true" {
		log.Println("skipping the token auth", DisableAuth)
		return handler(ctx, req)
	}

	// Check incoming context for metadata for jwt token
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
	}

	if len(meta["authorization"]) != 1 {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}

	token := strings.TrimPrefix(meta["authorization"][0], "Bearer ")
	log.Println("Authenticating with token: ", token)
	// authResp, err := TokeValidate(token)

	// Set up a connection to the user server
	conn, err := grpc.Dial(UserAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := userProto.NewUserServiceClient(conn)
	authResp, err := c.ValidateToken(ctx, &userProto.Token{Token: token})
	log.Println("Auth resp:", authResp)
	if err != nil {
		log.Fatalf("could not authenticate: %v", err)
	}

	return handler(ctx, req)
}

// Getenv Helper func
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// Debug function for token validation
func TokeValidate(token string) (bool, error) {
	if len(token) < 1 {
		return false, errors.New("error missing token")
	}
	return token == "secret-token", nil
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
