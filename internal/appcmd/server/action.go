package server

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"

	"github.com/lwmacct/260630-go-hsr-auth/internal/config"
)

func action(ctx context.Context, cmd *cli.Command) error {
	cfg := cfgm.MustLoadCmd(cmd, config.DefaultConfig(), "")
	return run(ctx, cfg)
}
