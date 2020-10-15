package governor

import (
	"net/http"
	"os"
	"runtime"

	jsoniter "github.com/json-iterator/go"
)

func init() {
	HandleFunc("/configs", func(writer http.ResponseWriter, req *http.Request) {
		writer.Write([]byte("目前还没有实现,待完善中"))
	})
	HandleFunc("/debug/env", func(writer http.ResponseWriter, req *http.Request) {
		writer.WriteHeader(200)
		jsoniter.NewEncoder(writer).Encode(os.Environ())
	})

	HandleFunc("/build/info", func(writer http.ResponseWriter, req *http.Request) {
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknow"
		}
		jsoniter.NewEncoder(writer).Encode(map[string]string{
			"hostname":  hostName,
			"goversion": runtime.Version(),
			"other":     "其它信息待增加",
		})
	})
}
