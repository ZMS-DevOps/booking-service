package startup

import (
	"fmt"
	"github.com/ZMS-DevOps/booking-service/infrastructure/persistence/accommodation"
	"github.com/ZMS-DevOps/booking-service/infrastructure/persistence/reservation_request"
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"net"

	//"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ZMS-DevOps/booking-service/application"
	"github.com/ZMS-DevOps/booking-service/domain"
	"github.com/ZMS-DevOps/booking-service/infrastructure/api"
	"github.com/ZMS-DevOps/booking-service/infrastructure/persistence"
	"github.com/ZMS-DevOps/booking-service/startup/config"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
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
	reservationRequestStore := server.initReservationRequestStore(mongoClient)
	unavailabilityService := server.initUnavailabilityService(unavailabilityStore)
	reservationRequestService := server.initReservationRequestService(reservationRequestStore)
	unavailabilityHandler := server.initUnavailabilityHandler(unavailabilityService)
	reservationRequestHandler := server.initReservationRequestHandler(reservationRequestService)
	unavailabilityHandler.Init(server.router)
	reservationRequestHandler.Init(server.router)
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
	store := accommodation.NewUnavailabilityMongoDBStore(client)
	store.DeleteAll()
	for _, unavailability := range unavailabilities {
		err := store.Insert(unavailability)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initReservationRequestStore(client *mongo.Client) domain.ReservationRequestStore {
	store := reservation_request.NewReservationRequestMongoDBStore(client)
	store.DeleteAll()
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

func (server *Server) initReservationRequestService(store domain.ReservationRequestStore) *application.ReservationRequestService {
	return application.NewReservationRequestService(store)
}

func (server *Server) initUnavailabilityHandler(service *application.UnavailabilityService) *api.UnavailabilityHandler {
	return api.NewUnavailabilityHandler(service)
}

func (server *Server) initReservationRequestHandler(service *application.ReservationRequestService) *api.ReservationRequestHandler {
	return api.NewReservationRequestHandler(service)
}

func (server *Server) initGrpcHandler(service *application.UnavailabilityService) *api.BookingHandler {
	return api.NewBookingHandler(service)
}