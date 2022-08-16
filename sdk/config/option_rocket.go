package config

type RocketOptions struct {
	Urls              []string `yaml:"urls" json:"urls"`
	GroupName         string   `yaml:"group_name" json:"group_name"`
	MaxReconsumeTimes int32    `yaml:"max_reconsume_times" json:"max_reconsume_times"`
	RetryTimes        int      `yaml:"retry_times" json:"retry_times"`
}

func (e RocketOptions) GetRocketOptions() (*RocketOptions, error) {

	return nil, nil
}
