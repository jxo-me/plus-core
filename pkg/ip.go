package pkg

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"net"
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"
)

type LocationResp struct {
	Status      string  `json:"status"` // success
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

// GetLocation 获取外网ip地址
func GetLocation(ctx context.Context, ip string) string {
	if ip == "127.0.0.1" || ip == "localhost" {
		return "内部IP"
	}
	url := fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN", ip)
	c := g.Client().SetHeader("User-Agent", UserAgent)

	var res LocationResp
	err := c.PostVar(ctx, url).Scan(&res)
	if err != nil {
		glog.Errorf(ctx, "post remote api error: %s\n", err.Error())
		return ""
	}
	glog.Debug(ctx, "rpc request ip location result:", res)
	return fmt.Sprintf("%s-%s-%s-%s", res.Country, res.City, res.RegionName, res.Isp)
}

// GetLocalHost 获取局域网ip地址
func GetLocalHost() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}

	}
	return ""
}
