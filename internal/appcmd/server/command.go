package server

import (
	"github.com/urfave/cli/v3"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"

	"github.com/lwmacct/260630-go-hsr-auth/internal/config"
)

var (
	defaults = config.DefaultConfig()
	usage    = cfgm.Schema(defaults).Command("server")
)

var Command = &cli.Command{
	Name:            "server",
	Usage:           "start HTTP server",
	Action:          action,
	Commands:        []*cli.Command{version.Command},
	HideHelpCommand: true,
	Flags:           commandFlags(),
}

func commandFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: usage.MustUsage("debug"),
			Value: defaults.Server.Debug,
		},
		&cli.StringFlag{
			Name:  "database.type",
			Usage: usage.MustUsage("database.type"),
			Value: defaults.Server.Database.Type,
		},
		&cli.StringFlag{
			Name:  "database.sqlite",
			Usage: usage.MustUsage("database.sqlite"),
			Value: defaults.Server.Database.SQLite,
		},
		&cli.StringFlag{
			Name:  "database.pgsql.host",
			Usage: usage.MustUsage("database.pgsql.host"),
			Value: defaults.Server.Database.PGSQL.Host,
		},
		&cli.StringFlag{
			Name:  "database.pgsql.port",
			Usage: usage.MustUsage("database.pgsql.port"),
			Value: defaults.Server.Database.PGSQL.Port,
		},
		&cli.StringFlag{
			Name:  "database.pgsql.user",
			Usage: usage.MustUsage("database.pgsql.user"),
			Value: defaults.Server.Database.PGSQL.User,
		},
		&cli.StringFlag{
			Name:  "database.pgsql.database",
			Usage: usage.MustUsage("database.pgsql.database"),
			Value: defaults.Server.Database.PGSQL.Database,
		},
		&cli.StringFlag{
			Name:  "database.pgsql.password",
			Usage: usage.MustUsage("database.pgsql.password"),
			Value: defaults.Server.Database.PGSQL.Password,
		},
		&cli.StringSliceFlag{
			Name:  "auth.admins",
			Usage: usage.MustUsage("auth.admins"),
			Value: defaults.Server.Auth.Admins,
		},
		&cli.BoolFlag{
			Name:  "auth.local.login-enabled",
			Usage: usage.MustUsage("auth.local.login-enabled"),
			Value: defaults.Server.Auth.Local.LoginEnabled,
		},
		&cli.BoolFlag{
			Name:  "auth.local.registration-enabled",
			Usage: usage.MustUsage("auth.local.registration-enabled"),
			Value: defaults.Server.Auth.Local.RegistrationEnabled,
		},
		&cli.BoolFlag{
			Name:  "auth.oauth.enabled",
			Usage: usage.MustUsage("auth.oauth.enabled"),
			Value: defaults.Server.Auth.OAuth.Enabled,
		},
		&cli.BoolFlag{
			Name:  "auth.oauth.auto-register",
			Usage: usage.MustUsage("auth.oauth.auto-register"),
			Value: defaults.Server.Auth.OAuth.AutoRegister,
		},
		&cli.StringFlag{
			Name:  "auth.oauth.callback-base-url",
			Usage: usage.MustUsage("auth.oauth.callback-base-url"),
			Value: defaults.Server.Auth.OAuth.CallbackBaseURL,
		},
		&cli.BoolFlag{
			Name:  "auth.oauth.github.enabled",
			Usage: usage.MustUsage("auth.oauth.github.enabled"),
			Value: defaults.Server.Auth.OAuth.GitHub.Enabled,
		},
		&cli.StringFlag{
			Name:  "auth.oauth.github.client-id",
			Usage: usage.MustUsage("auth.oauth.github.client-id"),
			Value: defaults.Server.Auth.OAuth.GitHub.ClientID,
		},
		&cli.StringFlag{
			Name:  "auth.oauth.github.client-secret",
			Usage: usage.MustUsage("auth.oauth.github.client-secret"),
			Value: defaults.Server.Auth.OAuth.GitHub.ClientSecret,
		},
		&cli.StringSliceFlag{
			Name:  "auth.oauth.github.scopes",
			Usage: usage.MustUsage("auth.oauth.github.scopes"),
			Value: defaults.Server.Auth.OAuth.GitHub.Scopes,
		},
		&cli.StringFlag{
			Name:  "auth.oauth.github.auth-url",
			Usage: usage.MustUsage("auth.oauth.github.auth-url"),
			Value: defaults.Server.Auth.OAuth.GitHub.AuthURL,
		},
		&cli.StringFlag{
			Name:  "auth.oauth.github.token-url",
			Usage: usage.MustUsage("auth.oauth.github.token-url"),
			Value: defaults.Server.Auth.OAuth.GitHub.TokenURL,
		},
		&cli.StringFlag{
			Name:  "auth.oauth.github.userinfo-url",
			Usage: usage.MustUsage("auth.oauth.github.userinfo-url"),
			Value: defaults.Server.Auth.OAuth.GitHub.UserInfoURL,
		},
		&cli.BoolFlag{
			Name:  "auth.oauth.google.enabled",
			Usage: usage.MustUsage("auth.oauth.google.enabled"),
			Value: defaults.Server.Auth.OAuth.Google.Enabled,
		},
		&cli.StringFlag{
			Name:  "auth.oauth.google.client-id",
			Usage: usage.MustUsage("auth.oauth.google.client-id"),
			Value: defaults.Server.Auth.OAuth.Google.ClientID,
		},
		&cli.StringFlag{
			Name:  "auth.oauth.google.client-secret",
			Usage: usage.MustUsage("auth.oauth.google.client-secret"),
			Value: defaults.Server.Auth.OAuth.Google.ClientSecret,
		},
		&cli.StringSliceFlag{
			Name:  "auth.oauth.google.scopes",
			Usage: usage.MustUsage("auth.oauth.google.scopes"),
			Value: defaults.Server.Auth.OAuth.Google.Scopes,
		},
		&cli.StringFlag{
			Name:  "auth.oauth.google.auth-url",
			Usage: usage.MustUsage("auth.oauth.google.auth-url"),
			Value: defaults.Server.Auth.OAuth.Google.AuthURL,
		},
		&cli.StringFlag{
			Name:  "auth.oauth.google.token-url",
			Usage: usage.MustUsage("auth.oauth.google.token-url"),
			Value: defaults.Server.Auth.OAuth.Google.TokenURL,
		},
		&cli.StringFlag{
			Name:  "auth.oauth.google.userinfo-url",
			Usage: usage.MustUsage("auth.oauth.google.userinfo-url"),
			Value: defaults.Server.Auth.OAuth.Google.UserInfoURL,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.provider",
			Usage: usage.MustUsage("auth.challenge.provider"),
			Value: defaults.Server.Auth.Challenge.Provider,
		},
		&cli.IntFlag{
			Name:  "auth.challenge.image.max-challenges",
			Usage: usage.MustUsage("auth.challenge.image.max-challenges"),
			Value: defaults.Server.Auth.Challenge.Image.MaxChallenges,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.hcaptcha.sitekey",
			Usage: usage.MustUsage("auth.challenge.hcaptcha.sitekey"),
			Value: defaults.Server.Auth.Challenge.HCaptcha.SiteKey,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.hcaptcha.secret",
			Usage: usage.MustUsage("auth.challenge.hcaptcha.secret"),
			Value: defaults.Server.Auth.Challenge.HCaptcha.Secret,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.hcaptcha.verify-url",
			Usage: usage.MustUsage("auth.challenge.hcaptcha.verify-url"),
			Value: defaults.Server.Auth.Challenge.HCaptcha.VerifyURL,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.turnstile.sitekey",
			Usage: usage.MustUsage("auth.challenge.turnstile.sitekey"),
			Value: defaults.Server.Auth.Challenge.Turnstile.SiteKey,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.turnstile.secret",
			Usage: usage.MustUsage("auth.challenge.turnstile.secret"),
			Value: defaults.Server.Auth.Challenge.Turnstile.Secret,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.turnstile.verify-url",
			Usage: usage.MustUsage("auth.challenge.turnstile.verify-url"),
			Value: defaults.Server.Auth.Challenge.Turnstile.VerifyURL,
		},
		&cli.StringFlag{
			Name:  "http.listen",
			Usage: usage.MustUsage("http.listen"),
			Value: defaults.Server.HTTP.Listen,
		},
		&cli.StringFlag{
			Name:  "http.web-root",
			Usage: usage.MustUsage("http.web-root"),
			Value: defaults.Server.HTTP.WebRoot,
		},
		&cli.BoolFlag{
			Name:  "http.tls.enabled",
			Usage: usage.MustUsage("http.tls.enabled"),
			Value: defaults.Server.HTTP.TLS.Enabled,
		},
		&cli.StringFlag{
			Name:  "http.tls.cert-file",
			Usage: usage.MustUsage("http.tls.cert-file"),
			Value: defaults.Server.HTTP.TLS.CertFile,
		},
		&cli.StringFlag{
			Name:  "http.tls.key-file",
			Usage: usage.MustUsage("http.tls.key-file"),
			Value: defaults.Server.HTTP.TLS.KeyFile,
		},
		&cli.DurationFlag{
			Name:  "http.tls.poll-interval",
			Usage: usage.MustUsage("http.tls.poll-interval"),
			Value: defaults.Server.HTTP.TLS.PollInterval,
		},
		&cli.DurationFlag{
			Name:  "http.session-ttl",
			Usage: usage.MustUsage("http.session-ttl"),
			Value: defaults.Server.HTTP.SessionTTL,
		},
		&cli.StringSliceFlag{
			Name:  "http.trusted-proxies",
			Usage: usage.MustUsage("http.trusted-proxies"),
			Value: defaults.Server.HTTP.TrustedProxies,
		},
		&cli.DurationFlag{
			Name:  "http.read-timeout",
			Usage: usage.MustUsage("http.read-timeout"),
			Value: defaults.Server.HTTP.ReadTimeout,
		},
		&cli.DurationFlag{
			Name:  "http.write-timeout",
			Usage: usage.MustUsage("http.write-timeout"),
			Value: defaults.Server.HTTP.WriteTimeout,
		},
		&cli.DurationFlag{
			Name:  "http.idle-timeout",
			Usage: usage.MustUsage("http.idle-timeout"),
			Value: defaults.Server.HTTP.IdleTimeout,
		},
		&cli.Int64Flag{
			Name:  "http.max-api-body-bytes",
			Usage: usage.MustUsage("http.max-api-body-bytes"),
			Value: defaults.Server.HTTP.MaxAPIBodyBytes,
		},
	}
}
