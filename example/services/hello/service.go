package hello

import (
	"context"
	"log"

	"github.com/neonxp/geezer"
	"github.com/neonxp/geezer/render"
)

const ServiceName = "hello"

type Service struct {
}

func (s Service) Find(ctx context.Context, params geezer.Params) (render.Renderer, error) {
	//TODO implement me
	log.Printf("Find params=%+v", params)
	return render.Text("text/html", "Hello, world!"), nil
}

func (s Service) Get(ctx context.Context, id string, params geezer.Params) (render.Renderer, error) {
	//TODO implement me
	log.Printf("Get id=%s params=%+v", id, params)
	return nil, geezer.ErrMethodNotFound
}

func (s Service) Create(ctx context.Context, data geezer.Data, params geezer.Params) (render.Renderer, error) {
	//TODO implement me
	log.Printf("Create data=%s params=%+v", data, params)
	return nil, geezer.ErrMethodNotFound
}

func (s Service) Update(ctx context.Context, id string, data geezer.Data, params geezer.Params) (render.Renderer, error) {
	//TODO implement me
	log.Printf("Update id=%s data=%s params=%+v", id, data, params)
	return nil, geezer.ErrMethodNotFound
}

func (s Service) Patch(ctx context.Context, id string, data geezer.Data, params geezer.Params) (render.Renderer, error) {
	//TODO implement me
	log.Printf("Patch id=%s data=%s params=%+v", id, data, params)
	return nil, geezer.ErrMethodNotFound
}

func (s Service) Remove(ctx context.Context, id string, params geezer.Params) error {
	//TODO implement me
	log.Printf("Remove id=%s params=%+v", id, params)
	return geezer.ErrMethodNotFound
}

func (s Service) Setup(app geezer.AppKernel, path string) error {
	//TODO implement me
	log.Printf("Setup path=%s", path)
	return nil
}
