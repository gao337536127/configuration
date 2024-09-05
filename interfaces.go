package configuration

type ReadConfig interface {
	GetConfig(section, key, value string) (string, error)
	GetConfigWithEnvironment(section, key, environment, defaultValue string) (string, error)
}
