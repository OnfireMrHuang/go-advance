package conf

type Bootstrap struct {
	Db       DB       `ini:"db"`
	LocalDir LocalDir `ini:"localDir"`
	Oss      OSS      `ini:"oss"`
}

type DB struct {
	Database string `ini:"database"`
	Host     string `ini:"host"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Port     int    `ini:"port"`
}

type LocalDir struct {
	Csv string `ini:"csv"`
	Zip string `ini:"zip"`
}

type OSS struct {
	Endpoint     string `ini:"endpoint"`
	AccessKey    string `ini:"access_key"`
	AccessSecret string `ini:"access_secret"`
	Bucket       string `ini:"bucket"`
}
