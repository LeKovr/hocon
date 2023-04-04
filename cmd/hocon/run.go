package main

// See also: https://blog.logrocket.com/guide-to-grpc-gateway/

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"

	// importing generated stubs
	gen "github.com/LeKovr/hocon/zgen/go/proto"
	// importing implementation
	service "github.com/LeKovr/hocon"
	"github.com/LeKovr/hocon/static"

	"github.com/LeKovr/go-kit/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	// "google.golang.org/grpc/credentials/insecure"
)

// Config holds all config vars
type Config struct {
	// logger.Config
	ListenHost string `long:"listen" default:"0.0.0.0" description:"Addr which server listens at"`
	ListenPort int    `long:"grpc" default:"8080" description:"Port which server listens at"`
	Root       string `long:"root" env:"ROOT" default:""  description:"Static files root directory"`
	Token      string `long:"token" env:"TOKEN" description:"Authorization token"`

	//GracePeriod    time.Duration `long:"grace" default:"1m" description:"Stop grace period"`

	Service service.Config `group:"Service Options" namespace:"srv" env-namespace:"SRV"`
}

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

// Run app and exit via given exitFunc
func Run(version string, exitFunc func(code int)) {
	// Load config
	var cfg Config
	err := config.Open(&cfg)
	defer func() { config.Close(err, exitFunc) }()
	if err != nil {
		return
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureValidToken(cfg.Token)),
	}

	// create new gRPC server
	grpcSever := grpc.NewServer(opts...)
	// register the GreeterServerImpl on the gRPC server
	gen.RegisterHoconServiceServer(grpcSever, service.New(cfg.Service))
	// creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	mux := runtime.NewServeMux(
		// handle incoming headers
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authorization")
			// send all the headers received from the client
			md := metadata.Pairs("Authorization", header)
			return md
		}),
	)
	// opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	//  err := gw.RegisterYourServiceHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)

	listenAddr := fmt.Sprintf("%s:%d", cfg.ListenHost, cfg.ListenPort)
	clientAddr := chooseClientAddr(cfg.ListenHost, cfg.ListenPort)
	// setting up a dial up for gRPC service by specifying endpoint/target url
	err = gen.RegisterHoconServiceHandlerFromEndpoint(context.Background(), mux, clientAddr, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return
	}

	// static pages server
	hfs, _ := static.New("static") //cfg.Root)
	fileServer := http.FileServer(hfs)

	// Creating a normal HTTP server
	server := http.Server{
		Handler: withLogger(withGW(mux, fileServer)),
	}

	// creating a listener for server
	var listener net.Listener
	listener, err = net.Listen("tcp", listenAddr)
	if err != nil {
		return
	}
	m := cmux.New(listener)

	// a different listener for HTTP1
	httpL := m.Match(cmux.HTTP1Fast())

	// a different listener for HTTP2 since gRPC uses HTTP2
	grpcL := m.Match(cmux.HTTP2())

	// start server

	// passing dummy listener
	go server.Serve(httpL)
	// passing dummy listener
	go grpcSever.Serve(grpcL)

	fmt.Println("Started at ", listenAddr)
	// actual listener
	m.Serve()
}

// withLogger prints HTTP request log
func withLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		log.Printf("http[%d]-- %s -- %s\n", m.Code, m.Duration, request.RequestURI)
	})
}

// withGW routes /api requests to grpc gateway
func withGW(gwmux *runtime.ServeMux, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			gwmux.ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

// chooseClientAddr chooses localhost if server listens any ip
func chooseClientAddr(host string, port int) string {
	if host == "0.0.0.0" {
		host = "localhost"
	}
	return fmt.Sprintf("%s:%d", host, port)
}

// valid validates the authorization.
func valid(token string, authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	got := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == got
}

// ensureValidToken ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func ensureValidToken(token string) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errMissingMetadata
		}
		if !valid(token, md["authorization"]) {
			return nil, errInvalidToken
		}
		// Continue execution of handler after ensuring a valid token.
		return handler(ctx, req)
	}
}
