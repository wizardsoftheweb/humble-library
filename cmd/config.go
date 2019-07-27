package wotwhb

func BootstrapConfig() {
	ensureDirectoryExists(configDirectory)
	ensureDirectoryExists(downloadDirectory)
}
