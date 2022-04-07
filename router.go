package az

import (
	"log"
	"net/http"
	"strings"
)

type RouteStruct struct {
	path    string
	methods map[string]*MethodStruct
	*HandlerStruct
	middlewares middlewares
	namespace   string
}

type MethodStruct struct {
	MethodStruct string
	params       []*ParamStruct
	*HandlerStruct
	middlewares middlewares
}

type ParamStruct struct {
	requiredParams []string
	*HandlerStruct
	middlewares middlewares
}

type HandlerStruct struct {
	handlerFunc http.HandlerFunc
	middlewares middlewares
	doc         *doc
}

// Router ...
type Router struct {
	routes         map[string]*RouteStruct
	DefaultHandler http.HandlerFunc // default HandlerStruct
	CaseSensitive  bool
	middlewares    middlewares
	namespace      string
}

// New creates a new Router.
func New() *Router {

	newRouter := new(Router)
	newRouter.routes = make(map[string]*RouteStruct)
	newRouter.DefaultHandler = http.NotFound
	newRouter.CaseSensitive = true

	return newRouter
}

func (r *Router) Namespace(namespace string) *Router {

	if namespace[0] != '/' {
		namespace = string('/') + namespace
	}

	if namespace != "/" && namespace[len(namespace)-1:] == "/" {
		namespace = namespace[:len(namespace)-1]
	}

	namespaceRouter := new(Router)

	namespaceRouter.CaseSensitive = r.CaseSensitive
	// namespaceRouter.DefaultHandler = r.DefaultHandler
	// namespaceRouter.middlewares = r.middlewares
	namespaceRouter.routes = r.routes

	namespaceRouter.namespace = namespace

	return namespaceRouter
}

func (*Router) Method(methodName string, handlerOrParam ...interface{}) *MethodStruct {
	me := new(MethodStruct)

	methodName = strings.ToUpper(methodName)

	for _, hOrP := range handlerOrParam {
		switch v := hOrP.(type) {
		case *ParamStruct:
			//Param
			me.params = append(me.params, v)

		case *HandlerStruct:
			//Handler
			me.HandlerStruct = v
		default:
			log.Fatal("Method for \"" + methodName + "\" accepts only ParamStruct or HandlerStruct")
		}
	}

	me.MethodStruct = methodName
	return me
}

func (r *Router) Param(handlerOrParamName ...interface{}) *ParamStruct {
	p := new(ParamStruct)

	for _, hOrPN := range handlerOrParamName {
		switch v := hOrPN.(type) {
		case string:
			//Param
			if !r.CaseSensitive {
				v = strings.ToLower(v)
			}
			p.requiredParams = append(p.requiredParams, v)

		case *HandlerStruct:
			//Handler
			p.HandlerStruct = v
		default:
			log.Fatal("Param accepts only string or HandlerStruct")
		}
	}

	return p
}

func Handler(handlerFunc ...http.HandlerFunc) *HandlerStruct {
	h := new(HandlerStruct)
	h.doc = new(doc)
	if len(handlerFunc) < 1 {
		h.handlerFunc = func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(501)
			w.Write([]byte(`Handler Not Implemented.`))
		}
		return h
	}

	h.handlerFunc = handlerFunc[0]
	return h
}

func (*Router) Handler(handlerFunc ...http.HandlerFunc) *HandlerStruct {
	return Handler(handlerFunc...)
}

func (h *HandlerStruct) HandlerFunc(handlerFunc http.HandlerFunc) *HandlerStruct {
	h.handlerFunc = handlerFunc
	return h
}

// Add RouteStruct.
func (r *Router) Route(path string, handlerOrMethod ...interface{}) *RouteStruct {

	// add slash on beggining if it missing
	if path[0] != '/' {
		path = string('/') + path
	}

	if !r.CaseSensitive {
		path = strings.ToLower(path)
	}

	if r.namespace != "" {
		path = r.namespace + path
	}

	// remove trailing slash
	if path != "/" && path[len(path)-1:] == "/" {
		path = path[:len(path)-1]
	}

	if _, ok := r.routes[path]; ok {
		log.Fatal("Route \"" + path + "\" already registered")
	}

	ro := new(RouteStruct)
	ro.methods = make(map[string]*MethodStruct)

	for _, hOrM := range handlerOrMethod {
		switch v := hOrM.(type) {
		case *MethodStruct:
			//Method
			if _, ok := ro.methods[v.MethodStruct]; ok {
				log.Fatal("Method " + v.MethodStruct + " already registered for RouteStruct: " + path)
			}
			ro.methods[v.MethodStruct] = v

		case *HandlerStruct:
			//Handler
			ro.HandlerStruct = v
		default:
			log.Fatal("Route for \"" + path + "\" accepts only MethodStruct or HandlerStruct")
		}
	}

	ro.path = path
	if r.namespace != "" {
		ro.middlewares = r.middlewares
		ro.namespace = r.namespace
	}
	r.routes[path] = ro

	return ro
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	path := req.URL.Path
	MethodStruct := req.Method

	if !r.CaseSensitive {
		path = strings.ToLower(path)
	}

	// remove trailing slash
	if path != "/" && path[len(path)-1:] == "/" {
		path = path[:len(path)-1]
	}

	if RouteStruct, ok := r.routes[path]; ok {
		if m, ok := RouteStruct.methods[MethodStruct]; ok {

			params := req.URL.Query()

			bestMatch := new(ParamStruct)
			bestMatchTimes := 0

			for _, ParamStruct := range m.params {
				times := 0
				allMatch := true
				for _, reqParam := range ParamStruct.requiredParams {
					if _, ok := params[reqParam]; ok {
						times++
					} else {
						allMatch = false
						break
					}
				}
				if times > 0 && times > bestMatchTimes && allMatch {
					bestMatch = ParamStruct
					bestMatchTimes = times
				}

			}

			if bestMatch.HandlerStruct != nil {
				useMiddlewares(bestMatch.HandlerStruct.handlerFunc, r.middlewares, RouteStruct.middlewares, m.middlewares, bestMatch.middlewares, bestMatch.HandlerStruct.middlewares)(w, req)
				return
			}

			useMiddlewares(m.HandlerStruct.handlerFunc, r.middlewares, RouteStruct.middlewares, m.middlewares, m.HandlerStruct.middlewares)(w, req)
			return
		}
		if RouteStruct.HandlerStruct != nil {
			useMiddlewares(RouteStruct.HandlerStruct.handlerFunc, r.middlewares, RouteStruct.middlewares, RouteStruct.HandlerStruct.middlewares)(w, req)
			return
		}
	}

	useMiddlewares(r.DefaultHandler, r.middlewares)(w, req)

}
