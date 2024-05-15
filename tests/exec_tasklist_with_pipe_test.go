package tests

import (
  "os/exec"
  "redis-panel/llog"
  "testing"
)

func TestExeTaskList(t *testing.T) {
  // 创建cmd.exe进程
  cmd := exec.Command("cmd.exe", "/c", "tasklist | findstr redis-server")

  // 获取命令输出
  output, err := cmd.CombinedOutput()
  if err != nil {
    llog.Log.Error("Command finished with error:", err)
  } else {
    llog.Log.Info(string(output))
  }

}
