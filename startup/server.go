package startup

import (
	"fmt"
	"github.com/ZMS-DevOps/booking-service/infrastructure/persistence/reservation_request"
	"github.com/ZMS-DevOps/booking-service/infrastructure/persistence/unavailability"
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/mux"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"net"

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
	config        *config.Config
	router        *mux.Router
	traceProvider *sdktrace.TracerProvider
	loki          promtail.Client
}

func NewServer(config *config.Config, traceProvider *sdktrace.TracerProvider, loki promtail.Client) *Server {
	server := &Server{
		config:        config,
		router:        mux.NewRouter(),
		traceProvider: traceProvider,
		loki:          loki,
	}
	return server
}

func (server *Server) Start(producer *kafka.Producer) {
	mongoClient := server.initMongoClient()
	unavailabilityStore := server.initUnavailabilityStore(mongoClient)
	reservationRequestStore := server.initReservationRequestStore(mongoClient)
	unavailabilityService := server.initUnavailabilityService(unavailabilityStore, producer, reservationRequestStore)
	reservationRequestService := server.initReservationRequestService(reservationRequestStore, unavailabilityService, producer)
	unavailabilityHandler := server.initUnavailabilityHandler(unavailabilityService)
	reservationRequestHandler := server.initReservationRequestHandler(reservationRequestService)
	unavailabilityHandler.Init(server.router)
	reservationRequestHandler.Init(server.router)
	grpcHandler := server.initGrpcHandler(unavailabilityService, reservationRequestService)
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
	store := unavailability.NewUnavailabilityMongoDBStore(client)
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
	//for _, request := range reservationRequests {
	//	_, err := store.Insert(request)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
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

func (server *Server) initUnavailabilityService(store domain.UnavailabilityStore, producer *kafka.Producer, reservationRequestStore domain.ReservationRequestStore) *application.UnavailabilityService {
	return application.NewUnavailabilityService(store, producer, reservationRequestStore, server.loki)
}

func (server *Server) initReservationRequestService(store domain.ReservationRequestStore, unavailabilityService *application.UnavailabilityService, producer *kafka.Producer) *application.ReservationRequestService {
	return application.NewReservationRequestService(store, unavailabilityService, producer, server.loki)
}

func (server *Server) initUnavailabilityHandler(service *application.UnavailabilityService) *api.UnavailabilityHandler {
	return api.NewUnavailabilityHandler(service, server.traceProvider, server.loki)
}

func (server *Server) initReservationRequestHandler(service *application.ReservationRequestService) *api.ReservationRequestHandler {
	return api.NewReservationRequestHandler(service, server.traceProvider, server.loki)
}

func (server *Server) initGrpcHandler(unavailabilityService *application.UnavailabilityService, reservationRequestService *application.ReservationRequestService) *api.BookingHandler {
	return api.NewBookingHandler(unavailabilityService, reservationRequestService, server.traceProvider, server.loki)
}
