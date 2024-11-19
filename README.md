# Integration Boilerplate (Golang)

### Overview

The Qiscus Integration team's Go project boilerplate is designed to speed up application development and ensure consistency across all integration projects. By providing a standardized structure and set of guidelines, this boilerplate makes it easier for developers to understand, review, and maintain projects. By using this boilerplate, developers can ensure that all integration projects follow a consistent structure and coding style, making it easier to onboard new team members and maintain projects in the long term.

This boilerplate is used to standardize the directory structure for projects of medium to large complexity or with the potential for it. However, **for other cases that only handle one or a few processes, it is not necessary to implement this boilerplate** in order to avoid over-abstraction. For example, you can use a single file like main.go or a flat architecture where all components such as handlers, services, and repositories are placed in a single directory for simplicity and ease of navigation. Here’s a repository that implements this kind of flat structure: [Repo](https://bitbucket.org/qiscus/panin/src/main/).

One significant change from the [v1](https://bitbucket.org/qiscus/integration-go/src/v1/) is moving away from grouping code by function and, instead, organizing it by module. This approach offers several advantages:

- **Single Responsibility Principle**: One of the SOLID principles of object-oriented design, states that a class or module should have only one reason to change.
- **Reusability**: Modules become more reusable across different clients, promoting code sharing and reducing duplication.
- **Loose Coupling, High Cohesion**: Two words that describe how easy or difficult it is to change a piece of software. Grouping by module enforces loose coupling between different parts of the code while promoting high cohesion within each module.
- **Faster Contribution**: Developers can contribute to specific modules without causing **collateral damage** in unrelated areas, speeding up the development process.
- **Ease of Understanding**: The codebase becomes more accessible and understandable as it's organized around modules and use cases. A use case repository clarifies what each module does.

### Create New Module

[TODO]

### Sample Use Case

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
#Database Log Level: info, warn, debug, error. by default database log are disabled
DATABASE_LOG_LEVEL=
REDIS_URL=
QISCUS_APP_ID=
QISCUS_SECRET_KEY=
QISCUS_OMNICHANNEL_URL=
```

### Run Locally

To run the project locally, follow these steps:

- Clone this repository
- Navigate to the directory
- Format code and tidy modfile: `make tidy`
- Run test: `make test`, make sure that all tests are passing
- Run the server: `make run bin=api`, or run the application with reloading on file changes with: `make run/live bin=api`. You can also apply this to the cron application by changing the parameter to `bin=cron`
- The backend server will be accessible at `http://localhost:8080`
- You can find another usefull commands in `Makefile`

### Generate Mock from Interface

- Install [Mockery](https://github.com/vektra/mockery)
- Add the following code in the interface code file: `//go:generate mockery --case snake --name XXXX`
- Run go generate using `make generate`

### Handle HTTP Client Exceptions

`client.Error` complies with Go standard error. which support Error, Unwrap, Is, As

```go
// Sample using errors.As
err := s.omni.AssignAgent(ctx, agentID, roomID)
if err != nil {
    var cerr *client.Error
    if errors.As(err, &cerr) {
        fmt.Println(cerr.Message)     		// General error message
        fmt.Println(cerr.StatusCode)  		// HTTP status code e.g: 400, 401 etc.
        fmt.Println(cerr.RawError)    		// Raw Go error object
        fmt.Println(cerr.RawAPIResponse)  // Raw API response body in byte
    }
}

```
