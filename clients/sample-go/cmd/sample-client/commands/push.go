package commands

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	managerapi "github.com/uor-framework/uor-client-go/api/services/collectionmanager/v1alpha1"
)

// PushOptions describe configuration options that can
// be set using the push subcommand.
type PushOptions struct {
	*RootOptions
	RootDir     string
	Destination string
	DSConfig    string
}

// NewPushCmd creates a new cobra.Command for the push subcommand.
func NewPushCmd(rootOpts *RootOptions) *cobra.Command {
	o := PushOptions{RootOptions: rootOpts}

	cmd := &cobra.Command{
		Use:           "push SOCKET-LOCATION SRC DST",
		Short:         "Build and push a UOR collection from a workspace",
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	cmd.Flags().StringVar(&o.DSConfig, "dsconfig", o.DSConfig, "dataset configuration path")

	return cmd
}

func (o *PushOptions) Complete(args []string) error {
	if len(args) < 3 {
		return errors.New("not enough arguments")
	}
	o.ServerAddress = args[0]
	o.RootDir = args[1]
	o.Destination = args[2]
	return nil
}

func (o *PushOptions) Run(ctx context.Context) error {
	client, err := clientSetup(ctx, o.ServerAddress)
	if err != nil {
		return err
	}

	var config []byte
	if o.DSConfig != "" {
		config, err = ioutil.ReadFile(o.DSConfig)
		if err != nil {
			return err
		}
	}

	req := managerapi.Publish_Request{
		Source:      o.RootDir,
		Destination: o.Destination,
		Json:        config,
	}
	resp, err := client.PublishContent(ctx, &req)
	if err != nil {
		return err
	}

	fmt.Fprintln(o.IOStreams.Out, resp.Digest)

	return nil
}
