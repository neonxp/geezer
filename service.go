package geezer

import (
	"encoding/json"

	"github.com/neonxp/geezer/render"
)

type Service interface {
	Find(params Params) (render.Renderer, error)
	Get(id string, params Params) (render.Renderer, error)
	Create(data Data, params Params) (render.Renderer, error)
	Update(id string, data Data, params Params) (render.Renderer, error)
	Patch(id string, data Data, params Params) (render.Renderer, error)
	Remove(id string, params Params) error
	Setup(app Kernel, path string) error
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
