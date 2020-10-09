package data

type RedisStruct struct {
	URL        string   `mapstructure:"url"` // deprecated
	Servers    []string `mapstructure:"servers"`
	Cluster    bool     `mapstructure:"cluster"`
	MasterName string   `mapstructure:"master_name"`
	PoolSize   int      `mapstructure:"pool_size"`
	Password   string   `mapstructure:"password"`
	Database   int      `mapstructure:"database"`
}
