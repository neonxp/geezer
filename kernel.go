package geezer

import (
	"errors"
	"strings"

	"github.com/neonxp/geezer/render"
)

var (
	ErrServiceNotFound = errors.New("service not found")
	ErrMethodNotFound  = errors.New("method not found")
)

type defaultKernel struct {
	routes map[string]Service
	hooks  map[string]map[HookLifecycle]map[HookType][]Hook
}

func newKernel() *defaultKernel {
	return &defaultKernel{
		routes: map[string]Service{},
		hooks:  map[string]map[HookLifecycle]map[HookType][]Hook{},
	}
}

func (s *defaultKernel) Register(name string, service Service) error {
	name = strings.ToLower(name)
	s.routes[name] = service
	if _, exist := s.hooks[name]; !exist {
		s.hooks[name] = map[HookLifecycle]map[HookType][]Hook{}
	}
	if err := service.Setup(s, name); err != nil {
		return err
	}
	return nil
}

func (s *defaultKernel) Hook(service string, lifecycle HookLifecycle, hookType HookType, hook Hook) {
	service = strings.ToLower(service)
	if _, exist := s.hooks[service]; !exist {
		s.hooks[service] = map[HookLifecycle]map[HookType][]Hook{}
	}
	if _, exist := s.hooks[service][lifecycle]; !exist {
		s.hooks[service][lifecycle] = map[HookType][]Hook{}
	}
	if _, exist := s.hooks[service][lifecycle][hookType]; !exist {
		s.hooks[service][lifecycle][hookType] = []Hook{}
	}
	s.hooks[service][lifecycle][hookType] = append(s.hooks[service][lifecycle][hookType], hook)
}

func (s *defaultKernel) Service(name string) Service {
	if service, exist := s.routes[name]; exist {
		return service
	}
	return nil
}

func (s *defaultKernel) Call(method Method, name, id string, data Data, params Params) (render.Renderer, error) {
	name = strings.ToLower(name)
	service := s.Service(name)
	if service == nil {
		return nil, ErrServiceNotFound
	}

	hookCtx, result, err := s.callBeforeHooks(method, name, id, data, params)
	if err != nil {
		return result, err
	}
	switch hookCtx.Method {
	case MethodFind:
		result, err = service.Find(hookCtx.Params)
	case MethodGet:
		result, err = service.Get(hookCtx.ID, hookCtx.Params)
	case MethodCreate:
		result, err = service.Create(hookCtx.Data, hookCtx.Params)
	case MethodUpdate:
		result, err = service.Update(hookCtx.ID, hookCtx.Data, hookCtx.Params)
	case MethodPatch:
		result, err = service.Patch(hookCtx.ID, hookCtx.Data, hookCtx.Params)
	case MethodRemove:
		err = service.Remove(hookCtx.ID, hookCtx.Params)
	default:
		return nil, ErrMethodNotFound
	}
	return s.callAfterHooks(method, name, hookCtx, result, err)
}

func (s *defaultKernel) callBeforeHooks(method Method, name string, id string, data Data, params Params) (*HookContext, render.Renderer, error) {
	var beforeHooks []Hook
	if hooks, ok := s.hooks[name][HookBefore]; ok {
		if allHooks, ok := hooks[HookAll]; ok {
			beforeHooks = append(beforeHooks, allHooks...)
		}
		if methodHooks, ok := hooks[hookTypeFromMethod[method]]; ok {
			beforeHooks = append(beforeHooks, methodHooks...)
		}
	}

	hookCtx := &HookContext{
		App:        s,
		Method:     method,
		Type:       HookBefore,
		ID:         id,
		Params:     params,
		Data:       data,
		Err:        nil,
		Result:     nil,
		StatusCode: 0,
	}
	for _, hook := range beforeHooks {
		if err := hook(hookCtx); err != nil {
			return nil, nil, err
		}
	}
	return hookCtx, nil, nil
}

func (s *defaultKernel) callAfterHooks(method Method, name string, hookCtx *HookContext, result render.Renderer, err error) (render.Renderer, error) {
	var afterHooks []Hook
	if hooks, ok := s.hooks[name][HookAfter]; ok {
		if allHooks, ok := hooks[HookAll]; ok {
			afterHooks = append(afterHooks, allHooks...)
		}
		if methodHooks, ok := hooks[hookTypeFromMethod[method]]; ok {
			afterHooks = append(afterHooks, methodHooks...)
		}
	}
	var errorHooks []Hook
	if hooks, ok := s.hooks[name][HookError]; ok {
		if allHooks, ok := hooks[HookAll]; ok {
			errorHooks = append(errorHooks, allHooks...)
		}
		if methodHooks, ok := hooks[hookTypeFromMethod[method]]; ok {
			errorHooks = append(errorHooks, methodHooks...)
		}
	}
	hookCtx.Result = result
	hookCtx.Err = err
	if err != nil {
		for _, hook := range errorHooks {
			if err := hook(hookCtx); err != nil {
				return nil, err
			}
		}
		return hookCtx.Result, hookCtx.Err
	}
	for _, hook := range afterHooks {
		if err := hook(hookCtx); err != nil {
			return nil, err
		}
	}
	return hookCtx.Result, hookCtx.Err
}

type Kernel interface {
	Register(name string, service Service) error
	Hook(service string, lifecycle HookLifecycle, hookType HookType, hook Hook)
	Service(name string) Service
	Call(method Method, name, id string, data Data, params Params) (render.Renderer, error)
}
