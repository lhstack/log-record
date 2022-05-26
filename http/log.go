package http

import (
	bytes2 "bytes"
	"fmt"
	"io"
	"log-record/store"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func queryRemoteAllApp(res http.ResponseWriter, req *http.Request) {
	var applications []string
	store.DB.Model(&store.RemoteLog{}).Select("application").Group("application").Find(&applications)
	writeJson(&applications, &res)
}

func queryRemoteAllLevel(res http.ResponseWriter, req *http.Request) {
	parse, err := url.Parse(req.RequestURI)
	if err != nil {
		writeError(err, &res)
		return
	}
	var level []string
	db := store.DB.Model(&store.RemoteLog{}).Select("level").Group("level")
	query := parse.Query()
	if query.Has("app") {
		app := getOrDefault(query, "app", "default", func(str string) (string, error) { return str, nil })
		db = db.Where("application = ?", app)
	}

	if query.Has("ip") {
		app := getOrDefault(query, "ip", "default", func(str string) (string, error) { return str, nil })
		db = db.Where("address = ?", app)
	}

	db.Find(&level)
	writeJson(&level, &res)
}

func queryRemoteAllIp(res http.ResponseWriter, req *http.Request) {
	parse, err := url.Parse(req.RequestURI)
	if err != nil {
		writeError(err, &res)
		return
	}
	var addresses []string
	db := store.DB.Model(&store.RemoteLog{}).Select("address").Group("address")
	query := parse.Query()
	if query.Has("app") {
		app := getOrDefault(query, "app", "default", func(str string) (string, error) { return str, nil })
		db = db.Where("application = ?", app)
	}
	db.Find(&addresses)
	writeJson(&addresses, &res)
}

func queryRemoteLogs(res http.ResponseWriter, req *http.Request) {
	parse, err := url.Parse(req.RequestURI)
	if err != nil {
		writeError(err, &res)
		return
	}
	query := parse.Query()
	page := getOrDefault(query, "page", 1, func(str string) (int, error) {
		return strconv.Atoi(str)
	})
	size := getOrDefault(query, "size", 10, func(str string) (int, error) {
		return strconv.Atoi(str)
	})
	if size <= 0 {
		size = 10
	}
	if page <= 0 {
		page = 1
	}
	var remoteLogs []store.RemoteLog
	db := store.DB.Offset(page - 1).Limit(size)
	if query.Has("startTime") && query.Has("endTime") {
		startTime := getOrDefault(query, "startTime", 0, func(str string) (int, error) {
			return strconv.Atoi(str)
		})
		endTime := getOrDefault(query, "endTime", 0, func(str string) (int, error) {
			return strconv.Atoi(str)
		})
		db = db.Where("timestamp BETWEEN ? AND ?", startTime, endTime)
	} else if query.Has("startTime") && !query.Has("endTime") {
		startTime := getOrDefault(query, "startTime", 0, func(str string) (int, error) {
			return strconv.Atoi(str)
		})
		db = db.Where("timestamp >= ?", startTime)
	} else if !query.Has("startTime") && query.Has("endTime") {
		endTime := getOrDefault(query, "endTime", 0, func(str string) (int, error) {
			return strconv.Atoi(str)
		})
		db = db.Where("timestamp <= ?", endTime)
	}
	if query.Has("app") {
		app := getOrDefault(query, "app", "default", func(str string) (string, error) { return str, nil })
		db = db.Where("application = ?", app)
	}

	if query.Has("ip") {
		ip := getOrDefault(query, "ip", "0.0.0.0", func(str string) (string, error) { return str, nil })
		db = db.Where("address = ?", ip)
	}

	if query.Has("level") {
		level := getOrDefault(query, "level", "INFO", func(str string) (string, error) { return str, nil })
		db = db.Where("level = ?", level)
	}
	var count int64
	db.Order("id ASC").Find(&remoteLogs)
	buffer := bytes2.NewBufferString("")
	for _, log := range remoteLogs {
		message := fmt.Sprintf(`<span style='font-size: 12px'>%v [%v]|-<span style='color: #358b99'>%v</span> <span style='color: green'>%v</span>-<span style='color: red'>%v</span></span><br/><br/>`,
			time.UnixMilli(log.Timestamp).Format("2006-01-02 15:04:05"),
			log.Thread,
			log.Level,
			log.LoggerName,
			log.Message)
		buffer.Write([]byte(message))
	}
	db.Count(&count)
	res.Header().Set("total", fmt.Sprintf("%v", count))
	res.Header().Set("Content-Type", "text/html;charset=UTF-8")
	bytes, err := io.ReadAll(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _ = res.Write(bytes)
}

func queryRemoteLogList(res http.ResponseWriter, req *http.Request) {
	parse, err := url.Parse(req.RequestURI)
	if err != nil {
		writeError(err, &res)
		return
	}
	query := parse.Query()
	page := getOrDefault(query, "page", 1, func(str string) (int, error) {
		return strconv.Atoi(str)
	})
	size := getOrDefault(query, "size", 10, func(str string) (int, error) {
		return strconv.Atoi(str)
	})
	if size <= 0 {
		size = 10
	}
	if page <= 0 {
		page = 1
	}
	var remoteLogs []store.RemoteLog
	db := store.DB.Offset(page - 1).Limit(size)
	if query.Has("startTime") && query.Has("endTime") {
		startTime := getOrDefault(query, "startTime", 0, func(str string) (int, error) {
			return strconv.Atoi(str)
		})
		endTime := getOrDefault(query, "endTime", 0, func(str string) (int, error) {
			return strconv.Atoi(str)
		})
		db = db.Where("timestamp BETWEEN ? AND ?", startTime, endTime)
	} else if query.Has("startTime") && !query.Has("endTime") {
		startTime := getOrDefault(query, "startTime", 0, func(str string) (int, error) {
			return strconv.Atoi(str)
		})
		db = db.Where("timestamp >= ?", startTime)
	} else if !query.Has("startTime") && query.Has("endTime") {
		endTime := getOrDefault(query, "endTime", 0, func(str string) (int, error) {
			return strconv.Atoi(str)
		})
		db = db.Where("timestamp <= ?", endTime)
	}
	db.Order("timestamp DESC").Find(&remoteLogs)
	writeJson(remoteLogs, &res)
}
