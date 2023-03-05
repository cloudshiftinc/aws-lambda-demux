package lambdaDemux

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestHandlerValidationWorks(t *testing.T) {
	fn := func(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest) (
		*events.APIGatewayProxyResponse,
		error) {
		return nil, nil
	}
	ctx, err := processHandler(fn)
	assert.NoError(t, err)
	assert.NotNil(t, ctx)

	// can't directly compare function types; instead, compare string representation
	assert.Equal(t, reflect.TypeOf(fn).String(), ctx.fn.Type().String())

	// ensure event type is as expected
	assert.Equal(t, reflect.TypeOf((*events.APIGatewayWebsocketProxyRequest)(nil)).Elem(), ctx.eventType)
}

func TestHandlerFailsOnTooManyArgs(t *testing.T) {
	ctx, err := processHandler(
		func(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest, foo any) (
			*events.APIGatewayProxyResponse,
			error) {
			return nil, nil
		})
	assert.EqualError(
		t,
		err,
		"lambda demux handler function: expected to take two parameters and return two parameters; got func(context.Context, *events.APIGatewayWebsocketProxyRequest, interface {}) (*events.APIGatewayProxyResponse, error)")
	assert.Nil(t, ctx)
}

func TestHandlerFailsOnTooFewArgs(t *testing.T) {
	ctx, err := processHandler(
		func(event *events.APIGatewayWebsocketProxyRequest) (
			*events.APIGatewayProxyResponse,
			error) {
			return nil, nil
		})
	assert.EqualError(
		t,
		err,
		"lambda demux handler function: expected to take two parameters and return two parameters; got func(*events.APIGatewayWebsocketProxyRequest) (*events.APIGatewayProxyResponse, error)")
	assert.Nil(t, ctx)
}

func TestHandlerFailsOnTooManyReturnValues(t *testing.T) {
	ctx, err := processHandler(
		func(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest) (
			*events.APIGatewayProxyResponse,
			error, int) {
			return nil, nil, -1
		})
	assert.EqualError(
		t,
		err,
		"lambda demux handler function: expected to take two parameters and return two parameters; got func(context.Context, *events.APIGatewayWebsocketProxyRequest) (*events.APIGatewayProxyResponse, error, int)")
	assert.Nil(t, ctx)
}

func TestHandlerFailsOnTooFewReturnValues(t *testing.T) {
	ctx, err := processHandler(
		func(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest) *events.APIGatewayProxyResponse {
			return nil
		})
	assert.EqualError(
		t,
		err,
		"lambda demux handler function: expected to take two parameters and return two parameters; got func(context.Context, *events.APIGatewayWebsocketProxyRequest) *events.APIGatewayProxyResponse")
	assert.Nil(t, ctx)
}

func TestHandlerFailsOnFirstArgNotContext(t *testing.T) {
	ctx, err := processHandler(
		func(ctx *context.Context, event *events.APIGatewayWebsocketProxyRequest) (
			*events.APIGatewayProxyResponse,
			error) {
			return nil, nil
		})
	assert.EqualError(
		t,
		err,
		"lambda demux handler function: expected first argument of handler function to be context.Context; got func(*context.Context, *events.APIGatewayWebsocketProxyRequest) (*events.APIGatewayProxyResponse, error)")
	assert.Nil(t, ctx)
}

func TestHandlerFailsOnSecondArgNotPointer(t *testing.T) {
	ctx, err := processHandler(
		func(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (
			*events.APIGatewayProxyResponse,
			error) {
			return nil, nil
		})
	assert.EqualError(
		t,
		err,
		"lambda demux handler function: expected second argument of handler function to be pointer to event struct; got func(context.Context, events.APIGatewayWebsocketProxyRequest) (*events.APIGatewayProxyResponse, error)")
	assert.Nil(t, ctx)
}

func TestHandlerFailsOnSecondArgNotStruct(t *testing.T) {
	ctx, err := processHandler(
		func(ctx context.Context, err error) (
			*events.APIGatewayProxyResponse,
			error) {
			return nil, nil
		})
	assert.EqualError(
		t,
		err,
		"lambda demux handler function: expected second argument of handler function to be pointer to event struct; got func(context.Context, error) (*events.APIGatewayProxyResponse, error)")
	assert.Nil(t, ctx)
}

func TestHandlerFailsOnFirstReturnValueNotPointer(t *testing.T) {
	ctx, err := processHandler(
		func(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest) (
			events.APIGatewayProxyResponse,
			error) {
			return events.APIGatewayProxyResponse{}, nil
		})
	assert.EqualError(
		t,
		err,
		"lambda demux handler function: expected first return value to be pointer to structure; got func(context.Context, *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error)")
	assert.Nil(t, ctx)
}

func TestHandlerFailsOnSecondReturnValueNotError(t *testing.T) {
	ctx, err := processHandler(
		func(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest) (
			*events.APIGatewayProxyResponse,
			int) {
			return nil, -1
		})
	assert.EqualError(
		t,
		err,
		"lambda demux handler function: expected second return value to be 'error'; got func(context.Context, *events.APIGatewayWebsocketProxyRequest) (*events.APIGatewayProxyResponse, int)")
	assert.Nil(t, ctx)
}

func foo() {
	h := &someHandler{}
	cfg := &Cfg{
		Factories: []Factory{
			func(ctx *EventContext) any {
				return nil
			},
		},
		Handlers: []any{
			func(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest) (
				*events.APIGatewayProxyResponse,
				error) {
				return h.HandleEvent(ctx, event)
			},
		},
	}

	lambda.Start(NewDemuxHandler(cfg))
}

type someHandler struct {
}

func (s *someHandler) HandleEvent(
	_ context.Context,
	_ *events.APIGatewayWebsocketProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return nil, nil
}
