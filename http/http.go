package http

import (
	"encoding/json"
	"fmt"
	"log-record/utils"
	"net/http"
	"net/url"
)

func init() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	listenPort := utils.EnvGetOrDefaultStringValue("HTTP_LISTEN_PORT", "8080")
	http.HandleFunc("/remoteLogList", queryRemoteLogList)
	http.HandleFunc("/remoteLogs", queryRemoteLogs)
	http.HandleFunc("/remoteLogApplications", queryRemoteAllApp)
	http.HandleFunc("/remoteLogAllIp", queryRemoteAllIp)
	http.HandleFunc("/remoteLogAllLevel", queryRemoteAllLevel)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%v", listenPort), nil)
		if err != nil {
			panic(err)
		}
	}()
	fmt.Println(fmt.Sprintf("init http server,HTTP_LISTEN_PORT %v \r\n", listenPort))
}

func writeError(err error, res *http.ResponseWriter) {
	response := *res
	header := response.Header()
	header.Set("content-type", "application/json;charset=utf-8")
	response.WriteHeader(200)
	result := make(map[string]any, 3)
	result["code"] = 400
	result["message"] = err.Error()
	marshal, err := json.Marshal(&result)
	if err != nil {
		_, _ = response.Write([]byte(err.Error()))
	} else {
		_, _ = response.Write(marshal)
	}
}

func writeJson(data interface{}, res *http.ResponseWriter) {
	response := *res
	header := response.Header()
	header.Set("content-type", "application/json;charset=utf-8")
	response.WriteHeader(200)
	result := make(map[string]any, 3)
	result["code"] = 200
	result["message"] = "success"
	result["data"] = data
	marshal, err := json.Marshal(&result)
	if err != nil {
		_, _ = response.Write([]byte(err.Error()))
	} else {
		_, _ = response.Write(marshal)
	}
}

func getOrDefault[T comparable](query url.Values, key string, defaultValue T, apply func(str string) (T, error)) T {
	has := query.Has(key)
	if has {
		value := query.Get(key)
		result, err := apply(value)
		if err != nil {
			fmt.Printf("get key %v,parse %v failure,err %v \r\n", key, value, err.Error())
			return defaultValue
		}
		return result
	}
	return defaultValue
}
