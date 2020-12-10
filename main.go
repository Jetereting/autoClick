package main

import (
	"encoding/json"
	"fmt"
	"github.com/Jetereting/bmob"
	"github.com/Jetereting/go_util"
	"github.com/go-vgo/robotgo"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Point struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

func ReadMouseLeftClickPoint() Point {
	// 录制单次鼠标点击事件
	defer time.Sleep(time.Millisecond * 300)
	p := Point{}
	if robotgo.AddEvent("mleft") {
		p.X, p.Y = robotgo.GetMousePos()
		fmt.Printf("---你按下左键, 坐标为(%d, %d)---\n", p.X, p.Y)
	}
	return p
}
func WriteMouseLeftClickList(q *bool, mousePointList *[]Point, dur int) {
	// 重放鼠标点击事件队列
	for i := 1; !*q; i++ {
		for _, p := range *mousePointList {
			robotgo.MovesClick(p.X, p.Y, "left", false)
			time.Sleep(time.Millisecond * time.Duration(dur))
		}
		fmt.Println("执行", i, "次")
	}
}
func main() {
	//检查支付
	if !bmob.IsPay("鼠标自动点击") {
		fmt.Println("未付费")
		return
	}

	//间断时间
	dur := 500
	if len(os.Args) > 1 {
		//脚本文件
		o1 := util.NewT(os.Args[1]).ToString()
		if strings.HasSuffix(o1, "json") {
			var mousePointList []Point
			jsonFile, err := os.Open(o1)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer jsonFile.Close()
			byteValue, _ := ioutil.ReadAll(jsonFile)
			err = json.Unmarshal(byteValue, &mousePointList)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(mousePointList)
			quit := false
			go WriteMouseLeftClickList(&quit, &mousePointList, dur)
			quit = robotgo.AddEvent("mright")
			return
		}
		//点击间隔
		if d := util.NewT(os.Args[1]).ToInt(); d > 12 {
			dur = d
		}
	}
	//点击次数
	times := 3
	if len(os.Args) > 2 {
		if d := util.NewT(os.Args[2]).ToInt(); d != 0 {
			times = d
		}
	}

	if startMsg := robotgo.ShowAlert("提示", "点击确定开始录制"+fmt.Sprint(times)+"个按键", "确定", "取消"); startMsg == 0 { //确定0，取消1
		var mousePointList []Point
		for i := 0; i < times; i++ {
			mousePointList = append(mousePointList, ReadMouseLeftClickPoint())
		}
		if endMsg := robotgo.ShowAlert("提示", "录制完毕, 点击确定开始播放，开始后单击右键退出", "确定", "取消"); endMsg == 0 { //确定0，取消1
			quit := false
			go WriteMouseLeftClickList(&quit, &mousePointList, dur)
			quit = robotgo.AddEvent("mright")
		}
	}
}
