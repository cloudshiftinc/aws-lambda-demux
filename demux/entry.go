package demux

import (
	"context"
)

func NewHandler(cfg *Cfg) any {
	demuxCfg, err := newDemuxCfg(cfg)
	if err != nil {
		return errorHandler(err)
	}

	return func(ctx context.Context, rawEvent map[string]any) (any, error) {
		return processEvent(demuxCfg, ctx, rawEvent)
	}
}

func errorHandler(err error) func(context.Context, any) (any, error) {
	return func(_ context.Context, _ any) (any, error) {
		return nil, err
	}
}
