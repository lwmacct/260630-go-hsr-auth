package config

import (
	"time"

	"github.com/lwmacct/260614-go-pkg-tlsreload/pkg/tlsreload"
)

type Config struct {
	Server Server `json:"server" desc:"服务端配置"`
}

type Server struct {
	Debug    bool           `json:"debug"    desc:"启用调试日志和诊断信息"`
	Database ServerDatabase `json:"database" desc:"数据库配置"`
	Auth     ServerAuth     `json:"auth"     desc:"认证配置"`
	HTTP     ServerHTTP     `json:"http"     desc:"HTTP 服务配置"`
}

type ServerDatabase struct {
	Type   string              `json:"type"   desc:"数据库类型：sqlite、pgsql"`
	SQLite string              `json:"sqlite" desc:"SQLite 数据库文件路径"`
	PGSQL  ServerDatabasePGSQL `json:"pgsql"  desc:"PostgreSQL 连接参数"`
}

type ServerDatabasePGSQL struct {
	Host     string `json:"host"     desc:"PostgreSQL 主机"`
	Port     string `json:"port"     desc:"PostgreSQL 端口"`
	User     string `json:"user"     desc:"PostgreSQL 用户名"`
	Database string `json:"database" desc:"PostgreSQL 数据库名"`
	Password string `json:"password" desc:"PostgreSQL 密码"`
}

type ServerAuth struct {
	Admins    []string            `json:"admins"    desc:"运行时管理员用户名列表"`
	Local     ServerAuthLocal     `json:"local"     desc:"本地账号认证配置"`
	OAuth     ServerAuthOAuth     `json:"oauth"     desc:"第三方登录配置"`
	Challenge ServerAuthChallenge `json:"challenge" desc:"认证挑战配置"`
}

type ServerAuthLocal struct {
	LoginEnabled        bool `json:"login-enabled"        desc:"是否启用本地账号登录"`
	RegistrationEnabled bool `json:"registration-enabled" desc:"是否允许用户名密码公开注册"`
}

type ServerAuthOAuth struct {
	Enabled         bool                    `json:"enabled"           desc:"是否启用 OAuth 第三方登录"`
	AutoRegister    bool                    `json:"auto-register"     desc:"OAuth 首次登录是否自动创建用户"`
	CallbackBaseURL string                  `json:"callback-base-url" desc:"OAuth 回调外部基准地址，留空则根据当前请求推断"`
	GitHub          ServerAuthOAuthProvider `json:"github"            desc:"GitHub OAuth 配置"`
	Google          ServerAuthOAuthProvider `json:"google"            desc:"Google OAuth 配置"`
}

type ServerAuthOAuthProvider struct {
	Enabled      bool     `json:"enabled"       desc:"是否启用该 OAuth 提供方"`
	ClientID     string   `json:"client-id"     desc:"OAuth Client ID"`
	ClientSecret string   `json:"client-secret" desc:"OAuth Client Secret"`
	Scopes       []string `json:"scopes"        desc:"OAuth 授权范围"`
	AuthURL      string   `json:"auth-url"      desc:"OAuth 授权地址"`
	TokenURL     string   `json:"token-url"     desc:"OAuth Token 地址"`
	UserInfoURL  string   `json:"userinfo-url"  desc:"OAuth 用户信息地址"`
}

type ServerAuthChallenge struct {
	Provider  string                    `json:"provider"  desc:"认证挑战提供方：image、hcaptcha、turnstile"`
	Image     ServerAuthChallengeImage  `json:"image"     desc:"图片验证码配置"`
	HCaptcha  ServerAuthChallengeRemote `json:"hcaptcha"  desc:"hCaptcha 挑战配置"`
	Turnstile ServerAuthChallengeRemote `json:"turnstile" desc:"Cloudflare Turnstile 挑战配置"`
}

type ServerAuthChallengeImage struct {
	MaxChallenges int `json:"max-challenges" desc:"内存中图片验证码最大数量，0 表示不限制"`
}

type ServerAuthChallengeRemote struct {
	SiteKey   string `json:"sitekey"    desc:"认证挑战站点公钥"`
	Secret    string `json:"secret"     desc:"认证挑战服务端密钥"`
	VerifyURL string `json:"verify-url" desc:"认证挑战服务端验证地址"`
}

type ServerHTTP struct {
	Listen          string           `json:"listen"             desc:"HTTP 服务监听地址"`
	WebRoot         string           `json:"web-root"           desc:"静态 Web 根目录，留空则不托管前端"`
	TLS             tlsreload.Config `json:"tls"                desc:"HTTPS TLS 配置"`
	SessionTTL      time.Duration    `json:"session-ttl"        desc:"HTTP 登录会话有效期"`
	TrustedProxies  []string         `json:"trusted-proxies"    desc:"可信 HTTP 反向代理 CIDR/IP 列表，仅这些来源可提供真实客户端 IP 头"`
	ReadTimeout     time.Duration    `json:"read-timeout"       desc:"HTTP 读取超时时间"`
	WriteTimeout    time.Duration    `json:"write-timeout"      desc:"HTTP 写入超时时间"`
	IdleTimeout     time.Duration    `json:"idle-timeout"       desc:"HTTP 空闲连接超时时间"`
	MaxAPIBodyBytes int64            `json:"max-api-body-bytes" desc:"HTTP API 最大请求体字节数，0 表示不限制"`
}

func DefaultConfig() Config {
	return Config{
		Server: Server{
			Database: ServerDatabase{
				Type:   "sqlite",
				SQLite: "${APP_DATA:-.local/data}/sqlite.db",
				PGSQL: ServerDatabasePGSQL{
					Host:     "${PGHOST}",
					Port:     "${PGPORT}",
					User:     "${PGUSER}",
					Database: "${PGDATABASE}",
					Password: "${PGPASSWORD}",
				},
			},
			Auth: ServerAuth{
				Admins: []string{"admin"},
				Local: ServerAuthLocal{
					LoginEnabled:        true,
					RegistrationEnabled: true,
				},
				OAuth: ServerAuthOAuth{
					Enabled:      false,
					AutoRegister: true,
					GitHub: ServerAuthOAuthProvider{ // #nosec G101 - provider endpoint defaults, not credentials.
						Scopes:      []string{"read:user", "user:email"},
						AuthURL:     "https://github.com/login/oauth/authorize",
						TokenURL:    "https://github.com/login/oauth/access_token",
						UserInfoURL: "https://api.github.com/user",
					},
					Google: ServerAuthOAuthProvider{ // #nosec G101 - provider endpoint defaults, not credentials.
						Scopes:      []string{"openid", "email", "profile"},
						AuthURL:     "https://accounts.google.com/o/oauth2/v2/auth",
						TokenURL:    "https://oauth2.googleapis.com/token",
						UserInfoURL: "https://openidconnect.googleapis.com/v1/userinfo",
					},
				},
				Challenge: ServerAuthChallenge{
					Provider: "image",
					Image: ServerAuthChallengeImage{
						MaxChallenges: 1024,
					},
					HCaptcha: ServerAuthChallengeRemote{
						VerifyURL: "https://api.hcaptcha.com/siteverify",
					},
					Turnstile: ServerAuthChallengeRemote{
						VerifyURL: "https://challenges.cloudflare.com/turnstile/v0/siteverify",
					},
				},
			},
			HTTP: ServerHTTP{
				Listen:  ":40318",
				WebRoot: "${WEB_ROOT:-dist}",
				TLS: tlsreload.Config{
					Enabled:      false,
					CertFile:     "${APP_DATA:-.local/data}/ssl/fullchain.pem",
					KeyFile:      "${APP_DATA:-.local/data}/ssl/privkey.pem",
					PollInterval: 3 * time.Second,
				},
				SessionTTL:      7 * 24 * time.Hour,
				ReadTimeout:     30 * time.Second,
				WriteTimeout:    30 * time.Second,
				IdleTimeout:     120 * time.Second,
				MaxAPIBodyBytes: 1 << 20,
			},
		},
	}
}
