package hello

import (
	"log"

	"github.com/neonxp/geezer"
)

func RegisterHooks(app geezer.AppKernel) {
	app.Hook(ServiceName, geezer.HookBefore, geezer.HookFind, func(ctx *geezer.HookContext) error {
		log.Printf("Hook before find")
		return nil
	})
	app.Hook(ServiceName, geezer.HookAfter, geezer.HookFind, func(ctx *geezer.HookContext) error {
		log.Printf("Hook after find")
		return nil
	})
}
