package configs

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"log"
)

var Cfg struct {
	JWTKey string
}

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	_, err := toml.DecodeFile("configs/config.toml", &Cfg)
	if err != nil {
		panic(err)
	}
	b, _ := json.MarshalIndent(Cfg, "", "  ")
	log.Println("Cfg", string(b))
}
