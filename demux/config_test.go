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

func TestConfigFailsOnNilCfg(t *testing.T) {
	cfg, err := newDemuxCfg(nil)
	assert.EqualError(t, err, "missing demuxer configuration")
	assert.Nil(t, cfg)
}

func TestConfigFailsOnMissingFactories(t *testing.T) {
	externalCfg := &Cfg{
		Handlers: []any{
			func(ctx context.Context, event *events.APIGatewayProxyRequest) (
				*events.APIGatewayProxyResponse,
				error) {
				return nil, nil
			},
		},
	}

	cfg, err := newDemuxCfg(externalCfg)
	assert.EqualError(t, err, "no demux factories provided")
	assert.Nil(t, cfg)
}

func TestConfigFailsOnMissingHandlers(t *testing.T) {
	externalCfg := &Cfg{
		Factories: []Factory{
			func(ctx *EventContext) any {
				return nil
			},
		},
	}

	cfg, err := newDemuxCfg(externalCfg)
	assert.EqualError(t, err, "no demux handlers provided")
	assert.Nil(t, cfg)
}

func TestConfigFailsOnBadHandler(t *testing.T) {
	externalCfg := &Cfg{
		Handlers: []any{
			func(ctx context.Context, event error) (
				*events.APIGatewayProxyResponse,
				error) {
				return nil, nil
			},
		},
		Factories: []Factory{
			func(ctx *EventContext) any {
				return nil
			},
		},
	}

	cfg, err := newDemuxCfg(externalCfg)
	assert.EqualError(
		t,
		err,
		"lambda demux handler function: expected second argument of handler function to be pointer to event struct; got func(context.Context, error) (*events.APIGatewayProxyResponse, error)")
	assert.Nil(t, cfg)
}
