package main

import (
	"fmt"
	"time"

	"github.com/Owen-Zhang/zsf/common/flag"

	"github.com/Owen-Zhang/zsf/config"
)

//App 应用struct
type App struct {
}

//LoadConfig 加载配制文件
func (a App) LoadConfig() {
	enablWatch := flag.Bool("watch")
	config.Init(enablWatch)
}

func main() {
	fmt.Println(string(config.Get("mysql.yaml")))
	for {
		time.Sleep(10 * time.Second)
	}
}
