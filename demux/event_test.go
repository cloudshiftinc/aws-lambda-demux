package demux

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var handlers = []any{
	func(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest) (
		*events.APIGatewayProxyResponse,
		error) {
		return nil, nil
	},
	func(ctx context.Context, event *events.APIGatewayProxyRequest) (
		*events.APIGatewayProxyResponse,
		error) {
		return nil, nil
	},
}

func TestCreateEventWorks(t *testing.T) {
	eventFactories := []Factory{
		func(ctx *EventContext) any {
			return &events.APIGatewayProxyRequest{}
		},
	}

	event := createEvent(nil, eventFactories)
	assert.NotNil(t, event)
	assert.Equal(t, reflect.TypeOf((*events.APIGatewayProxyRequest)(nil)).Elem(), reflect.TypeOf(event).Elem())

	event = createEvent(map[string]any{}, eventFactories)
	assert.NotNil(t, event)
	assert.Equal(t, reflect.TypeOf((*events.APIGatewayProxyRequest)(nil)).Elem(), reflect.TypeOf(event).Elem())
}

func TestCreateEventNoFactory(t *testing.T) {
	eventFactories := []Factory{
		func(ctx *EventContext) any {
			return nil
		},
	}

	event := createEvent(nil, eventFactories)
	assert.Nil(t, event)
}

func TestProcessEventWorks(t *testing.T) {
	var handlerCalled = false
	handlerMap, err := createHandlerMap(
		[]any{
			func(ctx context.Context, event *events.APIGatewayProxyRequest) (
				*events.APIGatewayProxyResponse,
				error) {
				handlerCalled = true
				return &events.APIGatewayProxyResponse{}, nil
			}})
	assert.NoError(t, err)
	assert.NotNil(t, handlerMap)

	var factoryCalled = false
	cfg := &demuxCfg{
		factories: []Factory{
			func(ctx *EventContext) any {
				factoryCalled = true
				return &events.APIGatewayProxyRequest{}
			},
		},
		handlerMap: handlerMap,
	}

	rawEvent := map[string]any{}

	resp, err := processEvent(cfg, context.TODO(), rawEvent)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, handlerCalled)
	assert.True(t, factoryCalled)
}
