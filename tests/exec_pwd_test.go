package tests

import (
  "golang.org/x/sys/windows"
  "os/exec"
  "redis-panel/llog"
  "testing"
)

func TestExecPwd(t *testing.T) {
  //测试卡死,应该是启动了程序,
  cmd := exec.Command("D:\\dev_program\\Redis-x64-3.2.100\\redis-server.exe", "--daemonize", "yes")

  // 设置命令执行时不显示窗口
  cmd.SysProcAttr = &windows.SysProcAttr{
    HideWindow: true,
  }

  outputBytes, err := cmd.CombinedOutput()
  if err != nil {
    llog.Log.Error(err)
  } else {
    llog.Log.Info(string(outputBytes))
  }

}
