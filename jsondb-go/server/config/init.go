package config

import "context"

func InitConfig(ctx context.Context, args []string, flags Flags) (*Config, context.Context, error) {

	if err := parseFlags(args, &flags); err != nil {
		return nil, nil, err
	}

	configuration = Config{
		Flags: flags,
		DSN:   ":memory",
		Port:  8080,
	}

	return &configuration, ctx, nil
}
