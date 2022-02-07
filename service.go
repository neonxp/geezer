package geezer

import (
	"context"
	"encoding/json"

	"github.com/neonxp/geezer/render"
)

type Service interface {
	Find(ctx context.Context, params Params) (render.Renderer, error)
	Get(ctx context.Context, id string, params Params) (render.Renderer, error)
	Create(ctx context.Context, data Data, params Params) (render.Renderer, error)
	Update(ctx context.Context, id string, data Data, params Params) (render.Renderer, error)
	Patch(ctx context.Context, id string, data Data, params Params) (render.Renderer, error)
	Remove(ctx context.Context, id string, params Params) error
	Setup(app AppKernel, path string) error
}

type Method int

const (
	MethodFind Method = iota
	MethodGet
	MethodCreate
	MethodUpdate
	MethodPatch
	MethodRemove
)

type Data json.RawMessage
