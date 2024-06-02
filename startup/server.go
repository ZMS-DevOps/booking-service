package startup

import (
	"fmt"
	"github.com/ZMS-DevOps/booking-service/application"
	"github.com/ZMS-DevOps/booking-service/domain"
	"github.com/ZMS-DevOps/booking-service/infrastructure/api"
	"github.com/ZMS-DevOps/booking-service/infrastructure/persistence"
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/ZMS-DevOps/booking-service/startup/config"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

type Server struct {
	config *config.Config
	router *mux.Router
}

func NewServer(config *config.Config) *Server {
	server := &Server{
		config: config,
		router: mux.NewRouter(),
	}
	return server
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	unavailabilityStore := server.initUnavailabilityStore(mongoClient)
	unavailabilityService := server.initUnavailabilityService(unavailabilityStore)
	unavailabilityHandler := server.initUnavailabilityHandler(unavailabilityService)
	unavailabilityHandler.Init(server.router)
	grpcHandler := server.initGrpcHandler(unavailabilityService)
	go server.startGrpcServer(grpcHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.router))
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.BookingDBUsername, server.config.BookingDBPassword, server.config.BookingDBHost, server.config.BookingDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initUnavailabilityStore(client *mongo.Client) domain.UnavailabilityStore {
	store := persistence.NewUnavailabilityMongoDBStore(client)
	store.DeleteAll()
	for _, unavailability := range unavailabilities {
		err := store.Insert(unavailability)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) startGrpcServer(bookingHandler *api.BookingHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	booking.RegisterBookingServiceServer(grpcServer, bookingHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initUnavailabilityService(store domain.UnavailabilityStore) *application.UnavailabilityService {
	return application.NewUnavailabilityService(store)
}

func (server *Server) initUnavailabilityHandler(service *application.UnavailabilityService) *api.UnavailabilityHandler {
	return api.NewUnavailabilityHandler(service)
}

func (server *Server) initGrpcHandler(service *application.UnavailabilityService) *api.BookingHandler {
	return api.NewBookingHandler(service)
}
