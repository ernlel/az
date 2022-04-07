package az

import (
	"bytes"
	"log"
	"strings"
	"text/template"
)

type doc struct {
	name        string
	description string
	docParams   []*docParam
	docHeaders  []*docHeader
	docCookies  []*docCookie
}

type docParam struct {
	Name        string
	ParamType   string
	Description string
}

type docHeader struct {
	Name        string
	Description string
}
type docCookie struct {
	Name        string
	Description string
}

type documentation struct {
	Name        string
	Description string
	Namespace   map[string][]documentationDoc
}

type documentationDoc struct {
	Name           string
	Description    string
	Namespace      string
	Path           string
	Method         string
	RequiredParams []docParam
	DocParams      []docParam
	DocHeaders     []docHeader
	DocCookies     []docCookie
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (h *HandlerStruct) Doc(name string, description string) *HandlerStruct {
	d := new(doc)
	d.name = name
	d.description = description
	h.doc = d
	return h
}

func (h *HandlerStruct) Param(name string, paramType string, description string) *HandlerStruct {
	dp := new(docParam)
	dp.Name = name
	dp.ParamType = paramType
	dp.Description = description
	h.doc.docParams = append(h.doc.docParams, dp)
	return h
}

func (h *HandlerStruct) Header(name string, description string) *HandlerStruct {
	he := new(docHeader)
	he.Name = name
	he.Description = description
	h.doc.docHeaders = append(h.doc.docHeaders, he)
	return h
}
func (h *HandlerStruct) Cookie(name string, description string) *HandlerStruct {
	co := new(docCookie)
	co.Name = name
	co.Description = description
	h.doc.docCookies = append(h.doc.docCookies, co)
	return h
}

func (r *Router) Documentation(name string, description string) []byte {

	d := documentation{}
	d.Name = name
	d.Description = description
	d.Namespace = make(map[string][]documentationDoc)

	for _, ro := range r.routes {
		dd := documentationDoc{}
		// dd.RequiredParams = make(map[string]string)
		// dd.DocParams = make(map[string]string)
		// dd.DocHeaders = make(map[string]string)
		dd.Namespace = ro.namespace
		dd.Path = ro.path
		if ro.HandlerStruct != nil {
			dd.DocHeaders = []docHeader{}
			dd.DocCookies = []docCookie{}
			dd.DocParams = []docParam{}
			dd.RequiredParams = []docParam{}
			for _, p := range ro.HandlerStruct.doc.docParams {
				dd.DocParams = append(dd.DocParams, *p)
			}
			for _, h := range ro.HandlerStruct.doc.docHeaders {
				dd.DocHeaders = append(dd.DocHeaders, *h)
			}
			for _, h := range ro.HandlerStruct.doc.docCookies {
				dd.DocCookies = append(dd.DocCookies, *h)
			}
			dd.Name = ro.HandlerStruct.doc.name
			dd.Description = ro.HandlerStruct.doc.description
			d.Namespace[ro.namespace] = append(d.Namespace[ro.namespace], dd)
		}
		for methodName, MethodStruct := range ro.methods {
			if MethodStruct.HandlerStruct != nil {
				dd.DocHeaders = []docHeader{}
				dd.DocCookies = []docCookie{}
				dd.DocParams = []docParam{}
				dd.RequiredParams = []docParam{}
				for _, p := range MethodStruct.HandlerStruct.doc.docParams {
					dd.DocParams = append(dd.DocParams, *p)
				}
				for _, h := range MethodStruct.HandlerStruct.doc.docHeaders {
					dd.DocHeaders = append(dd.DocHeaders, *h)
				}
				for _, h := range MethodStruct.HandlerStruct.doc.docCookies {
					dd.DocCookies = append(dd.DocCookies, *h)
				}
				dd.Method = methodName
				dd.Name = MethodStruct.HandlerStruct.doc.name
				dd.Description = MethodStruct.HandlerStruct.doc.description
				d.Namespace[ro.namespace] = append(d.Namespace[ro.namespace], dd)
			}
			for _, ParamStruct := range MethodStruct.params {
				if ParamStruct.HandlerStruct != nil {
					dd.DocHeaders = []docHeader{}
					dd.DocCookies = []docCookie{}
					dd.DocParams = []docParam{}
					dd.RequiredParams = []docParam{}

					for _, p := range ParamStruct.doc.docParams {
						if stringInSlice(p.Name, ParamStruct.requiredParams) {
							dd.RequiredParams = append(dd.RequiredParams, *p)
						} else {
							dd.DocParams = append(dd.DocParams, *p)
						}
					}

					for _, rp := range ParamStruct.requiredParams {
						found := false
						for _, rpp := range dd.RequiredParams {
							if rpp.Name == rp {
								found = true
								break
							}
						}
						if !found {
							dd.RequiredParams = append(dd.RequiredParams, docParam{Name: rp})
						}
					}

					for _, h := range ParamStruct.doc.docHeaders {
						dd.DocHeaders = append(dd.DocHeaders, *h)
					}
					for _, h := range ParamStruct.doc.docCookies {
						dd.DocCookies = append(dd.DocCookies, *h)
					}
					dd.Method = methodName
					dd.Name = ParamStruct.doc.name
					dd.Description = ParamStruct.doc.description
					d.Namespace[ro.namespace] = append(d.Namespace[ro.namespace], dd)
				}
			}
		}

	}

	// log.Println(d)

	var documentationHTML bytes.Buffer

	funcMap := template.FuncMap{
		"replace": func(input, from, to string) string {
			return strings.Replace(input, from, to, -1)
		},
		"HTMLbr": func(s string) string {
			return strings.Replace(s, "\n", "<br>", -1)
		},
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}

	t, err := template.New("documentation").Funcs(funcMap).Parse(documentationTemplate)
	if err != nil {
		log.Println(err)
		return []byte("Error parsing documentation template. " + err.Error())
	}

	err = t.Execute(&documentationHTML, d)
	if err != nil {
		log.Println(err)
		return []byte("Error executing documentation template. " + err.Error())
	}

	return documentationHTML.Bytes()
}
