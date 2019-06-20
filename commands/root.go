package commands

import (
	"context"
	"os"
	"strings"

	"github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/cihub/seelog"
	"github.com/spf13/cobra"
)

const (
	// CommandStart cli command "start"
	CommandStart = "start"

	// CommandMock cli command "mock"
	CommandMock = "mock"

	// CommandEvents cli command "events"
	CommandEvents = "events"

	// CommandTx cli command "tx"
	CommandTx = "tx"

	// CommandReset cli command "reset"
	CommandReset = "reset"

	// CommandVersion cli command "version"
	CommandVersion = "version"

	// CommandHelp cli command "help"
	CommandHelp = "help"
)

const (

	// DefaultEventSubscribe events 默认订阅条件
	DefaultEventSubscribe string = "tm.event='Tx' AND qcp.to='qos'"
)

// Runner 通过配置数据执行方法，返回运行过程中出现的错误，如果返回空则代表运行成功。
type Runner func(conf *config.Config) (context.CancelFunc, error)

// NewRootCommand 创建 root/默认 命令
//
// 实现默认功能，显示帮助信息，预处理配置初始化，日志配置初始化。
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "cassini",
		Short: "the relay of cross-chain",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if strings.EqualFold(cmd.Use, CommandVersion) ||
				strings.HasPrefix(cmd.Use, CommandHelp) {
				// doesn't need init log and config
				return nil
			}
			// 初始化日志
			var logger seelog.LoggerInterface
			logger, err = log.LoadLogger(config.GetConfig().LogConfigFile)
			if err != nil {
				log.Warn("Used the default logger because error: ", err)
			} else {
				log.Replace(logger)
			}
			if strings.EqualFold(cmd.Use, CommandEvents) {
				// doesn't need init config
				return nil
			}
			err = initConfig()
			if err != nil {
				return err
			}
			return
		},
	}
	return root
}

func initConfig() error {
	// init config
	err := config.GetConfig().Load()
	if err != nil {
		log.Error("Init config error: ", err.Error())
		return err
	}
	log.Debug("Init config: ", config.GetConfig().ConfigFile)
	return nil
}

func commandRunner(run Runner, isKeepRunning bool) error {
	cancel, err := run(config.GetConfig())
	if err != nil {
		log.Error("Run command error: ", err.Error())
		return err
	}
	if isKeepRunning {
		common.KeepRunning(func(sig os.Signal) {
			defer log.Flush()
			if cancel != nil {
				cancel()
			}
			log.Debug("Stopped by signal: ", sig)
		})
	}
	return nil
}

func reconfigMock(node string) (mock *config.MockConfig) {
	conf := config.GetConfig()
	if len(conf.Mocks) < 1 {
		mock = &config.MockConfig{
			RPC: &config.RPCConfig{
				NodeAddress: node}}
		conf.Mocks = []*config.MockConfig{mock}
	}
	if mock == nil {
		conf.Mocks = conf.Mocks[:1]
		mock = conf.Mocks[0]
		mock.RPC.NodeAddress = node
	}
	return
}
