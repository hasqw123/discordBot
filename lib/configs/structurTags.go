package configs

type TgBot struct {
	Host  string `yaml:"tgBotHost"`
	Token string `yaml:"token"`
}

type Config struct {
	BatchSize     int    `yaml:"batchSize"`
	AmountHandler int    `yaml:"amountHandler"`
	TgBot         TgBot  `yaml:"tgBot"`
	DscBot        DscBot `yaml:"dscBot"`
}

type DscBot struct {
	Token string `yaml:"token"`
}

func New() Config {
	return Config{}
}
