package constants

import (
	"path"
)

const (
	MongoConnectionString = "mongodb+srv://chemotherapy2000:4J8m1975%23@simplyfly-cluster.yjmgyvd.mongodb.net/?retryWrites=true&w=majority&appName=simplyfly-cluster"
	KafkaBroker           = "outgoing-stingray-7987-eu2-kafka.upstash.io:9092"

	// Response status
	Success = "SUCCESS"
	Failure = "FAILURE"

	// Update Type
	Cancellation = "CANCELLATION"
	Delay        = "DELAY"
	GateChange   = "GATE_CHANGE"

	RedisExpiry = 24
)

var (
	//prefix
	Root = "/"
	//entities
	Flight  = "flight"
	Flights = "flights"
	User    = "user"
	Users   = "users"
	Id      = ":id"
	Details = "details"
	//actions
	Register = "register"
	Update   = "update"
	Upsert   = "upsert"
	Book     = "book"

	// flights
	RegisterFlight   = path.Join(Root, Flight, Register)
	UpdateFlight     = path.Join(Root, Flight, Update)
	GetFlightDetails = path.Join(Root, Flight, Id)
	GetAllFlights    = path.Join(Root, Flights)

	// users
	GetUserDetails   = path.Join(Root, User, Details)
	UpsertUser       = path.Join(Root, User, Upsert)
	GetAllUsers      = path.Join(Root, Users)
	BookFlightByUser = path.Join(Root, User, Flight, Book)
)
