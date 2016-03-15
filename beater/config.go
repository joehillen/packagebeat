package beater

type PackageConfig struct {
	Period *int64
	Dpkg   *bool
	Rpm    *bool
}

type ConfigSettings struct {
	Input PackageConfig
}
