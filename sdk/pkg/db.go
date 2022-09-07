package pkg

import "strings"

type Mysql struct {
	Path          string `json:"path" yaml:"path"`
	Config        string `json:"config" yaml:"config"`
	Dbname        string `json:"dbname" yaml:"db-name"`
	Username      string `json:"username" yaml:"username"`
	Password      string `json:"password" yaml:"password"`
	MaxIdleConnes int    `json:"maxIdleConnes" yaml:"max-idle-connes"`
	MaxOpenConnes int    `json:"maxOpenConnes" yaml:"max-open-connes"`
	LogMode       bool   `json:"logMode" yaml:"log-mode"`
	LogZap        string `json:"logZap" yaml:"log-zap"`
}

func GetDBType(link string) string {
	a := strings.Split(link, ":")
	if len(a) > 0 {
		return a[0]
	}
	return ""
}

func Dsn(m *Mysql) string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
}

func GetByLink(link string) Mysql {
	var result Mysql
	a := strings.Split(link, ":")
	if len(a) == 4 {
		result.Username = a[1] // root
		b := strings.Split(a[2], "@tcp(")
		c := strings.Split(a[3], ")/")
		if len(b) == 2 || len(c) == 2 {
			result.Password = b[0]          // gdkid,,..
			result.Path = b[1] + ":" + c[0] // 127.0.0.1:13307
			result.Dbname = c[1]
		}
		result.Config = "charset=utf8mb4&parseTime=True&loc=Local"
		result.LogZap = ""
		result.LogMode = false
		result.MaxIdleConnes = 10
		result.MaxOpenConnes = 100
		return result
	}
	return result
}
