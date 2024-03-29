package demux

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

type Factory func(ctx *EventContext) any

type EventContext struct {
	Event           map[string]any
	ResourceContext map[string]any
}

func processEvent(
	cfg *demuxCfg,
	ctx context.Context,
	rawEvent map[string]any) (any, error) {
	event := createEvent(rawEvent, cfg.factories)
	if event == nil {
		return nil, errors.New("unable to determine event type for demux")
	}
	eventType := reflect.TypeOf(event).Elem()
	// TODO - verify that event is *type

	// lookup handler for event type
	handlerFn, ok := cfg.handlerMap[eventType]

	if !ok {
		return nil, fmt.Errorf("unable to find handler function for event type %s", eventType)
	}

	var returnValues []reflect.Value
	err := xray.Capture(
		ctx, fmt.Sprintf("event: %s", eventType), func(ctx context.Context) error {
			// map raw event onto event
			decoderCfg := &mapstructure.DecoderConfig{
				Metadata: nil,
				Result:   event,
				TagName:  "json",
			}
			decoder, err := mapstructure.NewDecoder(decoderCfg)
			if err != nil {
				return err
			}
			err = decoder.Decode(rawEvent)
			if err != nil {
				return err
			}

			// dispatch event to handler function
			args := make([]reflect.Value, 2)
			args[0] = reflect.ValueOf(ctx)
			args[1] = reflect.ValueOf(event)
			returnValues = handlerFn.Call(args)
			return nil
		})

	if err != nil {
		return nil, err
	}

	if errVal, ok := returnValues[1].Interface().(error); ok && errVal != nil {
		return nil, errVal
	}
	return returnValues[0].Interface(), nil
}

func createEvent(rawEvent map[string]any, eventFactories []Factory) any {
	if rawEvent == nil {
		rawEvent = map[string]any{}
	}

	rawResourceContext, ok := rawEvent["resourceContext"].(map[string]any)
	if !ok {
		rawResourceContext = map[string]any{}
	}

	ctx := &EventContext{
		Event:           rawEvent,
		ResourceContext: rawResourceContext,
	}

	for _, eventFactory := range eventFactories {
		event := eventFactory(ctx)
		if event != nil {
			return event
		}
	}
	return nil
}
