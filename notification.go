package mediatorfx

import (
	"github.com/Oleexo/mediator-go"
	"go.uber.org/fx"
)

func AsNotificationHandler[TNotification mediator.Notification, TNotificationHandler mediator.NotificationHandler[TNotification]](f any) []interface{} {
	return []interface{}{
		f,
		fx.Annotate(mediator.NewNotificationHandlerDefinition[TNotification],
			fx.As(new(mediator.NotificationHandlerDefinition)),
			fx.From(new(TNotificationHandler)),
			fx.ResultTags(`group:"mediator_notification_handlers"`),
		),
	}
}
