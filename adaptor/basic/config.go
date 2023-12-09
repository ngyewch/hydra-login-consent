package basic

type Config struct {
	Name               string `koanf:"name"`
	ForgotPasswordText string `koanf:"forgotPasswordText"`
	ForgotPasswordUri  string `koanf:"forgotPasswordUri"`
}
