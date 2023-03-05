package demux

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigWorks(t *testing.T) {
	externalCfg := &Cfg{
		Factories: []Factory{
			func(ctx *EventContext) any {
				return nil
			},
		},
		Handlers: []any{
			func(ctx context.Context, event *events.APIGatewayProxyRequest) (
				*events.APIGatewayProxyResponse,
				error) {
				return nil, nil
			},
		},
	}

	cfg, err := newDemuxCfg(externalCfg)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}
