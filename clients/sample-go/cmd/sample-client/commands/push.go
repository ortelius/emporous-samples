package commands

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	managerapi "github.com/uor-framework/uor-client-go/api/services/collectionmanager/v1alpha1"
	"google.golang.org/protobuf/types/known/structpb"
)

// PushOptions describe configuration options that can
// be set using the push subcommand.
type PushOptions struct {
	*RootOptions
	RootDir     string
	Destination string
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

	absRootDir, err := filepath.Abs(o.RootDir)
	if err != nil {
		return err
	}

	req := managerapi.Publish_Request{
		Source:      absRootDir,
		Destination: o.Destination,
		Collection:  &managerapi.Collection{},
	}

	// Add sample client specific attributes
	sampleClientAttributes := map[string]map[string]interface{}{
		"*.jpg": {
			"image": true,
		},
		"*.json": {
			"metadata": true,
		},
	}

	for file, attr := range sampleClientAttributes {
		for k, v := range attr {
			fmt.Fprintf(o.IOStreams.Out, "Adding attributes %v=%v to file pattern %s\n", k, v, file)
		}

	}

	for file, attr := range sampleClientAttributes {
		a, err := structpb.NewStruct(attr)
		if err != nil {
			return err
		}
		f := &managerapi.File{
			File:       file,
			Attributes: a,
		}
		req.Collection.Files = append(req.Collection.Files, f)
	}

	resp, err := client.PublishContent(ctx, &req)
	if err != nil {
		return err
	}

	fmt.Fprintln(o.IOStreams.Out, resp.Digest)

	return nil
}
