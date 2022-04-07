package az

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

const toLower = 'a' - 'A'

type CorsParams struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
	Debug            bool
}

func (cp CorsParams) isOriginAllowed(origin string) bool {

	if cp.isOriginAllowedAll() {
		return true
	}

	origin = strings.ToLower(origin)

	for _, o := range cp.AllowedOrigins {
		if strings.ToLower(o) == origin {
			return true
		}
	}
	return false
}

func (cp CorsParams) isOriginAllowedAll() bool {
	for _, o := range cp.AllowedOrigins {
		if o == "*" {
			return true
		}
	}
	return false
}

func (cp CorsParams) isMethodAllowed(MethodStruct string) bool {
	if len(cp.AllowedMethods) == 0 {
		return false
	}

	MethodStruct = strings.ToUpper(MethodStruct)

	if MethodStruct == "OPTIONS" {
		return true
	}

	for _, m := range cp.AllowedMethods {
		if strings.ToUpper(m) == strings.ToUpper(MethodStruct) {
			return true
		}
	}
	return false
}

func (cp CorsParams) areHeadersAllowed(requestedHeaders []string) bool {
	log.Println(len(requestedHeaders))
	if len(requestedHeaders) == 1 && len(requestedHeaders[0]) == 0 {
		return true
	}
	for _, header := range requestedHeaders {
		found := false
		for _, h := range cp.AllowedHeaders {
			if h == header {
				found = true
			}
		}
		if found {
			return true
		}
	}
	return false
}

func handlePreflight(w http.ResponseWriter, r *http.Request, p CorsParams) {

	origin := r.Header.Get("Origin")

	w.Header().Add("Vary", "Origin")
	w.Header().Add("Vary", "Access-Control-Request-Method")
	w.Header().Add("Vary", "Access-Control-Request-Headers")

	if origin == "" {
		if p.Debug {
			log.Printf("CORS > Preflight aborted: empty origin")
		}
		return
	}
	if !p.isOriginAllowed(origin) {
		if p.Debug {
			log.Printf("CORS > Preflight aborted: origin '%s' not allowed", origin)
		}
		return
	}

	reqMethod := r.Header.Get("Access-Control-Request-Method")
	if !p.isMethodAllowed(reqMethod) {
		if p.Debug {
			log.Printf("CORS > Preflight aborted: MethodStruct '%s' not allowed", reqMethod)
		}
		return
	}

	reqHeaders := strings.Split(strings.Replace(r.Header.Get("Access-Control-Request-Headers"), " ", "", -1), ",")
	if !p.areHeadersAllowed(reqHeaders) {
		if p.Debug {
			log.Printf("CORS > Preflight aborted: headers '%v' not allowed", reqHeaders)
		}
		return
	}

	if p.isOriginAllowedAll() && !p.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", strings.ToUpper(reqMethod))
	if len(reqHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(reqHeaders, ", "))
	}
	if p.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	if p.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(p.MaxAge))
	}
	if p.Debug {
		log.Printf("CORS > Preflight response headers: %v", w.Header())
	}

	return
}

func handleCors(w http.ResponseWriter, r *http.Request, p CorsParams) {
	if r.Method == "OPTIONS" {
		if p.Debug {
			log.Printf("CORS > No headers added: MethodStruct == %s", r.Method)
		}
		return
	}
	origin := r.Header.Get("Origin")

	w.Header().Add("Vary", "Origin")

	if origin == "" {
		if p.Debug {
			log.Printf("CORS > No headers added: missing origin")
		}
		return
	}
	if !p.isOriginAllowed(origin) {
		if p.Debug {
			log.Printf("CORS > No headers added: origin '%s' not allowed", origin)
		}
		return
	}

	if !p.isMethodAllowed(r.Method) {
		if p.Debug {
			log.Printf("CORS > No headers added: MethodStruct '%s' not allowed", r.Method)
		}
		return
	}

	if p.isOriginAllowedAll() && !p.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	if len(p.ExposedHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(p.ExposedHeaders, ", "))
	}
	if p.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	if p.Debug {
		log.Printf("CORS > Response added headers: %v", w.Header())
	}
}

func CORS(h http.HandlerFunc, params ...interface{}) http.HandlerFunc {
	var p CorsParams
	if len(params) > 0 {
		if par, ok := params[0].(CorsParams); ok {
			p = par
		} else {
			log.Fatal("wrong cors params")

		}
	} else {
		log.Fatal("not provided cors params")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
			if p.Debug {
				log.Printf("CORS > Preflight request")
			}
			handlePreflight(w, r, p)
			return
		}
		handleCors(w, r, p)
		h(w, r)
	}
}
