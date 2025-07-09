package config

type Config struct {
	General    GeneralConfig    `mapstructure:"general"`
	HTTPServer HTTPServerConfig `mapstructure:"server"`
	Hydra      HydraConfig      `mapstructure:"hydra"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	OAuth      OAuthConfig      `mapstructure:"oauth"`
	Redis      RedisConfig      `mapstructure:"redis"`
}

type GeneralConfig struct {
	FrontendURL string `mapstructure:"frontend_url"`
}

type HTTPServerConfig struct {
	Port    string        `mapstructure:"port"`
	Session SessionConfig `mapstructure:"session"`
}

type SessionConfig struct {
	CookieName string   `mapstructure:"cookie_name"`
	KeyPairs   []string `mapstructure:"key_pairs"`
	Timeout    int      `mapstructure:"timeout"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type HydraConfig struct {
	AdminURL  string `mapstructure:"admin_url"`
	PublicURL string `mapstructure:"public_url"`
}

type OAuthConfig struct {
	Microsoft MicrosoftOAuthConfig `mapstructure:"microsoft"`
}

type MicrosoftOAuthConfig struct {
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	RedirectURL  string   `mapstructure:"redirect_url"`
	Scopes       []string `mapstructure:"scopes"`
}
