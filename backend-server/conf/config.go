package conf

var EnvironmentConfig Config

//var DatabaseConnect Database[gorm.DB]

func init() {
	EnvironmentConfig = NewEnvConfig()
	//DatabaseConnect = NewDataBase(EnvironmentConfig)
}
