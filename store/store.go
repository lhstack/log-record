package store

import (
	"encoding/json"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"log-record/server"
	"net"
	"time"
)

var channelItem = make(chan rxgo.Item)

func init() {
	rxgo.FromChannel(channelItem).BufferWithTimeOrCount(rxgo.WithDuration(time.Second*2), 50).DoOnNext(func(items interface{}) {
		itemArray := items.([]interface{})
		logs := make([]*RemoteLog, len(itemArray))
		for i, item := range itemArray {
			logs[i] = item.(*RemoteLog)
		}
		DB.Begin().CreateInBatches(logs, len(logs)).Commit()
	})
}

// DataStore 存储数据
func DataStore(data []byte, header map[string]string, conn *server.Conn) {
	log := &RemoteLog{}
	err := json.Unmarshal(data, log)
	addr := conn.Conn.RemoteAddr().(*net.TCPAddr)
	log.Address = addr.IP.String()
	if err != nil {
		fmt.Println("data store failure ", string(data), err)
		return
	}
	channelItem <- rxgo.Item{V: log}
}
