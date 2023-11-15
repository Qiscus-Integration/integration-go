# Integration Boilerplate (Golang)

## Changelog

- **v1**: checkout to the [v1 branch](https://bitbucket.org/qiscus/integration-go/src/v1/)
  Layered architecture (delivery, usecase, repository) or driven by tech

## Overview

The Qiscus Integration team's Go project boilerplate is designed to speed up application development and ensure consistency across all integration projects. By providing a standardized structure and set of guidelines, this boilerplate makes it easier for developers to understand, review, and maintain projects. By using this boilerplate, developers can ensure that all integration projects follow a consistent structure and coding style, making it easier to onboard new team members and maintain projects in the long term.

This boilerplate is used to standardize the directory structure for projects of medium to large complexity or with the potential for it. However, for other cases that only handle one or a few processes, it is not necessary to implement this boilerplate in order to avoid over-abstraction. For example, you can use a single file like main.go or a flat architecture instead.

## Directories

[TODO]

## Create New Module

[TODO]

## Sample Use Case

You can find the documentation for this application [here](/docs/README.md)

Create tagging and storing room data in the database when a user initiates a chat on Qiscus Omnichannel by utilizing a new session webhook, besides that there is a cronjob that does an auto resolved room when the on going room has reached 10 minutes. And also create an API to get all rooms that are stored in the database with an API key authentication.

### Environment Variables

To run this project, you will need to add the following environment variables to your `.env` file

```
APP_SECRET_KEY=
DATABASE_HOST=
DATABASE_PORT=
DATABASE_USER=
DATABASE_PASSWORD=
DATABASE_NAME=
QISCUS_APP_ID=
QISCUS_SECRET_KEY=
QISCUS_OMNICHANNEL_URL=
```

### Run Locally

To run the project locally, follow these steps:

- Clone this repository: `git clone git@bitbucket.org:qiscus/integration-go.git`
- Navigate to the directory: `cd integration-go`
- Format code and tidy modfile: `make tidy`
- Run test: `make test`, make sure that all tests are passing
- Run the server: `make run bin=server`, or run the application with reloading on file changes with: `make run/live bin=server`. You can also apply this to the cron application by changing the parameter to `bin=cron`
- The backend server will be accessible at `http://localhost:8080`
- You can find another usefull commands in `Makefile`

### Generate Mock for Service

- Install [Mockery](https://github.com/vektra/mockery)
- Add the following code in the service file: `//go:generate mockery --all --case snake --output ./mocks --exported`
- Run go generate using `make generate`
