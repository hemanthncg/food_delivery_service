package common

type Config struct {
    DBUrl    string
    RedisUrl string
    KafkaUrl string
}

// LoadConfig loads configuration values
func LoadConfig() *Config {
    return &Config{
        DBUrl:    "localhost:5432",
        RedisUrl: "localhost:6379",
        KafkaUrl: "localhost:9092",
    }
}
