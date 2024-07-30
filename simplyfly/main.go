package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"simplifly/common/constants"
	"simplifly/common/kafka/kafkaConsumer"
	"simplifly/internal/api"
	"simplifly/internal/handler"
	"strings"
	"sync"
	"time"
)

func main() {
	logLevel := strings.ToLower(strings.TrimSpace("INFO"))
	lvlMap := map[string]log.Lvl{
		"debug":   log.DEBUG,
		"info":    log.INFO,
		"warning": log.WARN,
		"error":   log.ERROR,
	}
	if _, isPresent := lvlMap[logLevel]; !isPresent {
		logLevel = "debug"
	}
	log.Printf("Log level: %s", logLevel)
	log.SetLevel(lvlMap[logLevel])

	e := HttpServer()
	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server: ", err)
		}
	}()

	kafkaTopicUserUpdates := "user_updates_topic"
	err := kafkaConsumer.StartConsumer(context.Background(), kafkaTopicUserUpdates, KafkaMessageProcessor)
	if err != nil {
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func HttpServer() *echo.Echo {
	h := handler.Handler{API: api.API{}}
	e := echo.New()
	e.HideBanner = true
	// middleware to recover from panic
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// maximum limit 25mb
	e.Use(middleware.BodyLimit("25M"))

	// Enabling Cross Origin Resource sharing
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // todo update before deploying
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderOrigin,
			echo.HeaderXRequestedWith,
			"*",
		},
	}))

	// flights
	e.POST(constants.RegisterFlight, h.RegisterFlightHandler)
	e.POST(constants.UpdateFlight, h.UpdateFlightHandler)
	e.GET(constants.GetFlightDetails, h.GetFlightDetailsHandler)
	e.GET(constants.GetAllFlights, h.GetAllFlightsHandler)

	// users
	e.GET(constants.GetUserDetails, h.GetUserDetailsHandler)
	e.POST(constants.UpsertUser, h.UpsertUserHandler)
	e.GET(constants.GetAllUsers, h.GetAllUsersHandler)
	e.POST(constants.BookFlightByUser, h.BookFlightByUserHandler)

	// websocket
	e.GET("/ws", func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		connMutex.Lock()
		connections = append(connections, ws)
		connMutex.Unlock()

		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				connMutex.Lock()
				// Remove the connection from the slice when it's closed
				for i, conn := range connections {
					if conn == ws {
						connections = append(connections[:i], connections[i+1:]...)
						break
					}
				}
				connMutex.Unlock()
				break
			}
		}

		return nil
	})

	return e
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connMutex sync.Mutex
var connections = make([]*websocket.Conn, 0)

func KafkaMessageProcessor(message []byte) {

	fmt.Println("Kafka Message Received - ", string(message))

	connMutex.Lock()
	defer connMutex.Unlock()

	for _, conn := range connections {
		err := conn.WriteJSON(string(message))
		if err != nil {
			log.Error("Error writing JSON:", err)
			conn.Close()
		}
	}
}
