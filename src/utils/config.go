package utils

import (
	"flag"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

func readYamlFile(cfg interface{}, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func LoadConfig(cfg interface{}) error {
	cfgfile := flag.String("cfg", "", "config file")
	flag.Parse()
	err := readYamlFile(cfg, *cfgfile)
	if err != nil {
		log.Println("fail to parse the config file:", err)
		return err
	}
	return nil
}
