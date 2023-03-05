package demux

import (
	"context"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

type Factory func(ctx *EventContext) any

type EventContext struct {
	event           RawEvent
	resourceContext RawResourceContext
}

func processEvent(
	cfg *demuxCfg,
	ctx context.Context,
	rawEvent RawEvent) (any, error) {
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

	// map raw event onto event
	decoderCfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   eventType,
		TagName:  "json",
	}
	decoder, err := mapstructure.NewDecoder(decoderCfg)
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(rawEvent)
	if err != nil {
		return nil, err
	}

	// dispatch event to handler function
	args := make([]reflect.Value, 2)
	args[0] = reflect.ValueOf(ctx)
	args[1] = reflect.ValueOf(event)
	returnValues := handlerFn.Call(args)

	return returnValues[0].Interface(), returnValues[1].Interface().(error)
}

func createEvent(rawEvent RawEvent, eventFactories []Factory) any {
	if rawEvent == nil {
		rawEvent = RawEvent{}
	}

	rawResourceContext, ok := rawEvent["resourceContext"].(RawResourceContext)
	if !ok {
		rawResourceContext = RawResourceContext{}
	}

	ctx := &EventContext{
		event:           rawEvent,
		resourceContext: rawResourceContext,
	}

	for _, eventFactory := range eventFactories {
		event := eventFactory(ctx)
		if event != nil {
			return event
		}
	}
	return nil
}
