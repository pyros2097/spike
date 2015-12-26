package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/tylerb/graceful"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"strings"
	"time"
)

const SigningKey = "jwt@secret"

// type (
// 	GameHub struct {
// 		// Registered connections.
// 		connections map[*websocket.Conn]bool

// 		// Inbound messages from the connections.
// 		broadcast chan []byte

// 		// Register requests from the connections.
// 		register chan *websocket.Conn

// 		// Unregister requests from connections.
// 		unregister chan *websocket.Conn
// 	}
// )

// func NewHub() *GameHub {
// 	return &GameHub{
// 		broadcast:   make(chan []byte),
// 		register:    make(chan *websocket.Conn),
// 		unregister:  make(chan *websocket.Conn),
// 		connections: make(map[*websocket.Conn]bool),
// 	}
// }

func createToken(userId string, scopes []string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["userId"] = userId
	// token.Claims["scopes"] = scopes
	token.Claims["exp"] = time.Now().Add(time.Hour * 96).Unix()
	tokenString, _ := token.SignedString([]byte(SigningKey))
	return tokenString
}

func checkToken(str string) {
	t, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		// Always check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Return the key for validation
		return []byte(SigningKey), nil
	})
	if err == nil && t.Valid {
		// Store token claims in echo.Context
		// c.Set("claims", t.Claims)
		// return nil
	}
	// return unauthorized
}

func uuid4() string {
	var id [16]byte
	rand.Read(id[:])
	id[6] = (id[6] & 0x0f) | (4 << 4)
	id[8] = (id[8] & 0xbf) | 0x80
	return hex.EncodeToString(id[:])
}

func main() {
	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			entry := log.WithFields(log.Fields{
				"uuid":   uuid4(),
				"path":   c.Request().RequestURI,
				"method": c.Request().Method,
				"ip":     c.Request().RemoteAddr,
			})

			if reqID := c.Request().Header.Get("X-Request-Id"); reqID != "" {
				entry = entry.WithField("request_id", reqID)
			}

			entry.Info("started handling request")

			if err := next(c); err != nil {
				c.Error(err)
			}

			latency := time.Since(start)

			entry.WithFields(log.Fields{
				"status":      c.Response().Status(),
				"text_status": http.StatusText(c.Response().Status()),
				"took":        latency,
				fmt.Sprintf("measure#%s.latency", "web"): latency.Nanoseconds(),
			}).Info("completed handling request")

			return nil
		}
	})
	e.Use(mw.Recover())
	startTime := time.Now()

	if time.Since(startTime).Minutes() >= 3 {
		// finish game
	}

	e.Static("/", "public")
	e.WebSocket("/ws", func(c *echo.Context) (err error) {
		ws := c.Socket()
		if err = websocket.Message.Send(ws, "Please Login"); err != nil {
			return err
		}
		msg := ""
		for {
			// if err = websocket.Message.Send(ws, "Hello, Client!"); err != nil {
			// return err
			// }
			if err = websocket.Message.Receive(ws, &msg); err != nil {
				if err == io.EOF {
					return nil
				}
			}
			if strings.Contains(msg, "login") {
				// split get tokenString
				// checkToken(tokenString)
				// if valid {
				// login
				// } else {
				// return
				// }
			} else if strings.Contains(msg, "join_game") {
			} else if strings.Contains(msg, "leave_game") {
			}
			log.Info(msg)
		}
		ws.Close()
		return nil
	})
	log.Info("Starting websocket server at localhost:4000")
	graceful.ListenAndServe(e.Server(":4000"), 3*time.Minute)
}
