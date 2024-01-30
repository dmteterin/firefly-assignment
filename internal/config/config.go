package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ValidatorMinWordLength int    `mapstructure:"VALIDATOR_MIN_WORD_LENGTH"`
	ValidatorRegExp        string `mapstructure:"VALIDATOR_REG_EXP"`

	WordBankFilePath string `mapstructure:"WORD_BANK_FILE_PATH"`

	CrawlerLinkListPath             string `mapstructure:"CRAWLER_LINK_LIST_PATH"`
	CrawlerProxyListPath            string `mapstructure:"CRAWLER_PROXY_LIST_PATH"`
	CrawlerURLRetryLimit            int    `mapstructure:"CRAWLER_URL_RETRY_LIMIT"`
	CrawlerQueueConsumerThreadCount int    `mapstructure:"CRAWLER_QUEUE_CONSUMER_THREAD_COUNT"`
	CrawlerProxyTypePrefix          string `mapstructure:"CRAWLER_PROXY_TYPE_PREFIX"`
	CrawlerEnableExpBackoff         bool   `mapstructure:"CRAWLER_ENABLE_EXP_BACKOFF"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
