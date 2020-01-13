package config

type Config struct {
	Service struct {
		TimeStart string `yaml:"time_start"`
		TimeEnd   string `yaml:"time_end"`
	} `yaml:"service"`
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Store struct {
		Dirpass string `yaml:"dirpass"`
	} `yaml:"store"`
}
