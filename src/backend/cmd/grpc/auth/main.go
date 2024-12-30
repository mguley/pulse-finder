package main

import "application"

func main() {
	app := application.NewContainer().InfrastructureContainer.Get()
	authServer := app.AuthServer.Get()
	authService := app.AuthServiceServer.Get()

	// Register the Auth Service
	authServer.RegisterService(authService)

	// Start the gRPC server
	authServer.Start()
	authServer.WaitForShutdown()
}
