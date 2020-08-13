package normal

import "runtime"

//SystemP 系统参数
type SystemP struct {
	MaxProc int `yaml:"maxproc"` //使用几个Process
}

func defaultConfig() *SystemP {
	return &SystemP{
		MaxProc: runtime.NumCPU(),
	}
}
