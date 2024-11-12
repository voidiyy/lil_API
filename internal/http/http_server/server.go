package http_server

import (
	"fmt"
	"gigaAPI/config"
	"gigaAPI/internal/db_boss"
	"gigaAPI/internal/http/http_handler"
	"gigaAPI/internal/http/http_server/http_security"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type serverHTTPSConfig struct {
	listenAddr   string
	certHostname string
	timeOut      string
	idleTimeout  string
	certFile     string
	keyFile      string
	jwtSecret    []byte
}

type ServerHTTPS struct {
	config  *serverHTTPSConfig
	echo    *echo.Echo
	handler http_handler.HandlerWorker
	*HTML
}

//TODO: 1.add html template path to config

func InitServer(gc config.GlobalConfig) *ServerHTTPS {
	//config
	conf := &serverHTTPSConfig{
		listenAddr:   gc.HttpConf.ListenAddr,
		certHostname: gc.HttpConf.CertHostname,
		timeOut:      gc.HttpConf.TimeOut,
		idleTimeout:  gc.HttpConf.IdleTimeout,
		certFile:     gc.HttpConf.CertFile,
		keyFile:      gc.HttpConf.KeyFile,
		jwtSecret:    []byte(gc.HttpConf.JWTSecret),
	}
	//database
	db, err := db_boss.InitDB(gc.DbConf.DbType, gc.DbConf.DbURL)
	if err != nil {
		panic("db init: " + err.Error())
	}

	//html templates
	templ, er := NewHTML("../../internal/http/html")
	if er != nil {
		panic("templates init: " + er.Error())
	}

	//router
	e := echo.New()

	//data validation
	e.Validator = &http_handler.CustomValidator{Validator: validator.New()}
	//render html
	e.Renderer = templ

	srv := &ServerHTTPS{
		config:  conf,
		echo:    e,
		handler: http_handler.NewHandlerWorker(db, conf.jwtSecret),
		HTML:    templ,
	}
	srv.setMiddleware(
		middleware.Logger(),
		middleware.RequestID(),
		middleware.Recover(),
		http_security.CSRFTokenHeaderInjector(),
		http_security.SSetCSRF(),
	)

	return srv
}

func (s *ServerHTTPS) Run() error {

	//middlewares

	s.dbPath()
	err := s.echo.StartTLS(s.config.listenAddr, s.config.certFile, s.config.keyFile)

	return err
}

func (s *ServerHTTPS) setMiddleware(middlewares ...echo.MiddlewareFunc) {
	if len(middlewares) == 0 {
		fmt.Println("No middleware functions provided.")
		return
	}

	for i, m := range middlewares {
		s.echo.Use(m)
		fmt.Printf("Added middleware function [%d] of type %s\n", i, fmt.Sprintf("%T", m))
	}
}
