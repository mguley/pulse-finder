#### Vacancy Management System

This project is a modular application for managing and displaying vacancies. 

It consists of two main parts: a frontend built with React and TypeScript and a backend written in Go. 

The application uses PostgreSQL for data storage and gRPC for communication with the [Pulse Finder Bot](https://github.com/mguley/pulse-finder-bot), which supplies vacancy data.

#### Features
- **Frontend**:
  - Built with React, TypeScript, and MUI.
  - Displays a grid of vacancies with filtering options by company name and vacancy title.
  - Deployed via GitHub Pages for easy access.
- **Backend**:
  - Written in Go, following a Domain-Driven Design (DDD) approach.
  - Provides REST API endpoints:
    - `/v1/jwt`: For authentication.
    - `/v1/vacancies`: For retrieving vacancy data.
  - Handles gRPC communication to receive data from the [Pulse Finder Bot](https://github.com/mguley/pulse-finder-bot).
  - Stores vacancy data in PostgreSQL.
- **Infrastructure**:
  - Hosted on a Digital Ocean droplet.
  - Deployed using Terraform for consistent and automated infrastructure management.
 
#### Workflow

1. **Data Ingestion**:
   - The [Pulse Finder Bot](https://github.com/mguley/pulse-finder-bot) parses vacancy data and transmits it via gRPC to this backend.
   - The backend saves the data into a PostgreSQL database.
2. **Frontend Interaction**:
   - A user accesses the GitHub Pages deployment.
   - The frontend requests a JWT token from the `/v1/jwt` endpoint.
   - Using the token, the frontend calls the `/v1/vacancies` endpoint to retrieve vacancy data.
   - The frontend renders the vacancies using MUI components, allowing filtering by company name or title.

#### Directory Structure

##### Frontend (`src/domains`)
- `Dashboard`: Displays the grid of vacancies.
- `Jobs`: Handles job-specific UI logic.
- `NavigationBar`: Manages the navigation UI.

##### Backend (`src/backend`)
```
backend/
├── application/         # Core business logic
│   ├── auth/            # Authentication services
│   ├── vacancy/         # Vacancy-related logic
│   └── route/           # API endpoint definitions
├── cmd/                 # Entry points for services
│   ├── grpc/            # gRPC servers
│   └── main/            # Main service entry point
├── domain/              # Domain models and entities
├── infrastructure/      # Persistence, migrations, and gRPC handlers
├── interfaces/          # REST API endpoints and middleware
└── tests/               # Integration and unit tests
```

#### Demo

https://mguley.github.io/pulse-finder/
