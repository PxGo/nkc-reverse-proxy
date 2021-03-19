package test

import (
	"reflect"
	"testing"

	"github.com/tokisakiyuu/nkc-proxy-go-pure/pkg/config"
)

func TestParse(t *testing.T) {
	// 预期获得的实例
	want := config.Profile{
		Ports: config.Ports{
			HttpPort:  80,
			HttpsPort: 443,
		},
		Timeout: 10000,
		Servers: []config.Server{
			{
				Name:       "故园怀旧",
				Type:       "forward",
				Hostname:   []string{"localhost", "192.168.11.36"},
				Https:      true,
				HttpTarget: []string{"http://127.0.0.1:9000"},
				WsTarget:   []string{"http://127.0.0.1:2170", "http://127.0.0.1:2171"},
				// generate by preprocess
				NoResponsePage: "path/to/html",
				// generate by preprocess
				SSL: config.SSL{
					Cert: "path/to/cert",
					Key:  "path/to/key",
				},
			},
		},
		NoResponsePage: "path/to/html",
		SSL: config.SSL{
			Cert: "path/to/cert",
			Key:  "path/to/key",
		},
	}

	// 测试
	t.Run("测试配置文件json解析", func(t *testing.T) {
		if conf := config.Parse("proxy.config_test.json"); !reflect.DeepEqual(conf, want) {
			t.Errorf("Ops! \n%v\n%v", conf, want)
		}
	})
}
