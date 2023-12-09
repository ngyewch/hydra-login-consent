package basic

type Config struct {
	Name               string `koanf:"name"`
	BackgroundImageUri string `koanf:"backgroundImageUri"`
	ForgotPasswordText string `koanf:"forgotPasswordText"`
	ForgotPasswordUri  string `koanf:"forgotPasswordUri"`
}
