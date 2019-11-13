package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addEventsFlags(cmd *cobra.Command) {
	cmd.Flags().String("node", "127.0.0.1:26657", "node address")
	cmd.Flags().String("subscribe", "tm.event='Tx' AND qcp.to='qos'", "event subscribe query")
	cmd.Flags().Bool("exporter", false, "export to metrics gauge")
}

// NewEventsCommand 创建 events 命令
func NewEventsCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "Start web socket client and subscribe tx event",
		RunE: func(cmd *cobra.Command, args []string) error {
			mock := reconfigMock(viper.GetString("node"))
			mock.Subscribe = viper.GetString("subscribe")
			return commandRunner(run, isKeepRunning)
		},
	}

	addEventsFlags(cmd)
	return cmd
}
