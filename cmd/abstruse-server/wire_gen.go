// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jkuri/abstruse/internal/pkg/auth"
	"github.com/jkuri/abstruse/internal/pkg/config"
	"github.com/jkuri/abstruse/internal/pkg/http"
	"github.com/jkuri/abstruse/internal/pkg/log"
	"github.com/jkuri/abstruse/internal/server"
	"github.com/jkuri/abstruse/internal/server/controller"
	"github.com/jkuri/abstruse/internal/server/db"
	"github.com/jkuri/abstruse/internal/server/db/repository"
	"github.com/jkuri/abstruse/internal/server/etcd"
	"github.com/jkuri/abstruse/internal/server/grpc"
	"github.com/jkuri/abstruse/internal/server/service"
	"github.com/jkuri/abstruse/internal/server/websocket"
)

// Injectors from wire.go:

func CreateApp(cfg string) (*server.App, error) {
	viper, err := config.NewConfig(cfg)
	if err != nil {
		return nil, err
	}
	options, err := server.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	logOptions, err := log.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	logger, err := log.New(logOptions)
	if err != nil {
		return nil, err
	}
	httpOptions, err := http.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	websocketOptions, err := websocket.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	dbOptions, err := db.NewOptions(viper, logger)
	if err != nil {
		return nil, err
	}
	gormDB, err := db.NewDatabase(dbOptions)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewDBUserRepository(logger, gormDB)
	userService := service.NewUserService(logger, userRepository)
	userController := controller.NewUserController(logger, userService)
	versionService := service.NewVersionService(logger)
	versionController := controller.NewVersionController(logger, versionService)
	initControllers := controller.CreateInitControllersFn(userController, versionController)
	router := http.NewRouter(httpOptions, websocketOptions, initControllers)
	httpServer, err := http.NewServer(httpOptions, logger, router)
	if err != nil {
		return nil, err
	}
	etcdOptions, err := etcd.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	etcdServer := etcd.NewServer(etcdOptions, logger)
	websocketServer := websocket.NewServer(websocketOptions, logger)
	grpcOptions, err := grpc.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	app := websocket.NewApp(logger)
	grpcApp, err := grpc.NewApp(grpcOptions, app, logger)
	if err != nil {
		return nil, err
	}
	serverApp := server.NewApp(options, logger, httpServer, etcdServer, websocketServer, grpcApp)
	return serverApp, nil
}

// wire.go:

var providerSet = wire.NewSet(log.ProviderSet, config.ProviderSet, http.ProviderSet, etcd.ProviderSet, server.ProviderSet, db.ProviderSet, repository.ProviderSet, auth.ProviderSet, controller.ProviderSet, service.ProviderSet, websocket.ProviderSet, grpc.ProviderSet)
