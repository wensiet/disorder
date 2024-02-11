package config

import (
	"file-chunker/internal/service/files/encryptor"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"sync"
)

var once sync.Once

type config struct {
	Database struct {
		DSN string `yaml:"dsn" env-required:"true"`
	} `yaml:"database"`
	Discord struct {
		Token   string `yaml:"token" env-required:"true"`
		Channel string `yaml:"channel" env-required:"true"`
	} `yaml:"discord"`
	Bucket struct {
		Size int64 `yaml:"size" env-required:"true"`
	}
}

var conf config

func mustLoad() {
	err := cleanenv.ReadConfig("config/config.yaml", &conf)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat("keys/aes"); err != nil {
		gen, err := encryptor.KeyGen(16)
		if err != nil {
			panic(err)
		}
		err = encryptor.SaveKeyToFile(gen, "keys/aes")
	}
}

func GetConfig() config {
	once.Do(mustLoad)
	return conf
}

func GetKey() ([]byte, error) {
	once.Do(mustLoad)
	return encryptor.LoadKeyFromFile("keys/aes")
}
