package vacancy

import (
	"application/dependency"
	"application/event"
	"application/vacancy"
	"domain/vacancy/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	infraVacancy "infrastructure/vacancy"
	"interfaces/api/utils"
	apiHandlers "interfaces/api/vacancy/handlers"
	apiValidators "interfaces/api/vacancy/validators"
)

// Container provides a lazily initialized set of dependencies for the vacancy domain.
type Container struct {
	VacancyRepository dependency.LazyDependency[repository.VacancyRepository]
	VacancyService    dependency.LazyDependency[*vacancy.Service]
	VacancyValidator  dependency.LazyDependency[*apiValidators.RequestValidator]
	CreateHandler     dependency.LazyDependency[*apiHandlers.CreateVacancyHandler]
}

// NewContainer initializes and returns a new Container with lazy dependencies for the vacancy domain.
func NewContainer(db *pgxpool.Pool, d event.Dispatcher, h *utils.Handler, e *utils.Errors) *Container {
	c := &Container{
		VacancyRepository: dependency.LazyDependency[repository.VacancyRepository]{
			InitFunc: func() repository.VacancyRepository {
				return infraVacancy.NewPgxVacancyRepository(db)
			},
		},
	}
	c.VacancyService = dependency.LazyDependency[*vacancy.Service]{
		InitFunc: func() *vacancy.Service {
			return vacancy.NewService(c.VacancyRepository.Get(), d)
		},
	}
	c.VacancyValidator = dependency.LazyDependency[*apiValidators.RequestValidator]{
		InitFunc: apiValidators.NewRequestValidator,
	}
	c.CreateHandler = dependency.LazyDependency[*apiHandlers.CreateVacancyHandler]{
		InitFunc: func() *apiHandlers.CreateVacancyHandler {
			return apiHandlers.NewCreateVacancyHandler(h, e, c.VacancyService.Get(), c.VacancyValidator.Get())
		},
	}

	return c
}