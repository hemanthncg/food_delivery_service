package common

import "testing"

func TestLoadConfig(t *testing.T) {
    cfg := LoadConfig()
    if cfg.DBUrl == "" || cfg.RedisUrl == "" || cfg.KafkaUrl == "" {
        t.Error("expected config values to be set")
    }
}
