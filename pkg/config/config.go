package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

const (
	HOSTNAME       = "hostname"
	ETCD_ENDPOINTS = "etcd.endpoints"
	ETCD_TIMEOUT   = "etcd.timeout"
)

func Init(configFile string) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("failed to get hostname")
		hostname = "scheduler0"
	}

	viper.SetDefault(HOSTNAME, hostname)
	viper.SetDefault(ETCD_ENDPOINTS, []string{"etcd"})
	viper.SetDefault(ETCD_TIMEOUT, "5s")
}

func GetHostname() string {
	return viper.GetString(HOSTNAME)
}

func GetEtcdEndpoints() []string {
	return viper.GetStringSlice(ETCD_ENDPOINTS)
}

func GetEtcdTimeout() time.Duration {
	return viper.GetDuration(ETCD_TIMEOUT)
}
