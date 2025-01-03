package mediatorfx

import (
	"github.com/Oleexo/mediator-go"
	"go.uber.org/fx"
)

type SendContainerParams struct {
	fx.In

	RequestHandlers []mediator.RequestHandlerDefinition `group:"mediator_request_handlers"`
	Pipelines       []mediator.PipelineBehavior         `group:"mediator_pipelines"`
}

type PublisherParams struct {
	fx.In

	NotificationHandlers []mediator.NotificationHandlerDefinition `group:"mediator_notification_handlers"`
	PublishStrategy      mediator.PublishStrategy                 `optional:"true"`
}

func NewSendContainer(param SendContainerParams) mediator.SendContainer {
	return mediator.NewSendContainer(mediator.WithRequestDefinitionHandlers(param.RequestHandlers...),
		mediator.WithPipelineBehaviors(param.Pipelines))
}

func NewPublishContainer(param PublisherParams) mediator.PublishContainer {
	return mediator.NewPublishContainer(
		mediator.WithNotificationDefinitionHandlers(param.NotificationHandlers...),
		mediator.WithPublishStrategy(param.PublishStrategy),
	)
}

func NewPublisher(container mediator.PublishContainer) mediator.Publisher {
	return mediator.NewPublisher(container)
}

func NewSender(container mediator.SendContainer) mediator.Sender {
	return mediator.NewSender(container)
}

// NewModule returns a new fx.Option that provides the mediator components
// A PublishStrategy can be provided by the user, otherwise a synchronous strategy is used
func NewModule() fx.Option {
	return fx.Module("mediatorfx",
		fx.Provide(
			NewSendContainer,
			NewSender,
			NewPublishContainer,
			NewPublisher,
		),
	)
}
