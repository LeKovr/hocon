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
	gen "SELF/zgen/go/proto"
	// importing implementation
	app "SELF/service"
	"SELF/static"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	// "google.golang.org/grpc/credentials/insecure"
)

func main() {
	// create new gRPC server
	grpcSever := grpc.NewServer()
	// register the GreeterServerImpl on the gRPC server
	gen.RegisterHoconServiceServer(grpcSever, &app.HoconServiceImpl{})
	// creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	mux := runtime.NewServeMux(
		// handle incoming headers
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authorization")
			// send all the headers received from the client
			md := metadata.Pairs("auth", header)
			return md
		}),
	)
	// opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	//  err := gw.RegisterYourServiceHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)

	// setting up a dail up for gRPC service by specifying endpoint/target url
	err := gen.RegisterHoconServiceHandlerFromEndpoint(context.Background(), mux, "localhost:8081", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatal(err)
	}

	// static pages server
	hfs, _ := static.New("static") //cfg.Root)
	fileServer := http.FileServer(hfs)

	// Creating a normal HTTP server
	server := http.Server{
		Handler: withLogger(withGW(mux, fileServer)),
	}

	// creating a listener for server
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	m := cmux.New(l)

	// a different listener for HTTP1
	httpL := m.Match(cmux.HTTP1Fast())

	// a different listener for HTTP2 since gRPC uses HTTP2
	grpcL := m.Match(cmux.HTTP2())
	// start server

	// passing dummy listener
	go server.Serve(httpL)
	// passing dummy listener
	go grpcSever.Serve(grpcL)

	fmt.Println("Started")
	// actual listener
	m.Serve()
}

func withLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		log.Printf("http[%d]-- %s -- %s\n", m.Code, m.Duration, request.URL.Path)
	})
}

func withGW(gwmux *runtime.ServeMux, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			gwmux.ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
