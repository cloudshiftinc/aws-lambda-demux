package lambdaDemux

import (
	"errors"
	"reflect"
)

type Cfg struct {
	Factories []Factory
	Handlers  []any
}

type demuxCfg struct {
	handlerMap map[reflect.Type]reflect.Value
	factories  []Factory
}

func newDemuxCfg(cfg *Cfg) (*demuxCfg, error) {
	if cfg == nil {
		return nil, errors.New("missing demuxer configuration")
	}

	if len(cfg.Factories) == 0 {
		return nil, errors.New("no demux factories provided")
	}

	if len(cfg.Handlers) == 0 {
		return nil, errors.New("no demux handlers provided")
	}

	handlerMap, err := createHandlerMap(cfg.Handlers)
	if err != nil {
		return nil, err
	}

	return &demuxCfg{
		factories:  cfg.Factories,
		handlerMap: handlerMap,
	}, nil
}
