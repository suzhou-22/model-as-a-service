package property

type ConfigProperty struct {
	Amqp AmqpProperty `mapstructure:"amqp"`
	Gpt  GptProperty  `mapstructure:"gpt"`
}

type AmqpProperty struct {
	Url string `mapstructure:"url"`
}

type GptProperty struct {
	Token   string `mapstructure:"token"`
	BaseUrl string `mapstructure:"baseUrl"`
}
