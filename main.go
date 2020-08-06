package main

import (
	"fmt"
	"time"

	"github.com/Owen-Zhang/zsf/config"
)

//App 应用struct
type App struct {
}

//LoadConfig 加载配制文件
func (a App) LoadConfig() {
	//enablWatch := flag.Bool("watch")
	config.Init()
}

func main() {
	app := &App{}
	app.LoadConfig()

	content, err := config.Get("mysql.yaml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(content))
	for {
		time.Sleep(10 * time.Second)
	}
}
