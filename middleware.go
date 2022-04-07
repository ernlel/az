package az

import "net/http"

type middlewares []middleware

type middlewareFn func(h http.HandlerFunc, params ...interface{}) http.HandlerFunc

type middleware struct {
	middlewareFn middlewareFn
	params       []interface{}
}

func (r *Router) Use(mfn middlewareFn, params ...interface{}) *Router {
	m := middleware{middlewareFn: mfn, params: params}
	r.middlewares = append(r.middlewares, m)
	return r
}

func (h *HandlerStruct) Use(mfn middlewareFn, params ...interface{}) *HandlerStruct {
	m := middleware{middlewareFn: mfn, params: params}
	h.middlewares = append(h.middlewares, m)
	return h
}

func (ro *RouteStruct) Use(mfn middlewareFn, params ...interface{}) *RouteStruct {
	m := middleware{middlewareFn: mfn, params: params}
	ro.middlewares = append(ro.middlewares, m)
	return ro
}

func (me *MethodStruct) Use(mfn middlewareFn, params ...interface{}) *MethodStruct {
	m := middleware{middlewareFn: mfn, params: params}
	me.middlewares = append(me.middlewares, m)
	return me
}

func (p *ParamStruct) Use(mfn middlewareFn, params ...interface{}) *ParamStruct {
	m := middleware{middlewareFn: mfn, params: params}
	p.middlewares = append(p.middlewares, m)
	return p
}

func useMiddlewares(h http.HandlerFunc, m ...middlewares) http.HandlerFunc {

	for i := len(m) - 1; i >= 0; i-- {
		for ii := len(m[i]) - 1; ii >= 0; ii-- {
			h = m[i][ii].middlewareFn(h, m[i][ii].params...)
		}
	}

	return h
}

// Default Middlewares
var Middlewares = struct {
	Logger    middlewareFn
	BasicAuth middlewareFn
	CORS      middlewareFn
}{
	Logger:    logger,
	BasicAuth: basicAuth,
	CORS:      CORS,
}
