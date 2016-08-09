package conf

import "time"

type confType struct {
	Log struct {
		Debug    bool   `default:"true"`
		JSON     bool   `default:"false"`
		ToFile   bool   `default:"false"`
		Filename string `default:"/var/log/psmbtcp.log"`
	}
	Redis struct {
		Server      string        `default:"127.0.0.1"`
		Port        string        `default:"6379"`
		MaxIdel     int           `default:"3"`
		MaxActive   int           `default:"0"`
		IdelTimeout time.Duration `default:"30s"`
	}
	Mongo struct {
		Server            string        `default:"127.0.0.1"`
		Port              string        `default:"27017"`
		IsDrop            bool          `default:"true"`
		ConnectionTimeout time.Duration `default:"60s"`
		DbName            string        `default:"test"`
		Authentication    bool          `default:"false"`
		Username          string        `default:"username"`
		Password          string        `default:"password"`
	}
	MgoHistory struct {
		DbName         string `default:"test"`
		CollectionName string `default:"mbtcp:history"`
	}
	RedisHistory struct {
		HashName   string `default:"mbtcp:latest"`
		ZsetPrefix string `default:"mbtcp:data:"`
	}
	RedisWriter struct {
		HashName string `default:"mbtcp:writer"`
	}
	Psmbtcp struct {
		DefaultPort          string `default:"502"`
		MinConnectionTimeout int64  `default:"200000"`
		MinPollInterval      int    `default:"1"`
	}
	Zmq struct {
		Pub struct {
			Upstream   string `default:"ipc:///tmp/from.psmb"`
			Downstream string `default:"ipc:///tmp/to.modbus"`
		}
		Sub struct {
			Upstream   string `default:"ipc:///tmp/to.psmb"`
			Downstream string `default:"ipc:///tmp/from.modbus"`
		}
	}
}
