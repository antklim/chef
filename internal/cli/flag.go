package cli

// Inspired by https://github.com/auth0/auth0-cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Flag struct {
	LongForm   string
	ShortForm  string
	Help       string
	IsRequired bool
}

func (f *Flag) RegisterString(cmd *cobra.Command, value *string, defaultValue string) {
	registerString(cmd, f, value, defaultValue)
}

func registerString(cmd *cobra.Command, f *Flag, value *string, defaultValue string) {
	cmd.Flags().StringVarP(value, f.LongForm, f.ShortForm, defaultValue, f.Help)

	if err := markFlagRequired(cmd, f); err != nil {
		panic(errors.Wrap(err, "failed to register string flag"))
	}
}

func markFlagRequired(cmd *cobra.Command, f *Flag) error {
	if f.IsRequired {
		return cmd.MarkFlagRequired(f.LongForm)
	}

	return nil
}
