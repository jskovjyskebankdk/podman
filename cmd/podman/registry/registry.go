package registry

import (
	"context"

	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/containers/libpod/pkg/domain/infra"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// DefaultAPIAddress is the default address of the REST socket
const DefaultAPIAddress = "unix:/run/podman/podman.sock"

// DefaultVarlinkAddress is the default address of the varlink socket
const DefaultVarlinkAddress = "unix:/run/podman/io.podman"

type CliCommand struct {
	Mode    []entities.EngineMode
	Command *cobra.Command
	Parent  *cobra.Command
}

const ExecErrorCodeGeneric = 125

var (
	cliCtx          context.Context
	containerEngine entities.ContainerEngine
	exitCode        = ExecErrorCodeGeneric
	imageEngine     entities.ImageEngine

	// Commands holds the cobra.Commands to present to the user, including
	// parent if not a child of "root"
	Commands []CliCommand
)

func SetExitCode(code int) {
	exitCode = code
}

func GetExitCode() int {
	return exitCode
}

func ImageEngine() entities.ImageEngine {
	return imageEngine
}

// NewImageEngine is a wrapper for building an ImageEngine to be used for PreRunE functions
func NewImageEngine(cmd *cobra.Command, args []string) (entities.ImageEngine, error) {
	if imageEngine == nil {
		podmanOptions.FlagSet = cmd.Flags()
		engine, err := infra.NewImageEngine(&podmanOptions)
		if err != nil {
			return nil, err
		}
		imageEngine = engine
	}
	return imageEngine, nil
}

func ContainerEngine() entities.ContainerEngine {
	return containerEngine
}

// NewContainerEngine is a wrapper for building an ContainerEngine to be used for PreRunE functions
func NewContainerEngine(cmd *cobra.Command, args []string) (entities.ContainerEngine, error) {
	if containerEngine == nil {
		podmanOptions.FlagSet = cmd.Flags()
		engine, err := infra.NewContainerEngine(&podmanOptions)
		if err != nil {
			return nil, err
		}
		containerEngine = engine
	}
	return containerEngine, nil
}

func SubCommandExists(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return errors.Errorf("unrecognized command `%[1]s %[2]s`\nTry '%[1]s --help' for more information.", cmd.CommandPath(), args[0])
	}
	return errors.Errorf("missing command '%[1]s COMMAND'\nTry '%[1]s --help' for more information.", cmd.CommandPath())
}

// IdOrLatestArgs used to validate a nameOrId was provided or the "--latest" flag
func IdOrLatestArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 1 || (len(args) == 0 && !cmd.Flag("latest").Changed) {
		return errors.New(`command requires a name, id  or the "--latest" flag`)
	}
	return nil
}

type PodmanOptionsKey struct{}

func Context() context.Context {
	if cliCtx == nil {
		cliCtx = ContextWithOptions(context.Background())
	}
	return cliCtx
}

func ContextWithOptions(ctx context.Context) context.Context {
	cliCtx = context.WithValue(ctx, PodmanOptionsKey{}, podmanOptions)
	return cliCtx
}

// GetContextWithOptions deprecated, use  NewContextWithOptions()
func GetContextWithOptions() context.Context {
	return ContextWithOptions(context.Background())
}

// GetContext deprecated, use  Context()
func GetContext() context.Context {
	return Context()
}