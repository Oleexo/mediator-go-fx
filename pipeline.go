package mediatorfx

import (
	"github.com/Oleexo/mediator-go"
	"go.uber.org/fx"
)

func AsPipelineBehavior(f any) interface{} {
	return fx.Annotate(f,
		fx.As(new(mediator.PipelineBehavior)),
		fx.ResultTags(`group:"mediator_pipelines"`),
	)
}
