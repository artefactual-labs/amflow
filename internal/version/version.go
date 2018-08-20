package version

var version string

func Get() string {
	if version == "" {
		version = "(dev)"
	}
	return version
}
