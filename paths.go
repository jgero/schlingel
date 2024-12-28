package main

import (
	"context"
	"io"

	"github.com/jgero/schlingel/api"
	"github.com/jgero/schlingel/page"
)

func pages() []api.HtmxPartialRenderer {
	return []api.HtmxPartialRenderer{
		func() (string, func(context.Context, io.Writer, map[string]string) error) {
			return "/", func(ctx context.Context, w io.Writer, params map[string]string) error {
				return page.Index().Render(ctx, w)
			}
		},
	}
}
