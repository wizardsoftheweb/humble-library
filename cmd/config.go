package wotwhb

func BootstrapConfig(configPath, downloadPath string) {
	ensureDirectoryExists(configPath)
	ensureDirectoryExists(downloadPath)
}
