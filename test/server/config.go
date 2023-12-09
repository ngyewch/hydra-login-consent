package server

type Config struct {
	ListenAddr        string   `koanf:"listenAddr"`
	CsrfAuthKey       string   `koanf:"csrfAuthKey"`
	HydraAdminApiUrls []string `koanf:"hydraAdminApiUrls"`
}
