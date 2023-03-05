package demux

import (
	"context"
	"fmt"
	"reflect"
)

type handlerFunctionContext struct {
	fn        reflect.Value
	eventType reflect.Type
}

func processHandler(handler any) (*handlerFunctionContext, error) {
	handlerType := reflect.TypeOf(handler)

	if handlerType.Kind() != reflect.Func {
		return nil, handlerFunctionError(handlerType, "expected a handler function")
	}

	if handlerType.NumIn() != 2 || handlerType.NumOut() != 2 {
		return nil, handlerFunctionError(
			handlerType,
			"expected to take two parameters and return two parameters")
	}

	argTypes := make([]reflect.Type, 0, 2)
	for i := 0; i < handlerType.NumIn(); i++ {
		argType := handlerType.In(i)
		argTypes = append(argTypes, argType)
	}

	contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if argTypes[0] != contextType {
		return nil, handlerFunctionError(
			handlerType,
			"expected first argument of handler function to be context.Context")
	}
	if argTypes[1].Kind() != reflect.Pointer {
		return nil, handlerFunctionError(
			handlerType,
			"expected second argument of handler function to be pointer to event struct")
	}

	returnTypes := make([]reflect.Type, 0, 2)
	for i := 0; i < handlerType.NumOut(); i++ {
		returnType := handlerType.Out(i)
		returnTypes = append(returnTypes, returnType)
	}

	if returnTypes[0].Kind() != reflect.Pointer {
		return nil, handlerFunctionError(handlerType, "expected first return value to be pointer to structure")
	}

	errorType := reflect.TypeOf((*error)(nil)).Elem()

	if returnTypes[1] != errorType {
		return nil, handlerFunctionError(handlerType, "expected second return value to be 'error'")
	}

	return &handlerFunctionContext{
		fn:        reflect.ValueOf(handler),
		eventType: argTypes[1].Elem(),
	}, nil
}

func createHandlerMap(handlers []any) (map[reflect.Type]reflect.Value, error) {

	handlerMap := map[reflect.Type]reflect.Value{}
	for _, handler := range handlers {
		ctx, err := processHandler(handler)
		if err != nil {
			return nil, err
		}
		if _, ok := handlerMap[ctx.eventType]; ok {
			return nil, fmt.Errorf("event handler for type %s already provided; %s", ctx.eventType, handler)
		}
		handlerMap[ctx.eventType] = ctx.fn
	}
	return handlerMap, nil
}

func handlerFunctionError(handlerType reflect.Type, msg string) error {
	return fmt.Errorf("lambda demux handler function: %s; got %s", msg, handlerType)
}
