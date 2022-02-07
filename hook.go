package geezer

import "github.com/neonxp/geezer/render"

//go:generate stringer -type=HookLifecycle,HookType -output hook_string.go
type HookLifecycle int

const (
	HookBefore HookLifecycle = iota
	HookAfter
	HookError
)

type HookType int

const (
	HookAll HookType = iota
	HookFind
	HookGet
	HookCreate
	HookUpdate
	HookPatch
	HookRemove
)

var hookTypeFromMethod = map[Method]HookType{
	MethodFind:   HookFind,
	MethodGet:    HookGet,
	MethodCreate: HookCreate,
	MethodUpdate: HookUpdate,
	MethodPatch:  HookPatch,
	MethodRemove: HookRemove,
}

type HookContext struct {
	App        AppKernel
	Path       []string
	Method     Method
	Type       HookLifecycle
	ID         string
	Params     Params
	Data       Data
	Err        error
	Result     render.Renderer
	StatusCode int
}

type Hook func(ctx *HookContext) error
