# Integration Boilerplate (Golang)

## Overview
The Qiscus Integration team's Go project boilerplate is designed to speed up application development and ensure consistency across all integration projects. By providing a standardized structure and set of guidelines, this boilerplate makes it easier for developers to understand, review, and maintain projects. By using this boilerplate, developers can ensure that all integration projects follow a consistent structure and coding style, making it easier to onboard new team members and maintain projects in the long term.

This boilerplate is used to standardize the directory structure for projects of medium to large complexity or with the potential for it. However, for other cases that only handle one or a few processes, it is not necessary to implement this boilerplate in order to avoid over-abstraction. For example, you can use a single file like main.go or a flat architecture instead.

## Directories
### `/cmd`
Main applications for this project. The directory name for each application should match the name of the executable you want to have. Don't put a lot of code in the directory.

### `/common`
The directory is used to hold code that is shared across different parts of the application. The common directory may contain utility functions, constants, error types, database connection and other modules that are used by multiple packages within the application. The purpose of the common directory is to avoid code duplication and to keep the shared code organized in one place.

### `/cron`
The cron directory contains code related to running scheduled tasks or background jobs using the operating system's cron scheduler. The directory may include files for defining and configuring the cron jobs, as well as the code to execute the tasks.

### `/delivery`
This directory will act as the presenter layer. Decide how the data will presented. Could be as REST API, or HTML File, or gRPC whatever the delivery type. This layer also will accept the input from user. Sanitize the input and sent it to Usecase layer.

### `/domain`
The domain layer is a representation of the application's business logic in the code. It contains the entities and value objects that the application uses to model its business concepts and rules. The domain layer is independent of any specific technology or implementation details and is designed to be reusable and independent.

### `/repository`
This directory containing adapters to different storage implementations. A data source might be an adapter to a SQL database, an elastic search adapter, or REST API. A data source implements methods defined on the repository and stores the implementation of fetching and pushing the data.

### `/server`
The server directory is contains code related to setting up and running the application's HTTP server.

### `/usecase`
This directory will act as the business process layer, any process will handled here. This layer will decide, which repository layer will use. And have responsibility to provide data to serve into delivery. Process the data doing calculation or anything will done here. Usecase layer will accept any input from Delivery layer, that already sanitized, then process the input could be storing into DB , or Fetching from DB ,etc.

## Sample Use Case
Create tagging and storing room data in the database when a user initiates a chat on Qiscus Omnichannel by utilizing a new session webhook, besides that there is a cronjob that does an auto resolved room when the on going room has reached 10 minutes. And also create an API to get all rooms that are stored in the database with an API key authentication.