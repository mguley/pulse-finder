package main

import "application"

func main() {
	app := application.NewContainer().InfrastructureContainer.Get()
	vacancyServer := app.VacancyServer.Get()
	vacancyService := app.VacancyServiceServer.Get()

	// Register the Vacancy Service
	vacancyServer.RegisterService(vacancyService)

	// Start the gRPC server
	vacancyServer.Start()
	vacancyServer.WaitForShutdown()
}
