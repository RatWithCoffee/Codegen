package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func (h *Storage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	case "/user/profile":
		h.handlerProfile(w, r)

	case "/user/create":
		h.handlerCreate(w, r)

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		respJson, _ := json.Marshal(map[string]string{"error": "unknown method"})
		w.Write(respJson)
	}
}

func (h *Storage) handlerProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	var query url.Values
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("error decoding body")
		}
		query = r.Form
	} else if r.Method == "GET" {
		query = r.URL.Query()
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Header().Set("Content-Type", "application/json")
		respJson, _ := json.Marshal(map[string]string{"error": "bad method"})
		w.Write(respJson)
		return
	}

	receivedStruct, err := GetProfileParamsFromUrlQuery(query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		respJson, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.Write(respJson)
		return
	}

	res, err := h.Profile(ctx, receivedStruct)

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
	respJson, err := json.Marshal(map[string]interface{}{"error": "", "response": res})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		respJson, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.Write(respJson)
		return
	}
	w.Write(respJson)
}

// заполнение структуры params
func GetProfileParamsFromUrlQuery(urlQuery url.Values) (ProfileParams, error) {
	var res ProfileParams

	fieldsList := []string{
		"login",
	}

	var params []string
	for _, name := range fieldsList {
		if urlQuery.Get(name) != "" {
			params = urlQuery[name]
		} else {
			params = []string{""}
		}

		switch name {

		case "login":

			str, err := IsTagStrValid(map[string]string{"required": ""}, params[0], "login")
			if err != nil {
				return res, err
			}
			res.Login = str

		}
	}
	return res, nil
}
func (h *Storage) handlerCreate(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("X-Auth") != "100500" {
		w.WriteHeader(http.StatusForbidden)
		respJson, _ := json.Marshal(map[string]string{"error": "unauthorized"})
		w.Write(respJson)
		return
	}

	ctx := r.Context()

	var query url.Values
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("error decoding body")
		}
		query = r.Form
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Header().Set("Content-Type", "application/json")
		respJson, _ := json.Marshal(map[string]string{"error": "bad method"})
		w.Write(respJson)
		return
	}

	receivedStruct, err := GetCreateParamsFromUrlQuery(query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		respJson, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.Write(respJson)
		return
	}

	res, err := h.Create(ctx, receivedStruct)

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
	respJson, err := json.Marshal(map[string]interface{}{"error": "", "response": res})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		respJson, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.Write(respJson)
		return
	}
	w.Write(respJson)
}

// заполнение структуры params
func GetCreateParamsFromUrlQuery(urlQuery url.Values) (CreateParams, error) {
	var res CreateParams

	fieldsList := []string{
		"login",

		"full_name",

		"status",

		"age",
	}

	var params []string
	for _, name := range fieldsList {
		if urlQuery.Get(name) != "" {
			params = urlQuery[name]
		} else {
			params = []string{""}
		}

		switch name {

		case "login":

			str, err := IsTagStrValid(map[string]string{"required": "", "min": "10"}, params[0], "login")
			if err != nil {
				return res, err
			}
			res.Login = str

		case "full_name":

			str, err := IsTagStrValid(map[string]string{"paramname": "full_name"}, params[0], "full_name")
			if err != nil {
				return res, err
			}
			res.Name = str

		case "status":

			str, err := IsTagStrValid(map[string]string{"enum": "user|moderator|admin", "default": "user"}, params[0], "status")
			if err != nil {
				return res, err
			}
			res.Status = str

		case "age":

			num, err := IsTagIntValid(map[string]string{"min": "0", "max": "128"}, params[0], "age")
			if err != nil {
				return res, err
			}
			res.Age = num

		}
	}
	return res, nil
}
