package config

type PullerConfig struct {
	Vendors          []string
	WindowSize       int
	Concurrency      int
	RetryCount       int
	TimeoutPerWindow int
}

type MySQLConfig struct {
	DSN string
}

type MQConfig struct {
	DSN string
}

type FullConfig struct {
	Puller PullerConfig
	MySQL  MySQLConfig
	MQ     MQConfig
}
