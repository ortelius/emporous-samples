package commands

import (
	"context"
	"net"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	managerapi "github.com/uor-framework/uor-client-go/api/services/collectionmanager/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type RootOptions struct {
	IOStreams     genericclioptions.IOStreams
	ServerAddress string
}

// NewRootCmd creates a new cobra.Command for the command root.
func NewRootCmd() *cobra.Command {
	o := RootOptions{}

	o.IOStreams = genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	cmd := &cobra.Command{
		Use:           filepath.Base(os.Args[0]),
		Short:         "Sample Client",
		SilenceErrors: false,
		SilenceUsage:  false,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVarP(&o.ServerAddress, "socket-address", "s", "/var/run/uor.sock", "location of unix domain socket")

	cmd.AddCommand(NewPullCmd(&o))
	cmd.AddCommand(NewPushCmd(&o))

	return cmd
}

// clientSetup creates a new CollectionManagerClient instance from given inputs.
func clientSetup(ctx context.Context, serverAddress string) (managerapi.CollectionManagerClient, func() error, error) {
	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(unixConnect))
	if err != nil {
		return nil, nil, err
	}

	client := managerapi.NewCollectionManagerClient(conn)
	cleanup := func() error {
		return conn.Close()
	}
	return client, cleanup, nil
}

// unixConnect creates a unix address from a given input.
func unixConnect(_ context.Context, addr string) (net.Conn, error) {
	unixAddr, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		return nil, err
	}
	return net.DialUnix("unix", nil, unixAddr)
}
