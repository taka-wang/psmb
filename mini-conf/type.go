package conf

type confType struct {
	Log struct {
		Debug    bool   `default:"true"`
		JSON     bool   `default:"false"`
		ToFile   bool   `default:"false"`
		Filename string `default:"/var/log/psmbtcp.log"`
	}
	Redis struct {
		Server      string `default:"127.0.0.1"`
		Port        string `default:"6379"`
		MaxIdel     int    `default:"3"`
		MaxActive   int    `default:"0"`
		IdelTimeout int    `default:"30"` // time.Duration
	}
	Mongo struct {
		Server            string `default:"127.0.0.1"`
		Port              string `default:"27017"`
		IsDrop            bool   `default:"true"`
		ConnectionTimeout int    `default:"60"` // time.Duration
		DbName            string `default:"test"`
		Authentication    bool   `default:"false"`
		Username          string `default:"username"`
		Password          string `default:"password"`
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
	RedisFilter struct {
		HashName    string `default:"mbtcp:filter"`
		MaxCapacity int    `default:"32"`
	}
	MemFilter struct {
		MaxCapacity int `default:"32"`
	}
	MemReader struct {
		MaxCapacity int `default:"32"`
	}
	Psmbtcp struct {
		DefaultPort          string `default:"502"`
		MinConnectionTimeout int64  `default:"200000"`
		MinPollInterval      int    `default:"1"`
		MaxWorker            int    `default:"6"`
		MaxQueue             int    `default:"100"`
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
