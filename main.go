package main

import (
  _ "embed"
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/app"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
  "golang.org/x/sys/windows"
  "io/ioutil"
  "os/exec"
  "go-redis-panel/llog"
  "strconv"
)

//go:embed logo.png
var iconData []byte // 嵌入的图标数据
// 定义一个全局或适当作用域内的channel用于通信
var stopChan chan bool = make(chan bool, 1)
var messageChan chan string = make(chan string, 1)

func main() {
  myApp := app.New()
  myWindow := myApp.NewWindow("Go Redis Panel")
  myWindow.Resize(fyne.NewSize(600, 400)) // 设置窗口大小

  // 创建图标资源并设置
  iconResource := fyne.NewStaticResource("logo.png", iconData)
  myApp.SetIcon(iconResource)    // 设置应用图标
  myWindow.SetIcon(iconResource) // 设置窗口图标，如果你想要不同于应用图标的窗口图标

  // 输入框和输出框
  redisPathEntry := widget.NewEntry()
  redisPathEntry.SetPlaceHolder("Enter Redis installation path here")
  redisPathEntry.Text = loadRedisPath()

  outputEntry := widget.NewMultiLineEntry()
  outputEntry.SetPlaceHolder("Output will appear here")
  outputEntry.Wrapping = fyne.TextWrapWord
  go func() {
    for msg := range messageChan {
      outputEntry.SetText(msg)
    }
  }()
  // 创建按钮并绑定事件
  startButton := widget.NewButton("Start Redis", func() {
    redisPath := redisPathEntry.Text
    saveRedisPath(redisPath)
    go startRedisServer(redisPath)
  })
  stopButton := widget.NewButton("Stop Redis", func() {
    stopChan <- true // 发送停止信号
  })
  restartButton := widget.NewButton("Restart Redis", func() {
    stopChan <- true // 发送停止信号
    go startRedisServer(redisPathEntry.Text)
  })

  // 将控件添加到容器
  buttons := container.NewHBox(startButton, stopButton, restartButton) // 使用HBox代替VBox
  inputForm := container.NewVBox(widget.NewForm(widget.NewFormItem("Redis Path", redisPathEntry)))
  content := container.NewVBox(inputForm, buttons, outputEntry)

  myWindow.SetContent(content)
  myWindow.ShowAndRun()
}

func startRedisServer(redisPath string) {
  var cmd = getCommand(redisPath, redisPath+"\\redis-server", "--daemonize", "yes")
  err := cmd.Start() // 启动但不等待其完成
  if err != nil {
    llog.Log.Error("error:", err.Error())
    messageChan <- "error:" + err.Error()
    return
  } else {
    messageChan <- "running pid:" + strconv.Itoa(cmd.Process.Pid)
  }
  done := make(chan error)
  //messageChan <- cmd.ProcessState.String() + " pid:" + strconv.Itoa(cmd.ProcessState.Pid())
  go func() { done <- cmd.Wait() }() // 在另一个goroutine中等待命令完成
  select {
  case <-stopChan:
    //cmd.Cancel()
    err := cmd.Process.Kill() //关键即可,关闭会触发cmd.wait()
    if err != nil {
      messageChan <- "Faild to killed pid:" + strconv.Itoa(cmd.Process.Pid) + " error:"
      err.Error()
    } else {
      messageChan <- "Killed successfully pid :" + strconv.Itoa(cmd.Process.Pid)
    }
    //退出goroutine
    return
  case err := <-done: // 命令执行完毕或发生错误
    if err != nil {
      messageChan <- "error: " + err.Error()
    } else {
      messageChan <- "Redis server started successfully."
    }
    //退出goroutine
    return
  }
}

func getCommand(dir string, command string, args ...string) *exec.Cmd {
  cmd := exec.Command(command, args...)
  if dir != "" {
    cmd.Dir = dir
  }

  // 设置命令执行时不显示窗口
  cmd.SysProcAttr = &windows.SysProcAttr{
    HideWindow: true,
  }

  return cmd
}

// 命令执行函数
func execCommand(dir string, command string, args ...string) (string, error) {
  cmd := getCommand(dir, command, args...)

  outputBytes, err := cmd.CombinedOutput()

  return string(outputBytes), err
}

// 将Redis路径保存到文件
func saveRedisPath(path string) {
  err := ioutil.WriteFile("redis-path.txt", []byte(path), 0644)
  if err != nil {
    llog.Log.Info("Failed to save Redis path:", err)
  }
}

// 读取Redis路径
func loadRedisPath() string {
  bytes, err := ioutil.ReadFile("redis-path.txt")
  if err != nil {
    llog.Log.Info("Failed to load Redis path:", err)
    return ""
  }
  return string(bytes)
}
