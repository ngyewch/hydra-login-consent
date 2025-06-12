package basic

type Config struct {
	Name               string `koanf:"name" validate:"required"`
	BackgroundImageUri string `koanf:"backgroundImageUri"`
	ForgotPasswordText string `koanf:"forgotPasswordText"`
	ForgotPasswordUri  string `koanf:"forgotPasswordUri"`
}
