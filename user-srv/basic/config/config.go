package config

type AppConfig struct {
	Mysql Mysql `yaml:"Mysql"`
	Redis Redis `yaml:"Redis"`
	Jwt   Jwt   `yaml:"Jwt"`
	Nacos Nacos `yaml:"Nacos"`
}
type Mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}
type Jwt struct {
	SecretKey string `yaml:"secret_key"`
}
type Nacos struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	PublicId  string `yaml:"public_id"`
	NaCosName string `yaml:"nacos_name"`
	Password  string `yaml:"password"`
	DataId    string `yaml:"data_id"`
	Group     string `yaml:"group"`
}
