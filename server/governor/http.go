package governor

import (
	"encoding/json"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
)

//HandleFunc 注册路由
func HandleFunc(parttern string, handler http.HandlerFunc) {
	DefaultServeMux.HandleFunc(parttern, handler)
	routes = append(routes, parttern)
}

func init() {
	HandleFunc("/test", func(writer http.ResponseWriter, req *http.Request) {
		writer.Write([]byte("test"))
	})
	HandleFunc("/routes", func(writer http.ResponseWriter, req *http.Request) {
		json.NewEncoder(writer).Encode(routes)
	})
	HandleFunc("/debug/pprof/", pprof.Index)
	HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	HandleFunc("/debug/pprof/profile", pprof.Profile)
	HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	HandleFunc("/debug/pprof/trace", pprof.Trace)

	HandleFunc("/allocs", pprof.Handler("allocs").ServeHTTP)
	HandleFunc("/block", pprof.Handler("block").ServeHTTP)
	HandleFunc("/goroutine", pprof.Handler("goroutine").ServeHTTP)
	HandleFunc("/heap", pprof.Handler("heap").ServeHTTP)
	HandleFunc("/mutex", pprof.Handler("mutex").ServeHTTP)
	HandleFunc("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)

	if info, ok := debug.ReadBuildInfo(); ok {
		HandleFunc("/modinfo", func(writer http.ResponseWriter, req *http.Request) {
			json.NewEncoder(writer).Encode(info)
		})
	}
}
