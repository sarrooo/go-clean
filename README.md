<h1 align="center">Go-Clean</h1>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/sarrooo/go-clean"><img src="https://goreportcard.com/badge/github.com/sarrooo/go-clean" /></a>
  <a href="https://github.com/sarrooo/go-clean/discussions"><img src="https://img.shields.io/badge/chat-on%20discussion-yellow"></a>
  <a href="https://github.com/sarrooo/go-clean/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-GPL%20v3.0-brightgreen.svg" /></a>
</p>

<p align="center">
  Go API boilerplate, including Clean Architecture, Testing, Gin and GORM.
</p>

# Table of Contents

- [Getting Started](#getting-started)
- [Tools](#tools)
- [Makefile Targets](#makefile-targets)
- [Architecture](#architecture)
- [View Models](#view-models)
- [Repository](#repository)
- [Error Handling](#error-handling)
- [Tests](#tests)
- [Naming](#naming-1)
- [Feedbacks](#feedbacks)
- [License](#license)
- [Author](#author)

# Getting Started

```bash
git clone git@github.com:sarrooo/go-clean.git
cd go-clean
make build
```

# Tools

- Contenerization with [Docker](https://www.docker.com/), including [caching](https://docs.docker.com/build/cache/) for faster builds.
- [Gin](https://gin-gonic.com/) for routing.
- [GORM](https://gorm.io/) for ORM, using [PostgreSQL](https://www.postgresql.org/) as database.
- [Go-Swagger](https://github.com/go-swagger/go-swagger) for OpenAPI 2.0 (Swagger) specification.
- [Zap](https://github.com/uber-go/zap) for logging.
- [Viper](https://github.com/spf13/viper) for configuration files.
- [Testify](https://github.com/stretchr/testify) for testing.
- [Mockery](https://vektra.github.io/mockery/latest/) for mocks generation.
- [golangci-lint](https://golangci-lint.run/) for linting.
- [validator](https://github.com/go-playground/validator) for data validation.

# Architecture Overview

This architectural design adheres to the principles of the Clean Architecture, emphasizing separation of concerns and maintainability. The system is structured into three key layers:

- **Controllers:** Responsible for handling incoming requests, interacting with services, and managing the flow of data to and from the client.

- **Services (Use Cases):** Encompasses the business logic and use cases of the application. Services handle the core functionality, ensuring a clear distinction between application-specific rules and external concerns.

- **Repositories (Entities):** Manages the data storage and retrieval, providing a clean interface for the application to interact with the underlying database entities.

Key Features of Our Architecture:

- **Centralized Input and Output Handling:** The architecture ensures a centralized approach to handling request input and response output, promoting consistency and clarity in data flow.

- **Centralized and Normalized Error Handling:** We have implemented a unified error handling strategy across the entire system. By using normalized error codes, we enhance the predictability and manageability of error scenarios.

- **Robust Testing Strategy:** Our testing strategy is comprehensive, covering various levels:
  - Unit tests ensure the individual components function as expected.
  - Table Driven Tests enhance code readability and maintainability.
  - Mocks generation simplifies testing by creating mock objects for dependencies.
  - Code coverage analysis ensures a thorough examination of the codebase.

- **API Documentation:** A well-documented API is crucial for developers and users alike. This architecture includes API documentation, making it easier to understand and interact with the exposed endpoints.

## Folder Structure 📁

```bash
.
├── Dockerfile
├── .mockery.yaml => Mockery configuration file
├── LICENSE
├── Makefile
├── README.md
├── cmd
│   └── main.go
├── docker-compose.yml
├── docs
│   └── swagger.yaml => Generated YAML file
├── go.mod
├── go.sum
├── internal
│   ├── database => Handles database connection (Gorm)
│   ├── errcode => Defines all error types and codes
│   ├── controllers => Contains API handlers and middlewares
│   ├── logger => Configures the logger
│   ├── models => Houses database data models
│   ├── dto => Defines intermediate data models
│   ├── repositories => Manages repositories interacting with the database
│   ├── services => Implements business logic
│   ├── ... => Other packages
│   └── viewmodels => Defines request view models
└── mocks => Generated mock code
```

Our folder structure is designed for clarity and modularity, ensuring each component resides in its designated location for ease of navigation and maintenance. This organized structure contributes to a scalable and maintainable codebase.

## Resources 🪵

- [https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) : Blog post that explains principles of the **Clean Architecture**
- [https://www.youtube.com/watch?v=goC-gCNWhS4](https://youtu.be/goC-gCNWhS4?si=C1W1cdv5_oLD4hVI) : Youtube video that give us a good example of clean architecture implementation to make a Golang API
- [github.com/ruslantsyganok/clean_arcitecture_golang_example](http://github.com/ruslantsyganok/clean_arcitecture_golang_example) : Github repository of the video above
- [https://irahardianto.github.io/service-pattern-go/](https://irahardianto.github.io/service-pattern-go/) : Blog post that explains also Golang API implementation in lines with the principles of the **Clean Architecture**

# Makefile Targets

This project includes a Makefile to streamline common development tasks. Here are the available targets and their purposes:

- **`run`**: Starts the Docker containers for the application.

- **`build`**: Builds and starts the Docker containers, ensuring that the images are up-to-date.

- **`stop`**: Stops the running Docker containers.

- **`unit-test`**: Executes unit tests for the project, providing coverage information.

- **`func-test`**: Runs functional tests using the `tests/server.sh` script. Specify the desired port with `PORT=<port>`.

- **`lint`**: Performs linting using `goimports` and `golangci-lint` to ensure code quality.

- **`swagger`**: Generates the OpenAPI specification file (`docs/swagger.yaml`) using the `swagger` tool.

- **`serve-swagger`**: Serves the generated Swagger documentation on a local server.

- **`generate-mocks`**: Generates mock code using the `mockery` tool for easier testing.

# View Models
The **viewmodel** package defines all data type used by API handlers.

For each request we must define 2 view models : **Request (input)** if the route has at least one parameters and the **Response (output)**.

<aside>
⚠️ View model type are only used in **handler** and **service** packages, not at all in **repository** package.

</aside>

## Naming

each viewmodel must end with the suffix `Request` or `Response` , and the first part of the viewmodel name must match the first part of the Controller that uses its viewmodels.

e.g

```go
type RegisterOutput struct { ... }       // BAD, IT SHOULD FINISH BY Response
type RegisterUserResponse struct { ... } // BAD, THE CONTROLLER is registerController
                                         // not registerUserController
type RegisterResponse struct { ... }     // GOOD !
type RegisterRequest struct { ... }      // GOOD
```

## Documentation

View models allow us to document our request input and output using [go-swagger](https://github.com/go-swagger/go-swagger) annotation in the code. You can find out more details about annotation in this go-swagger’s [documentation part](https://goswagger.io/use/spec.html).

The naming convention for `// swagger:response` and `// swagger:response` is to name by the name of the controller

### Example

```go
// swagger:parameters **registerController**
type RegisterUserRequest struct { ... }

// swagger:response **registerController**
type RegisterUserRequest struct { ... }

// swagger:route POST /auth/register auth **registerController**
//
// Endpoint for user registration.
//
// responses:
//
//	200: **registerController**
//	400: errorResponse
func registerController(svc services.ServiceInterface) gin.HandlerFunc { ... }
```

<aside>
⚠️ All request and response view models must be documented

</aside>

## Binding

### Request

I wanted to have a behavior similar to [grpc](https://grpc.io/) where the request is binded in the handler, instead of verify all parameters in the handler. It's why there is a binding middleware.

Before each request, if the request has at least one parameter, you must call `requestViewmodelMiddleware`, this middleware will bind request parameters to request view model, apply verification tags, etc. It allow the developer to not check further in the handler. After the binding middleware, if there was no error, you can access to it in the handler by a context variable name *requestViewmodel.*

For this binding we use the [binding feature of Gin](https://gin-gonic.com/docs/examples/binding-and-validation/), and it has few limitations :

- Binding of URI parameters (`/country/:id` in this example *id* is an URI parameter) is made not made by `ShouldBind`, but I write a function that bind it.
- We use [ShouldBind](https://pkg.go.dev/github.com/gin-gonic/gin@v1.9.0#Context.ShouldBind) method and it bind depending of the **Method** and the **Content-Type** headers. It mean that you can’t bind **form** parameters (`/country?sort_by=name` in this example *sort_by*  is an form parameter) when the request use **POST** method or inversely you can’t bind **JSON/Body** parameters when the request use **GET** method. Be careful, test your code to ensure the binding is correct.

### Response

After each request, if there wasn’t error, a response binding middleware is called, `responseViewmodelMiddleware`. It gets two context value, *statusCode and* *****responseViewmodel*, and respond to the client. At the end of each handler you must define these two context values, using Gin method, `c.Set("responseViewmodel", response)`.

### Tags

Gin `ShouldBind()` method uses [validator](https://github.com/go-playground/validator) library. You must define binding rules using `binding` tags, check out the example below.

```go
type RegisterUser struct {
	// The email of the user.
	Email string `json:"email" binding:"required,email"`

	// The password of the user.
	Password string `json:"password" binding:"required,min=8,max=64"`

	// The first name of the user.
	FirstName string `json:"first_name" binding:"required"`

	// The last name of the user.
	LastName string `json:"last_name" binding:"required"`

	// The phone number of the user.
	Phone string `json:"phone" binding:"omitempty,e164"`

	// The birth date of the user.
	BirthDate string `json:"birth_date" binding:"omitempty,datetime"`
}
```

In this example each fields will be validate by validator library. By example if the email is not at the right format, a error will be triggered, in different language.

<aside>
💡 It’s recommended to use binding tags rather than in business logic methods. By example validator library provided translation for error.

</aside>

## Examples

<details>

<summary>Declaring view models</summary>

- `GET /country/list?page=1&limit=2&sort_by=name`
    
    ```go
    // swagger:parameters listCountryController
    type ListCountryRequest struct {
    	// The page number for pagination
    	// in:query
    	Page int `form:"page"`
    
    	// The number of items to retrieve (0 means all).
    	// in:query
    	Limit int `form:"limit"`
    
    	// The field to sort by (default is updated_at).
    	// in:query
    	SortBy string `form:"sort_by"`
    
    	// The order of sorting (asc or desc, default is desc).
    	// in:query
    	Order string `form:"order"`
    }
    ```
    
- `GET /country/1`
    
    ```go
    // swagger:parameters getCountryController
    type GetCountryRequest struct {
    	// The ID of the country.
    	// in:path
    	ID uint `uri:"id"`
    }
    ```
    
- `POST /country | body: {”name”: “France”}`
    
    ```go
    // swagger:parameters postCountryController
    type PostCountryRequest struct {
    	// in:body
    	Body struct {
    		// The name of the country.
    		Name string `json:"name"`
    	}
    }
    ```

<aside>
💡 You can split your view models in several other data type that will be used in other request and response view model. These sub data types must be documented.

</aside>

</details>

<details>

<summary>Binding middleware</summary>

**Request view model**

```go
func registerCountryRoutes(group *gin.RouterGroup, service service.ServiceInterface) {
	group.GET("/list", requestViewmodelMiddleware(&viewmodel.ListCountryRequest{}), listCountriesController(service))
}
```

**Response view model**

```go
func listCountriesController(service service.ServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := c.MustGet("requestViewmodel").(*viewmodel.ListCountryRequest)
		response := &viewmodel.ListCountryResponse{}

		...

		c.Set("statusCode", http.StatusOK)
		c.Set("responseViewmodel", response)
	}
}
```
</details>


# Repository

Each resources have it’s own **repository**. The repository is the only one-way to interact with database entity.

## Interface

The repository is an interface that defines all methods to interact with database like the example below :

```go
type CountryRepositoryInterface interface {
	CreateCountry(country *models.Country) error
	GetCountries() (*[]models.Country, error)
	GetCountryByID(id uint) (*models.Country, error)
	GetCountryByName(name string) (*models.Country, error)
	UpdateCountry(country *models.Country) error
}
```

The repository pattern helps in separating the logic that retrieves data from the database from the business logic of the application. This promotes cleaner, more maintainable code.

The business logic have access to a GlobalRepository that implements all resources repository.

# Error Handling

We decided to have a central error handling strategy, handled by one middleware called `errorHandlerMiddleware`.

## Package

We created **errcode** package to standardize errors. The package is very simple, it just define a new type called `GoCleanError` and a list of predefined `GoCleanError`.

```go
type GoCleanError struct {
	error
}

var (
	ErrInvalidParameters   = GoCleanError{errors.New("invalid parameters")}
	ErrNotFound            = GoCleanError{errors.New("not found")}
	ErrUnknown             = GoCleanError{errors.New("unknown error")}
	...
)
```

The Go 1.13 error handling introduce a new error feature, the **wrapping**. To wrap an error in other error you just could do that `fmt.Errorf("%w: %v", error_code.ErrInvalidParameters, err)`. And with method like `errors.Is(…)` or `errors.As(…)` you can know if the error is `GoCleanError` by example.

<aside>
💡 You are free to create new `GoCleanError` type according to your needs

</aside>

## Middleware

When an error from an external method, a method that you didn’t code, you must wrap this error with a Betrip error like the example below, and return it.

```go
err = repo.DB.Find(&countries).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", error_code.ErrDatabase, err)
	}
```

Then the error will go up to the top level, the handler, if you don’t want to handle the error at a sub level. In the handler you must set the Gin context with this error like that.

```go
countries, err := service.ListCountries()
		if err != nil {
			c.Error(err)
			return
		}
```

Then the `errorHandlerMiddleware` will catch the error and handle it. It logs the full error chain and if there isn’t `GoCleanError` in the chain it will return an internal error status to the client otherwise it will catch the most top level `GoCleanError` in the error chain and return it’s message to the client.

This strategy makes it safe because the client will not have to much information on the error but the developer will have all error information.

## Ressources 🪵
<aside>
💡 These ressources are very useful and we advise you to read/watch them to better understand the error management strategy in place.

</aside>

- [https://www.youtube.com/watch?v=IKoSsJFdRtI](https://www.youtube.com/watch?v=IKoSsJFdRtI): Video made by Software engineer at Google that explains how deal with error since Go 1.13.
- [https://go.dev/blog/go1.13-errors](https://go.dev/blog/go1.13-errors): Blog article related to the video.
- [https://www.joeshaw.org/error-handling-in-go-http-applications/](https://www.joeshaw.org/error-handling-in-go-http-applications/): Blog article that explain the sentinel error strategy to have a safe strategy error handling when developing API.
- [https://pkg.go.dev/errors](https://pkg.go.dev/errors): Documentation of the standard errors package


# Tests

Thanks to our architecture layout is quite simple to test our code because all service are mocked and the code is well splited.

We focus our test on handler and service package because it’s the core of our system.

We use [Testify](https://github.com/stretchr/testify) librairie to test our code. This librairie provide features like assertions, mocking, suite, etc.

In addition of Testify we use [Mockery](https://github.com/vektra/mockery). This librairie greatly simplifies mocking and avoid boilerplate, it generate mock type for each interface in our code and we can define the behaviour of each methods during test scenario.

## Rules

### Levels

For units we use 3 levels : Suite → Test → SubTest

- **Suite** : **There is 1 suite by package**. The suite is shared accros each tests and sub tests. We can store mocks, service, logger, etc, in the suite to use it in tests. By example `ControllerSuiteTest`. Check the suite documentation, and don’t hesitate to use suite hook like `SetupSubTest` or `TearDownSubTest` by example.
- **Test** : **There is 1 test by method**. The test focus on a method, by example `TestRegisterController` will test only the register controller, etc.
- **SubTest** : There are multiple sub tests by test. A subtest test one path, by example the sub test `Success` of `TestRegisterController` test how the behavior of the function if there aren’t error during the execution.

These 3 levels must be used for each package tested.

### **Table Driven Tests**

To improve the code readability we choose to use [table driven test](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests). Each table contains multiple subtest.

- Example
    
    ```go
    tests := controllerTestTable{
    		"Success": {
    			setupMock: func() {
    				suite.svc.On("RegisterUser", &user).Return(modelUser, nil)
    				suite.svc.On("GenerateToken", modelUser).Return("token", nil)
    			},
    			requestViewmodel: &viewmodel.RegisterUserRequest{
    				Body: user,
    			},
    			expected: controllerTestExpected{
    				status: http.StatusOK,
    				responseViewmodel: &viewmodel.RegisterUserResponse{
    					Body: struct {
    						Token        string "json:\"token\""
    					}{
    						Token:        "token",
    					},
    				},
    			},
    		},
    }
    ```
    

### Mock

- Before each subtest you must reset the expected calls of each mocks, Ex: `suite.svc.ExpectedCalls = nil`
- During each subtest you must check that all mock method was called, Ex: `suite.svc.AssertExpectations(suite.T())`

## Controllers

Because each controller have the same structure, subtest execution, assertion, table test type (expected, parameters, …) is handled globally managed for whole package.

```go
func xxxController(svc services.ServiceInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := ctx.MustGet(ContextKeyRequestViewmodel).(*viewmodel.xxxRequest)
		response := &viewmodel.xxxResponse{}
	}
}
```

### Exception

There is a exception for middleware testing. Middleware are is controllers folder but don’t have same testing structure. Because each middleware have it’s own context, parameters, output, etc, you can not use suite and handle test specifically for each middleware. Don’t hesitate to check the `middleware_test.go` file to know how to test a middleware.

## Service

Because each service method have specific structure (expected, parameters, …), each tests need to handle there own structure type, and to execute there own tests

```go
func (suite *ServiceSuiteTest) TestRegisterUser() {
	type parametersType struct {
		registerUser *dto.RegisterUser
	}

	type expectedType struct {
		user *models.User
		err  error
	}
	
	tests := map[string]struct {
		setupMock  func()
		parameters parametersType
		expected   expectedType
	}{ ... }

	...

	for testName, test := range tests {
		suite.Run(testName, func() {
			test.setupMock()

			_, err := suite.svc.RegisterUser(test.parameters.registerUser)

			suite.Assert().Equal(test.expected.err, err)
		})
	}
}
```

## Resources 🪵

- [https://pkg.go.dev/github.com/stretchr/testify](https://pkg.go.dev/github.com/stretchr/testify): Testify documentation
- [https://pkg.go.dev/github.com/stretchr/testify/suite](https://pkg.go.dev/github.com/stretchr/testify/suite): Testify suite documentation
- [https://vektra.github.io/mockery/latest/](https://vektra.github.io/mockery/latest/): Mockery documentation

# Naming

Naming rules must be respected to keep the code **homogeneous and harmonious**.

In addition, as a general rule, the name of an element must be **clear about the purpose** for which it is used, and **must not be open to various interpretations**.

e.g

```go
// DTO file

 type ListParams struct { ... } // GOOD, IT IS CLEAR ON WHAT IT CONTAINS
 type TokensResponse struct { ... } // BAD, DTO SHOULD NOT USE Response
                                    // SINCE IT IS USED FOR viewmodels PKG
```

## Interfaces

Each interface must end with the suffix `Interface` , e.g. `CountryRepositoryInterface`.

## Controllers

Each controller must end with the suffix `Controller` , e.g. `registerController`.

## Repositories

Each repository must end with the suffix `Repository` , e.g. `CountryRepository`.

## Viewmodels

each viewmodel must end with the suffix `Request` or `Response` , and the first part of the viewmodel name must match the first part of the Controller that uses its viewmodels.

e.g

```go
type RegisterOutput struct { ... }       // BAD, IT SHOULD FINISH BY Response
type RegisterUserResponse struct { ... } // BAD, THE CONTROLLER is registerController
                                         // not registerUserController
type RegisterResponse struct { ... }     // GOOD !
type RegisterRequest struct { ... }      // GOOD
```

## Errors

Each `errcode` start with `Err` and is followed by the kind of the error, e.g `ErrRestrictedArea`

## DTO (Data Transfer Object)

DTO structures do **not follow specific rules,** but they must **avoid collisions with the rules of other packages** (viewmodels, errcode, repositories, controllers ...).

# Feedbacks

This repository was made to share our experience and to solve many problems that we encountered during our projects. Feel free to give us your feedbacks in the [discussion section](https://github.com/sarrooo/go-clean/discussions).

# License

[GNU General Public License v3.0](https://github.com/resotto/goilerplate/blob/master/LICENSE).

# Author

[Sarrooo](https://github.com/sarrooo)
