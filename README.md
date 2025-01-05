# Mediator for Go - Fx Integration

This project facilitates seamless integration of the Mediator pattern
using [mediator-go](https://github.com/Oleexo/mediator-go) with the [Fx](https://github.com/uber-go/fx) dependency
injection framework. The integration enables the creation of clean, scalable, and maintainable Go applications.

## Getting Started

### Prerequisites

Mediator for Go requires **Go version 1.24** or above.

### Installing Mediator for Go - Fx integration

With [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules), dependencies are handled automatically when
you add the import in your code:

```go
import "github.com/Oleexo/mediator-go-fx"
```

Alternatively, to manually add the package to your project, run:

```sh
go get -u github.com/Oleexo/mediator-go-fx
```

## Usage

### Register a Request Handler

The `mediatorfx.AsRequestHandler` function returns an array of registrations. Below is an example of how to register a
request handler:

```go
package main

import (
	"context"
	"fmt"

	"github.com/Oleexo/mediator-go"
	mediatorfx "github.com/Oleexo/mediator-go-fx"
	"go.uber.org/fx"
)

type MyRequest struct {
	Name string
}

func (r MyRequest) String() string {
	return fmt.Sprintf("MyRequest{Name=%s}", r.Name)
}

type MyResponse struct {
	Result string
}

type MyRequestHandler struct {
}

func NewMyRequestHandler() *MyRequestHandler {
	return &MyRequestHandler{}
}

func main() {
    var constructors []any

    // Register handlers
    constructors = append(constructors, mediatorfx.AsRequestHandler[MyRequest, MyResponse](NewMyRequestHandler)...)

	app := fx.New(
		fx.Provide(constructors...),
		mediatorfx.NewModule(),
		fx.Invoke(func(container mediator.SendContainer) {
			request := MyRequest{}
			response, err := mediator.SendWithoutContext[MyRequest, MyResponse](container, request)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Response: %s", response.Result)
		}),
	)

	app.Run()
}
```

### Register a Notification Handler

The `mediatorfx.AsNotificationHandler` function also returns an array of registrations. Below is an example of how to
register a notification handler:

```go
package main

import (
	"context"
	"fmt"

	"github.com/Oleexo/mediator-go"
	mediatorfx "github.com/Oleexo/mediator-go-fx"
	"go.uber.org/fx"
)

type MyNotification struct {
	Name string
}

type MyNotificationHandler struct {
}

func (MyNotificationHandler) Handle(ctx context.Context, request MyNotification) error {
	fmt.Printf("Handler 1\n")
	return nil
}


func main() {
    var constructors []any

    // Register handlers
    constructors = append(constructors, mediatorfx.AsNotificationHandler[MyNotification, *MyNotificationHandler](NewMyNotificationHandler1)...)

	app := fx.New(
	    fx.Provide(constructors...),
	    mediatorfx.NewModule(),
	    fx.Invoke(func(publisher mediator.Publisher) {
	    	// Send a request
	    	notification := MyNotification{Name: "World"}
	    	err := publisher.Publish(context.Background(), notification)
	    	if err != nil {
	    		// Handle the error
	    		panic(err)
	    	}
	    }),
	)

	app.Run()
}
```

### Register a pipeline behavior

The `mediatorfx.AsPipelineBehavior` is used to register pipeline for request.

```go
package main

import (
	"context"
	
	"github.com/Oleexo/mediator-go"
	mediatorfx "github.com/Oleexo/mediator-go-fx"
	"go.uber.org/fx"
)

type MyPipelineBehavior struct {
	
}

func (p MyPipelineBehavior) Handle(ctx context.Context, request mediator.BaseRequest, next mediator.RequestHandlerFunc) (interface{}, error) {
	// your behavior
	return next()
}

func myPipelineConstructorFunc() MyPipelineBehavior {
	return MyPipelineBehavior{}
}

func main() {
	app := fx.New(
		fx.Provide(
			mediatorfx.AsPipelineBehavior(myPipelineConstructorFunc),
			mediatorfx.NewModule(),
		),
	)

	app.Run()
}
```

### Add module to Fx

The module will register the necessary object like `mediator.Publisher` or `mediator.Sender` to Fx dependency Injection.
The module is named `mediatorfx`.

See more about module in official Fx [documentation](https://uber-go.github.io/fx/modules.html)

```go
package main

import (
	mediatorfx "github.com/Oleexo/mediator-go-fx"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			// Register all request, request handler, pipeline, ...
		),
		mediatorfx.NewModule(),
	)

	app.Run()
}
```

### Complete Example

Below is a complete example that demonstrates how to set up both request/response handlers and notification handlers
with Fx for dependency injection and Mediator:

```go
package main

import (
	"context"
	"fmt"
	"github.com/Oleexo/mediator-go-fx/mediator"
	"go.uber.org/fx"
)

// Define a request and its handler
type MyRequest struct {
	Name string
}

func (r MyRequest) String() string {
	return fmt.Sprintf("MyRequest{Name=%s}", r.Name)
}

type MyResponse struct {
	Result string
}

type MyRequestHandler struct{}

func NewMyRequestHandler() *MyRequestHandler {
	return &MyRequestHandler{}
}

func (h *MyRequestHandler) Handle(ctx context.Context, cmd MyRequest) (MyResponse, error) {
	// Implement your business logic here
	return MyResponse{Result: fmt.Sprintf("Hello, %s!", cmd.Name)}, nil
}

// Define a notification and its handler
type MyNotification struct {
	Message string
}

type MyNotificationHandler1 struct{}

func NewMyNotificationHandler1() *MyNotificationHandler1 {
	return &MyNotificationHandler1{}
}

func (h *MyNotificationHandler1) Handle(ctx context.Context, notification MyNotification) error {
	fmt.Printf("Handler1 received notification: %s\n", notification.Message)
	return nil
}

type MyNotificationHandler2 struct{}

func NewMyNotificationHandler2() *MyNotificationHandler2 {
	return &MyNotificationHandler2{}
}

func (h *MyNotificationHandler2) Handle(ctx context.Context, notification MyNotification) error {
	fmt.Printf("Handler2 received notification: %s\n", notification.Message)
	return nil
}

func main() {
	var constructors []any

	// Register handlers
	constructors = append(constructors, mediatorfx.AsRequestHandler[MyRequest, MyResponse](NewMyRequestHandler)...)
	constructors = append(constructors, mediatorfx.AsNotificationHandler[MyNotification, *MyNotificationHandler1](NewMyNotificationHandler1)...)
	constructors = append(constructors, mediatorfx.AsNotificationHandler[MyNotification, *MyNotificationHandler2](NewMyNotificationHandler2)...)

	// Create and run the application
	app := fx.New(
		fx.Provide(constructors...),
		mediatorfx.NewModule(),
		fx.Invoke(func(container mediator.SendContainer) {
			// Send a request
			request := MyRequest{Name: "World"}
			response, err := mediator.SendWithoutContext[MyRequest, MyResponse](container, request)
			if err != nil {
				// Handle the error
				panic(err)
			}
			fmt.Printf("Response: %s\n", response.Result)

			// Send a notification
			notification := MyNotification{Message: "This is a test notification!"}
			if err := mediator.PublishWithoutContext(container, notification); err != nil {
				// Handle the error
				panic(err)
			}
		}),
	)

	app.Run()
}
```

Learn more in the [👉 main repository](https://github.com/Oleexo/mediator-go).

## Contributing

Contributions are welcome! If you have new ideas or discover ways to enhance the project, 
feel free to submit a pull request. 🌟

## License

This project is licensed under the MIT License. See the LICENSE file for further details.
