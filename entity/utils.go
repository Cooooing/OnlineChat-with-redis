package entity

import (
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var r *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	r = rand.New(source)
}

// GetRandomColor 获取随机 ANSI value 用于颜色
func GetRandomColor() string {
	color := 15
	for {
		color = r.Intn(232)
		i := (color - 34) % 36
		if i >= 0 && i <= 18 {
			break
		}
	}
	return strconv.Itoa(color)
}

// ClearTerminal 清屏
func ClearTerminal() {
	optSys := runtime.GOOS
	if strings.HasPrefix(optSys, "linux") {
		//执行clear指令清除控制台
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Panicln(err)
		}
	} else if strings.HasPrefix(optSys, "windows") {
		//执行clear指令清除控制台
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Panicln(err)
		}
	}
}

// TimeHandle 时间格式化处理 time->string
func TimeHandle(sendTime time.Time) string {
	if sendTime.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		return sendTime.Format("15:04:05")
	}
	return sendTime.Format("2006-01-02 15:04:05")
}
