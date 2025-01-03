package mediatorfx

import (
	"github.com/Oleexo/mediator-go"
	"go.uber.org/fx"
)

func AsRequestHandler[TRequest mediator.Request[TResponse], TResponse interface{}](f any) []interface{} {
	return []interface{}{
		fx.Annotate(
			f,
			fx.As(new(mediator.RequestHandler[TRequest, TResponse])),
		),
		fx.Annotate(
			mediator.NewRequestHandlerDefinition[TRequest, TResponse],
			fx.As(new(mediator.RequestHandlerDefinition)),
			fx.ResultTags(`group:"mediator_request_handlers"`),
		),
	}
}
