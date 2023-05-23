package main

import (
	"fmt"
	"os"

	"github.com/weyg0/hyperactivity-disorder/pkg/scheduler"
	"k8s.io/component-base/logs"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func main() {

	// 创建一个调度器命令，可以通过参数进行配置
	command := app.NewSchedulerCommand(
		app.WithPlugin(scheduler.Name, scheduler.New),
	)

	// 初始化日志记录
	logs.InitLogs()
	defer logs.FlushLogs()

	// 执行命令
	if err := command.Execute(); err != nil {
		// 如果执行命令过程中出现错误，将错误信息输出到标准错误输出，并退出程序
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}
