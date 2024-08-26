package main

import (
	"strings"
	"text/template"
)

// функция для конверта мапы название тега - его значение, чтобы можно было подставить мапу в шаблон
func MapToStr(tag map[string]string) string {
	var builder strings.Builder
	builder.WriteString("{")
	for key, value := range tag {
		builder.WriteString(`"` + key + `":"` + value + `",`)
	}
	result := builder.String()
	if len(result) > 1 {
		result = result[:len(result)-1]
	}
	result += "}"
	return "map[string]string" + result
}

var (
	tmplFuncs = template.FuncMap{
		"MapToStr": MapToStr,
	}

	serveHttpTemplate = template.Must(template.New("serverHttp").Parse(`
func (h *{{ (index . 0).StructName }} ) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path { 
	{{range $inf := .}}
		case "{{$inf.CommentInf.Url}}":
			h.handler{{$inf.FuncName}}(w, r) 
	{{end}}
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			respJson, _ := json.Marshal(map[string]string{"error": "unknown method"})
				w.Write(respJson)
	}
}

`))

	handlerTemplate = template.Must(template.New("handlerFunc").Parse(
		`func (h *{{.StructName}}) handler{{.FuncName}}(w http.ResponseWriter, r *http.Request) {
		{{$inf := .CommentInf}}
		
		w.Header().Set("Content-Type", "application/json")

		{{if $inf.Auth}} 
			if r.Header.Get("X-Auth") != "100500" {
				w.WriteHeader(http.StatusForbidden)
				respJson, _ := json.Marshal(map[string]string{"error": "unauthorized"})
				w.Write(respJson)
				return
			}
		{{end}}		

		ctx := r.Context()

		var query  url.Values
		if r.Method=="POST"{ 
			err := r.ParseForm()
			if err != nil {
				fmt.Println("error decoding body")
			}
			query = r.Form
		} {{if not $inf.Method}} else if r.Method=="GET" {
			query = r.URL.Query()
		} {{end}} else {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Header().Set("Content-Type", "application/json")
			respJson, _ := json.Marshal(map[string]string{"error": "bad method"})
			w.Write(respJson)
			return
		} 

			
		receivedStruct, err := Get{{.FuncName}}ParamsFromUrlQuery(query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			respJson, _ := json.Marshal(map[string]string{"error": err.Error()})
				w.Write(respJson)
			return
		}

		res, err := h.{{.FuncName}}(ctx, receivedStruct)

		if err != nil {
			var apiErr ApiError
			ok := errors.As(err, &apiErr)
			if ok {
				w.WriteHeader(apiErr.HTTPStatus)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
		}
		respJson, _ := json.Marshal(map[string]interface{}{"error": err.Error()})
		w.Write(respJson)
		return
	}

	w.WriteHeader(http.StatusOK)
	respJson, err := json.Marshal(map[string]interface{} {"error": "", "response": res})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		respJson, _ := json.Marshal(map[string]string{"error": err.Error()})
			w.Write(respJson)
		return
	}
	w.Write(respJson)
}
`))

	fillParamsGet = template.Must(template.New("getStructFromUrlQuery").Funcs(tmplFuncs).Parse(`
// заполнение структуры params
func Get{{.StructName}}FromUrlQuery(urlQuery url.Values) ({{.StructName}}, error) {
	var res {{.StructName}}

	fieldsList := []string{ {{range $field := .Fields}}
		"{{-  $field.Paramname -}}",	
	{{end}} }

	var params []string
	for _, name := range fieldsList {
		if urlQuery.Get(name) != "" {
			params = urlQuery[name]
		} else {
			params = []string{""}
		}
		
		switch name {
		{{range $field := .Fields}}
			case "{{$field.Paramname}}":

				{{$tag := MapToStr $field.ApigenTag}}

				{{if eq $field.FieldType "int"}}
					num, err := IsTagIntValid({{$tag}},params[0], "{{$field.Paramname}}")
					if err != nil {
						return res, err
					}
					res.{{$field.Name}} = num
				{{else}}
					str, err := IsTagStrValid({{$tag}},params[0],"{{$field.Paramname}}")
					if err != nil {
						return res, err
					}
					res.{{$field.Name}} = str
				{{end}}
		{{end}}
		}
	}
	return res, nil
}
`))

	imports = `
import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)
`
)
