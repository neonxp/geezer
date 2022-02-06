package geezer

import "context"

type Params struct {
	Ctx      context.Context
	Path     []string
	Query    Values
	Headers  Values
	Provider string
}
