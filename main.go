package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	// "strings"
	"time"
)

const (

	//AI權值
	SD4 = 94011
	ED4 = 5000
	SA3 = 2511
	EA3 = 2500
	ED3 = 100
	EA2 = 50
	ED2 = 20
	EA1 = 10
)

func main() {
	fmt.Println("start")

	http.HandleFunc("/", index)
	http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request:", r.URL.Path)
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.HandleFunc("/random_ai", randomAi)
	http.HandleFunc("/pichu_ai", pichuAi)

	go func() {
		err := http.ListenAndServe(":1920", nil)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	<-make(chan struct{})

}

func index(w http.ResponseWriter, r *http.Request) {

	var html = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8"/>
	</head>
	<body>
		<h3>五子棋的概念</h3>

		<canvas id="chess_panel" width="450" height="450"></canvas>
		<div id="log" style="display: inline-block;width:200px;height:450px;overflow: scroll;">Log:</div>
		<br/>
		<div id="choose_player_number">
			<button id="single_player">單人遊戲</button>
			<button id="multiple_player" disabled>雙人遊戲</button>
		</div>
		<div id="input_player_name" style="display:none;">
			請輸入玩家姓名：<input id="player_name"/>
			<button id="player_name_submit">確定</button>
		</div>
		<div id="choose_who_first" style="display:none;">
			<button id="ai_first">電腦先下</button>
			<button id="player_first">玩家先下</button>

		</div>
		<div id="game_status" style="display:none;">
			玩家 <span id="player_name_display">{{player_name}}</span> VS Pichu's AI <br/>
		</div>
		<div id="msg_bar">
			<span id="status_display">{{status_display}}</span>
		</div>
		<script src="//code.jquery.com/jquery-2.2.1.min.js"></script>
		<script src="js/main.js"></script>
	</body>
</html>
`

	w.Write([]byte(html))

}

func randomAi(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request:", r.URL.Path)

	var logs [][]int
	data, _ := ioutil.ReadAll(r.Body)
	// defer r.Body.Close()
	err := json.Unmarshal(data, &logs)
	if err != nil {
		fmt.Println("json error: " + err.Error())
	}

	panel := makePanel()

	loadPanel(panel, logs)

	var output []int

	for {
		var x = rand.Int() % 19
		var y = rand.Int() % 19
		fmt.Println("try (" + strconv.Itoa(x) + ", " + strconv.Itoa(y) + ")")
		if panel[x][y] == 0 {
			output = []int{x, y}
			break
		}
	}

	outputData, _ := json.Marshal(output)
	logs = append(logs, output)
	record, _ := json.Marshal(logs)
	logName := "logs/random_ai_" + r.RemoteAddr + "_" + time.Now().String() + ".log"
	// logName = strings.Replace(logName, " ", "", 0)
	// logName = strings.Replace(logName, "[", "", 0)
	// logName = strings.Replace(logName, "]", "", 0)
	// logName = strings.Replace(logName, "+", "", 0)
	fmt.Println("Log: " + logName)
	ioutil.WriteFile(logName, record, 0777)

	w.Write(outputData)

	// http.ServeFile(w, r, r.URL.Path[1:])
}
func pichuAi(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request:", r.URL.Path)

	var logs [][]int
	data, _ := ioutil.ReadAll(r.Body)
	// defer r.Body.Close()
	err := json.Unmarshal(data, &logs)
	if err != nil {
		fmt.Println("json error: " + err.Error())
	}

	panel := makePanel()

	loadPanel(panel, logs)

	var output []int = doPichuAi(panel)

	outputData, _ := json.Marshal(output)
	logs = append(logs, output)
	record, _ := json.Marshal(logs)
	logName := "logs/random_ai_" + r.RemoteAddr + "_" + time.Now().String() + ".log"
	// logName = strings.Replace(logName, " ", "", 0)
	// logName = strings.Replace(logName, "[", "", 0)
	// logName = strings.Replace(logName, "]", "", 0)
	// logName = strings.Replace(logName, "+", "", 0)
	fmt.Println("Log: " + logName)
	ioutil.WriteFile(logName, record, 0777)
	w.Write(outputData)
}

func doPichuAi(panel [][]int) []int {
	var x_max = 19
	var y_max = 19
	best := []int{1, 10, 10}
	stat := make([]int, 5, 5)
	a := makePanel()
	b := panel

	for ay := 0; ay < y_max; ay++ {
		for ax := 0; ax < x_max; ax++ {
			for k := 0; k < 5; k++ {
				if ax+k >= x_max {
					stat[k] = 3
					continue
				}
				stat[k] = b[ax+k][ay]
			}
			switch btov(stat) {
			case 80:
				a[ax][ay] += ED4
				break
			case 188:
				a[ax+1][ay] += ED4
				break
			case 224:
				a[ax+2][ay] += ED4
				break
			case 236:
				a[ax+3][ay] += ED4
				break
			case 240:
				a[ax+4][ay] += ED4
				break
			case 40:
				a[ax][ay] += SD4
				break
			case 94:
				a[ax+1][ay] += SD4
				break
			case 112:
				a[ax+2][ay] += SD4
				break
			case 118:
				a[ax+3][ay] += SD4
				break
			case 120:
				a[ax+4][ay] += SD4
				break
			case 39:
				a[ax+0][ay] += SA3
				a[ax+4][ay] += SA3
				break
			case 78:
				a[ax+0][ay] += EA3
				a[ax+4][ay] += EA3
				break
			case 26:
			case 13:
				a[ax+0][ay] += ED3
				a[ax+1][ay] += ED3
				break
			case 62:
			case 31:
				a[ax+0][ay] += ED3
				a[ax+2][ay] += ED3
				break
			case 74:
			case 37:
				a[ax+0][ay] += ED3
				a[ax+3][ay] += ED3
				break
			case 170:
			case 85:
				a[ax+1][ay] += ED3
				a[ax+2][ay] += ED3
				break
			case 182:
			case 91:
				a[ax+1][ay] += ED3
				a[ax+3][ay] += ED3
				break
			case 186:
			case 93:
				a[ax+1][ay] += ED3
				a[ax+4][ay] += ED3
				break
			case 218:
			case 109:
				a[ax+1][ay] += ED3
				a[ax+4][ay] += ED3
				break
			case 222:
			case 111:
				a[ax+2][ay] += ED3
				a[ax+4][ay] += ED3
				break
			case 234:
			case 117:
				a[ax+3][ay] += ED3
				a[ax+4][ay] += ED3
				break
			case 24:
			case 12:
				a[ax][ay] += EA2
				a[ax+1][ay] += EA2
				a[ax+4][ay] += EA2
				break
			case 60:
			case 30:
				a[ax][ay] += EA2
				a[ax+2][ay] += EA2
				a[ax+4][ay] += EA2
				break
			case 72:
			case 36:
				a[ax][ay] += EA2
				a[ax+3][ay] += EA2
				a[ax+4][ay] += EA2
				break
			case 8:
			case 4:
				a[ax][ay] += ED2
				a[ax+1][ay] += ED2
				a[ax+2][ay] += ED2
				break
			case 20:
			case 10:
				a[ax][ay] += ED2
				a[ax+1][ay] += ED2
				a[ax+3][ay] += ED2
				break
			case 56:
			case 28:
				a[ax][ay] += ED2
				a[ax+2][ay] += ED2
				a[ax+3][ay] += ED2
				break
			case 164:
			case 82:
				a[ax+1][ay] += ED2
				a[ax+2][ay] += ED2
				a[ax+3][ay] += ED2
				break
			case 168:
			case 84:
				a[ax+1][ay] += ED2
				a[ax+2][ay] += ED2
				a[ax+4][ay] += ED2
				break
			case 180:
			case 90:
				a[ax+1][ay] += ED2
				a[ax+3][ay] += ED2
				a[ax+4][ay] += ED2
				break
			case 216:
			case 108:
				a[ax+2][ay] += ED2
				a[ax+3][ay] += ED2
				a[ax+4][ay] += ED2
				break
			case 54:
			case 27:
				a[ax+0][ay] += EA1
				a[ax+2][ay] += EA1
				a[ax+3][ay] += EA1
				a[ax+4][ay] += EA1
				break
			case 18:
			case 9:
				a[ax+0][ay] += EA1
				a[ax+1][ay] += EA1
				a[ax+3][ay] += EA1
				a[ax+4][ay] += EA1
				break
			case 6:
			case 3:
				a[ax+0][ay] += EA1
				a[ax+1][ay] += EA1
				a[ax+2][ay] += EA1
				a[ax+4][ay] += EA1
				break
			case 162:
			case 81:
				a[ax+1][ay] += 1
				a[ax+2][ay] += 1
				a[ax+3][ay] += 1
				a[ax+4][ay] += 1
				break
			case 2:
			case 1:
				a[ax+0][ay] += 1
				a[ax+1][ay] += 1
				a[ax+2][ay] += 1
				a[ax+3][ay] += 1
				break
			}

			for k := 0; k < 5; k++ {
				if ay+k >= y_max {
					stat[k] = 3
					continue
				}
				stat[k] = b[ax][ay+k]
			}
			switch btov(stat) {
			case 80:
				a[ax][ay] += ED4
				break
			case 188:
				a[ax][ay+1] += ED4
				break
			case 224:
				a[ax][ay+2] += ED4
				break
			case 236:
				a[ax][ay+3] += ED4
				break
			case 240:
				a[ax][ay+4] += ED4
				break
			case 40:
				a[ax][ay] += SD4
				break
			case 94:
				a[ax][ay+1] += SD4
				break
			case 112:
				a[ax][ay+2] += SD4
				break
			case 118:
				a[ax][ay+3] += SD4
				break
			case 120:
				a[ax][ay+4] += SD4
				break
			case 39:
				a[ax][ay] += SA3
				a[ax][ay+4] += SA3
				break
			case 78:
				a[ax][ay+0] += EA3
				a[ax][ay+4] += EA3
				break
			case 26:
			case 13:
				a[ax][ay+0] += ED3
				a[ax][ay+1] += ED3
				break
			case 62:
			case 31:
				a[ax][ay+0] += ED3
				a[ax][ay+2] += ED3
				break
			case 74:
			case 37:
				a[ax][ay+0] += ED3
				a[ax][ay+3] += ED3
				break
			case 170:
			case 85:
				a[ax][ay+1] += ED3
				a[ax][ay+2] += ED3
				break
			case 182:
			case 91:
				a[ax][ay+1] += ED3
				a[ax][ay+3] += ED3
				break
			case 186:
			case 93:
				a[ax][ay+1] += ED3
				a[ax][ay+4] += ED3
				break
			case 218:
			case 109:
				a[ax][ay+1] += ED3
				a[ax][ay+4] += ED3
				break
			case 222:
			case 111:
				a[ax][ay+2] += ED3
				a[ax][ay+4] += ED3
				break
			case 234:
			case 117:
				a[ax][ay+3] += ED3
				a[ax][ay+4] += ED3
				break
			case 24:
			case 12:
				a[ax][ay] += EA2
				a[ax][ay+1] += EA2
				a[ax][ay+4] += EA2
				break
			case 60:
			case 30:
				a[ax][ay] += EA2
				a[ax][ay+2] += EA2
				a[ax][ay+4] += EA2
				break
			case 72:
			case 36:
				a[ax][ay] += EA2
				a[ax][ay+3] += EA2
				a[ax][ay+4] += EA2
				break
			case 8:
			case 4:
				a[ax][ay] += ED2
				a[ax][ay+1] += ED2
				a[ax][ay+2] += ED2
				break
			case 20:
			case 10:
				a[ax][ay] += ED2
				a[ax][ay+1] += ED2
				a[ax][ay+3] += ED2
				break
			case 56:
			case 28:
				a[ax][ay] += ED2
				a[ax][ay+2] += ED2
				a[ax][ay+3] += ED2
				break
			case 164:
			case 82:
				a[ax][ay+1] += ED2
				a[ax][ay+2] += ED2
				a[ax][ay+3] += ED2
				break
			case 168:
			case 84:
				a[ax][ay+1] += ED2
				a[ax][ay+2] += ED2
				a[ax][ay+4] += ED2
				break
			case 180:
			case 90:
				a[ax][ay+1] += ED2
				a[ax][ay+3] += ED2
				a[ax][ay+4] += ED2
				break
			case 216:
			case 108:
				a[ax][ay+2] += ED2
				a[ax][ay+3] += ED2
				a[ax][ay+4] += ED2
				break
			case 54:
			case 27:
				a[ax+0][ay] += EA1
				a[ax][ay+2] += EA1
				a[ax][ay+3] += EA1
				a[ax][ay+4] += EA1
				break
			case 18:
			case 9:
				a[ax+0][ay] += EA1
				a[ax][ay+1] += EA1
				a[ax][ay+3] += EA1
				a[ax][ay+4] += EA1
				break
			case 6:
			case 3:
				a[ax+0][ay] += EA1
				a[ax][ay+1] += EA1
				a[ax][ay+2] += EA1
				a[ax][ay+4] += EA1
				break
			case 162:
			case 81:
				a[ax][ay+1] += 1
				a[ax][ay+2] += 1
				a[ax][ay+3] += 1
				a[ax][ay+4] += 1
				break
			case 2:
			case 1:
				a[ax][ay] += 1
				a[ax][ay+1] += 1
				a[ax][ay+2] += 1
				a[ax][ay+3] += 1
				break
			}

			for k := 0; k < 5; k++ {
				if (ax+k >= x_max) || (ay+k >= y_max) {
					stat[k] = 3
					continue
				}
				stat[k] = b[ax+k][ay+k]
			}
			switch btov(stat) {
			case 80:
				a[ax][ay] += ED4
				break
			case 188:
				a[ax+1][ay+1] += ED4
				break
			case 224:
				a[ax+2][ay+2] += ED4
				break
			case 236:
				a[ax+3][ay+3] += ED4
				break
			case 240:
				a[ax+4][ay+4] += ED4
				break
			case 40:
				a[ax][ay] += SD4
				break
			case 94:
				a[ax+1][ay+1] += SD4
				break
			case 112:
				a[ax+2][ay+2] += SD4
				break
			case 118:
				a[ax+3][ay+3] += SD4
				break
			case 120:
				a[ax+4][ay+4] += SD4
				break
			case 39:
				a[ax+0][ay] += SA3
				a[ax+4][ay+4] += SA3
				break
			case 78:
				a[ax][ay+0] += EA3
				a[ax+4][ay+4] += EA3
				break
			case 26:
			case 13:
				a[ax][ay+0] += ED3
				a[ax+1][ay+1] += ED3
				break
			case 62:
			case 31:
				a[ax][ay+0] += ED3
				a[ax+2][ay+2] += ED3
				break
			case 74:
			case 37:
				a[ax][ay+0] += ED3
				a[ax+3][ay+3] += ED3
				break
			case 170:
			case 85:
				a[ax+1][ay+1] += ED3
				a[ax+2][ay+2] += ED3
				break
			case 182:
			case 91:
				a[ax+1][ay+1] += ED3
				a[ax+3][ay+3] += ED3
				break
			case 186:
			case 93:
				a[ax+1][ay+1] += ED3
				a[ax+4][ay+4] += ED3
				break
			case 218:
			case 109:
				a[ax+1][ay+1] += ED3
				a[ax+4][ay+4] += ED3
				break
			case 222:
			case 111:
				a[ax+2][ay+2] += ED3
				a[ax+4][ay+4] += ED3
				break
			case 234:
			case 117:
				a[ax+3][ay+3] += ED3
				a[ax+4][ay+4] += ED3
				break
			case 24:
			case 12:
				a[ax][ay] += EA2
				a[ax+1][ay+1] += EA2
				a[ax+4][ay+4] += EA2
				break
			case 60:
			case 30:
				a[ax][ay] += EA2
				a[ax+2][ay+2] += EA2
				a[ax+4][ay+4] += EA2
				break
			case 72:
			case 36:
				a[ax][ay] += EA2
				a[ax+3][ay+3] += EA2
				a[ax+4][ay+4] += EA2
				break
			case 8:
			case 4:
				a[ax][ay] += ED2
				a[ax+1][ay+1] += ED2
				a[ax+2][ay+2] += ED2
				break
			case 20:
			case 10:
				a[ax][ay] += ED2
				a[ax+1][ay+1] += ED2
				a[ax+3][ay+3] += ED2
				break
			case 56:
			case 28:
				a[ax][ay] += ED2
				a[ax+2][ay+2] += ED2
				a[ax+3][ay+3] += ED2
				break
			case 164:
			case 82:
				a[ax+1][ay+1] += ED2
				a[ax+2][ay+2] += ED2
				a[ax+3][ay+3] += ED2
				break
			case 168:
			case 84:
				a[ax+1][ay+1] += ED2
				a[ax+2][ay+2] += ED2
				a[ax+4][ay+4] += ED2
				break
			case 180:
			case 90:
				a[ax+1][ay+1] += ED2
				a[ax+3][ay+3] += ED2
				a[ax+4][ay+4] += ED2
				break
			case 216:
			case 108:
				a[ax+2][ay+2] += ED2
				a[ax+3][ay+3] += ED2
				a[ax+4][ay+4] += ED2
				break
			case 54:
			case 27:
				a[ax+0][ay] += EA1
				a[ax+2][ay+2] += EA1
				a[ax+3][ay+3] += EA1
				a[ax+4][ay+4] += EA1
				break
			case 18:
			case 9:
				a[ax+0][ay] += EA1
				a[ax+1][ay+1] += EA1
				a[ax+3][ay+3] += EA1
				a[ax+4][ay+4] += EA1
				break
			case 6:
			case 3:
				a[ax+0][ay] += EA1
				a[ax+1][ay+1] += EA1
				a[ax+2][ay+2] += EA1
				a[ax+4][ay+4] += EA1
				break
			case 162:
			case 81:
				a[ax+1][ay+1] += 1
				a[ax+2][ay+2] += 1
				a[ax+3][ay+3] += 1
				a[ax+4][ay+4] += 1
				break
			case 2:
			case 1:
				a[ax][ay] += 1
				a[ax+1][ay+1] += 1
				a[ax+2][ay+2] += 1
				a[ax+3][ay+3] += 1
				break
			}

			for k := 0; k < 5; k++ {
				if (ax+k >= x_max) || (ay-k < 0) {
					stat[k] = 3
					continue
				}
				stat[k] = b[ax+k][ay-k]
			}

			switch btov(stat) {
			case 80:
				a[ax][ay] += ED4
				break
			case 188:
				a[ax+1][ay-1] += ED4
				break
			case 224:
				a[ax+2][ay-2] += ED4
				break
			case 236:
				a[ax+3][ay-3] += ED4
				break
			case 240:
				a[ax+4][ay-4] += ED4
				break
			case 40:
				a[ax][ay] += SD4
				break
			case 94:
				a[ax+1][ay-1] += SD4
				break
			case 112:
				a[ax+2][ay-2] += SD4
				break
			case 118:
				a[ax+3][ay-3] += SD4
				break
			case 120:
				a[ax+4][ay-4] += SD4
				break
			case 39:
				a[ax+0][ay] += SA3
				a[ax+4][ay-4] += SA3
				break
			case 78:
				a[ax][ay+0] += EA3
				a[ax+4][ay-4] += EA3
				break
			case 26:
			case 13:
				a[ax][ay+0] += ED3
				a[ax+1][ay-1] += ED3
				break
			case 62:
			case 31:
				a[ax][ay+0] += ED3
				a[ax+2][ay-2] += ED3
				break
			case 74:
			case 37:
				a[ax][ay+0] += ED3
				a[ax+3][ay-3] += ED3
				break
			case 170:
			case 85:
				a[ax+1][ay-1] += ED3
				a[ax+2][ay-2] += ED3
				break
			case 182:
			case 91:
				a[ax+1][ay-1] += ED3
				a[ax+3][ay-3] += ED3
				break
			case 186:
			case 93:
				a[ax+1][ay-1] += ED3
				a[ax+4][ay-4] += ED3
				break
			case 218:
			case 109:
				a[ax+1][ay-1] += ED3
				a[ax+4][ay-4] += ED3
				break
			case 222:
			case 111:
				a[ax+2][ay-2] += ED3
				a[ax+4][ay-4] += ED3
				break
			case 234:
			case 117:
				a[ax+3][ay-3] += ED3
				a[ax+4][ay-4] += ED3
				break
			case 24:
			case 12:
				a[ax][ay] += EA2
				a[ax+1][ay-1] += EA2
				a[ax+4][ay-4] += EA2
				break
			case 60:
			case 30:
				a[ax][ay] += EA2
				a[ax+2][ay-2] += EA2
				a[ax+4][ay-4] += EA2
				break
			case 72:
			case 36:
				a[ax][ay] += EA2
				a[ax+3][ay-3] += EA2
				a[ax+4][ay-4] += EA2
				break
			case 8:
			case 4:
				a[ax][ay] += ED2
				a[ax+1][ay-1] += ED2
				a[ax+2][ay-2] += ED2
				break
			case 20:
			case 10:
				a[ax][ay] += ED2
				a[ax+1][ay-1] += ED2
				a[ax+3][ay-3] += ED2
				break
			case 56:
			case 28:
				a[ax][ay] += ED2
				a[ax+2][ay-2] += ED2
				a[ax+3][ay-3] += ED2
				break
			case 164:
			case 82:
				a[ax+1][ay-1] += ED2
				a[ax+2][ay-2] += ED2
				a[ax+3][ay-3] += ED2
				break
			case 168:
			case 84:
				a[ax+1][ay-1] += ED2
				a[ax+2][ay-2] += ED2
				a[ax+4][ay-4] += ED2
				break
			case 180:
			case 90:
				a[ax+1][ay-1] += ED2
				a[ax+3][ay-3] += ED2
				a[ax+4][ay-4] += ED2
				break
			case 216:
			case 108:
				a[ax+2][ay-2] += ED2
				a[ax+3][ay-3] += ED2
				a[ax+4][ay-4] += ED2
				break
			case 54:
			case 27:
				a[ax+0][ay] += EA1
				a[ax+2][ay-2] += EA1
				a[ax+3][ay-3] += EA1
				a[ax+4][ay-4] += EA1
				break
			case 18:
			case 9:
				a[ax+0][ay] += EA1
				a[ax+1][ay-1] += EA1
				a[ax+3][ay-3] += EA1
				a[ax+4][ay-4] += EA1
				break
			case 6:
			case 3:
				a[ax+0][ay] += EA1
				a[ax+1][ay-1] += EA1
				a[ax+2][ay-2] += EA1
				a[ax+4][ay-4] += EA1
				break
			case 162:
			case 81:
				a[ax+1][ay-1] += 1
				a[ax+2][ay-2] += 1
				a[ax+3][ay-3] += 1
				a[ax+4][ay-4] += 1
				break
			case 2:
			case 1:
				a[ax][ay] += 1
				a[ax+1][ay-1] += 1
				a[ax+2][ay-2] += 1
				a[ax+3][ay-3] += 1
				break
			}

			for k := 0; k < 5; k++ {
				if ax-k < 0 {
					stat[k] = 3
					continue
				}
				stat[k] = b[ax-k][ay]
			}

			switch btov(stat) {
			case 80:
				a[ax][ay] += ED4
				break
			case 188:
				a[ax-1][ay] += ED4
				break
			case 224:
				a[ax-2][ay] += ED4
				break
			case 236:
				a[ax-3][ay] += ED4
				break
			case 240:
				a[ax-4][ay] += ED4
				break
			case 40:
				a[ax][ay] += SD4
				break
			case 94:
				a[ax-1][ay] += SD4
				break
			case 112:
				a[ax-2][ay] += SD4
				break
			case 118:
				a[ax-3][ay] += SD4
				break
			case 120:
				a[ax-4][ay] += SD4
				break
			case 39:
				a[ax+0][ay] += SA3
				a[ax-4][ay] += SA3
				break
			case 78:
				a[ax][ay+0] += EA3
				a[ax-4][ay] += EA3
				break
			case 26:
			case 13:
				a[ax][ay+0] += ED3
				a[ax-1][ay] += ED3
				break
			case 62:
			case 31:
				a[ax][ay+0] += ED3
				a[ax-2][ay] += ED3
				break
			case 74:
			case 37:
				a[ax][ay+0] += ED3
				a[ax-3][ay] += ED3
				break
			case 170:
			case 85:
				a[ax-1][ay] += ED3
				a[ax-2][ay] += ED3
				break
			case 182:
			case 91:
				a[ax-1][ay] += ED3
				a[ax-3][ay] += ED3
				break
			case 186:
			case 93:
				a[ax-1][ay] += ED3
				a[ax-4][ay] += ED3
				break
			case 218:
			case 109:
				a[ax-1][ay] += ED3
				a[ax-4][ay] += ED3
				break
			case 222:
			case 111:
				a[ax-2][ay] += ED3
				a[ax-4][ay] += ED3
				break
			case 234:
			case 117:
				a[ax-3][ay] += ED3
				a[ax-4][ay] += ED3
				break
			case 24:
			case 12:
				a[ax][ay] += EA2
				a[ax-1][ay] += EA2
				a[ax-4][ay] += EA2
				break
			case 60:
			case 30:
				a[ax][ay] += EA2
				a[ax-2][ay] += EA2
				a[ax-4][ay] += EA2
				break
			case 72:
			case 36:
				a[ax][ay] += EA2
				a[ax-3][ay] += EA2
				a[ax-4][ay] += EA2
				break
			case 8:
			case 4:
				a[ax][ay] += ED2
				a[ax-1][ay] += ED2
				a[ax-2][ay] += ED2
				break
			case 20:
			case 10:
				a[ax][ay] += ED2
				a[ax-1][ay] += ED2
				a[ax-3][ay] += ED2
				break
			case 56:
			case 28:
				a[ax][ay] += ED2
				a[ax-2][ay] += ED2
				a[ax-3][ay] += ED2
				break
			case 164:
			case 82:
				a[ax-1][ay] += ED2
				a[ax-2][ay] += ED2
				a[ax-3][ay] += ED2
				break
			case 168:
			case 84:
				a[ax-1][ay] += ED2
				a[ax-2][ay] += ED2
				a[ax-4][ay] += ED2
				break
			case 180:
			case 90:
				a[ax-1][ay] += ED2
				a[ax-3][ay] += ED2
				a[ax-4][ay] += ED2
				break
			case 216:
			case 108:
				a[ax-2][ay] += ED2
				a[ax-3][ay] += ED2
				a[ax-4][ay] += ED2
				break
			case 54:
			case 27:
				a[ax+0][ay] += EA1
				a[ax-2][ay] += EA1
				a[ax-3][ay] += EA1
				a[ax-4][ay] += EA1
				break
			case 18:
			case 9:
				a[ax+0][ay] += EA1
				a[ax-1][ay] += EA1
				a[ax-3][ay] += EA1
				a[ax-4][ay] += EA1
				break
			case 6:
			case 3:
				a[ax+0][ay] += EA1
				a[ax-1][ay] += EA1
				a[ax-2][ay] += EA1
				a[ax-4][ay] += EA1
				break
			case 162:
			case 81:
				a[ax-1][ay] += 1
				a[ax-2][ay] += 1
				a[ax-3][ay] += 1
				a[ax-4][ay] += 1
				break
			case 2:
			case 1:
				a[ax][ay] += 1
				a[ax-1][ay] += 1
				a[ax-2][ay] += 1
				a[ax-3][ay] += 1
				break
			}

			for k := 0; k < 5; k++ {
				if ay-k < 0 {
					stat[k] = 3
					continue
				}
				stat[k] = b[ax][ay-k]
			}

			switch btov(stat) {
			case 80:
				a[ax][ay] += ED4
				break
			case 188:
				a[ax][ay-1] += ED4
				break
			case 224:
				a[ax][ay-2] += ED4
				break
			case 236:
				a[ax][ay-3] += ED4
				break
			case 240:
				a[ax][ay-4] += ED4
				break
			case 40:
				a[ax][ay] += SD4
				break
			case 94:
				a[ax][ay-1] += SD4
				break
			case 112:
				a[ax][ay-2] += SD4
				break
			case 118:
				a[ax][ay-3] += SD4
				break
			case 120:
				a[ax][ay-4] += SD4
				break
			case 39:
				a[ax+0][ay] += SA3
				a[ax][ay-4] += SA3
				break
			case 78:
				a[ax][ay+0] += EA3
				a[ax][ay-4] += EA3
				break
			case 26:
			case 13:
				a[ax][ay+0] += ED3
				a[ax][ay-1] += ED3
				break
			case 62:
			case 31:
				a[ax][ay+0] += ED3
				a[ax][ay-2] += ED3
				break
			case 74:
			case 37:
				a[ax][ay+0] += ED3
				a[ax][ay-3] += ED3
				break
			case 170:
			case 85:
				a[ax][ay-1] += ED3
				a[ax][ay-2] += ED3
				break
			case 182:
			case 91:
				a[ax][ay-1] += ED3
				a[ax][ay-3] += ED3
				break
			case 186:
			case 93:
				a[ax][ay-1] += ED3
				a[ax][ay-4] += ED3
				break
			case 218:
			case 109:
				a[ax][ay-1] += ED3
				a[ax][ay-4] += ED3
				break
			case 222:
			case 111:
				a[ax][ay-2] += ED3
				a[ax][ay-4] += ED3
				break
			case 234:
			case 117:
				a[ax][ay-3] += ED3
				a[ax][ay-4] += ED3
				break
			case 24:
			case 12:
				a[ax][ay] += EA2
				a[ax][ay-1] += EA2
				a[ax][ay-4] += EA2
				break
			case 60:
			case 30:
				a[ax][ay] += EA2
				a[ax][ay-2] += EA2
				a[ax][ay-4] += EA2
				break
			case 72:
			case 36:
				a[ax][ay] += EA2
				a[ax][ay-3] += EA2
				a[ax][ay-4] += EA2
				break
			case 8:
			case 4:
				a[ax][ay] += ED2
				a[ax][ay-1] += ED2
				a[ax][ay-2] += ED2
				break
			case 20:
			case 10:
				a[ax][ay] += ED2
				a[ax][ay-1] += ED2
				a[ax][ay-3] += ED2
				break
			case 56:
			case 28:
				a[ax][ay] += ED2
				a[ax][ay-2] += ED2
				a[ax][ay-3] += ED2
				break
			case 164:
			case 82:
				a[ax][ay-1] += ED2
				a[ax][ay-2] += ED2
				a[ax][ay-3] += ED2
				break
			case 168:
			case 84:
				a[ax][ay-1] += ED2
				a[ax][ay-2] += ED2
				a[ax][ay-4] += ED2
				break
			case 180:
			case 90:
				a[ax][ay-1] += ED2
				a[ax][ay-3] += ED2
				a[ax][ay-4] += ED2
				break
			case 216:
			case 108:
				a[ax][ay-2] += ED2
				a[ax][ay-3] += ED2
				a[ax][ay-4] += ED2
				break
			case 54:
			case 27:
				a[ax+0][ay] += EA1
				a[ax][ay-2] += EA1
				a[ax][ay-3] += EA1
				a[ax][ay-4] += EA1
				break
			case 18:
			case 9:
				a[ax+0][ay] += EA1
				a[ax][ay-1] += EA1
				a[ax][ay-3] += EA1
				a[ax][ay-4] += EA1
				break
			case 6:
			case 3:
				a[ax+0][ay] += EA1
				a[ax][ay-1] += EA1
				a[ax][ay-2] += EA1
				a[ax][ay-4] += EA1
				break
			case 162:
			case 81:
				a[ax][ay-1] += 1
				a[ax][ay-2] += 1
				a[ax][ay-3] += 1
				a[ax][ay-4] += 1
				break
			case 2:
			case 1:
				a[ax][ay] += 1
				a[ax][ay-1] += 1
				a[ax][ay-2] += 1
				a[ax][ay-3] += 1
				break
			}

			for k := 0; k < 5; k++ {
				if (ax-k < 0) || (ay-k < 0) {
					stat[k] = 3
					continue
				}
				stat[k] = b[ax-k][ay-k]
			}

			switch btov(stat) {
			case 80:
				a[ax][ay] += ED4
				break
			case 188:
				a[ax-1][ay-1] += ED4
				break
			case 224:
				a[ax-2][ay-2] += ED4
				break
			case 236:
				a[ax-3][ay-3] += ED4
				break
			case 240:
				a[ax-4][ay-4] += ED4
				break
			case 40:
				a[ax][ay] += SD4
				break
			case 94:
				a[ax-1][ay-1] += SD4
				break
			case 112:
				a[ax-2][ay-2] += SD4
				break
			case 118:
				a[ax-3][ay-3] += SD4
				break
			case 120:
				a[ax-4][ay-4] += SD4
				break
			case 39:
				a[ax+0][ay] += SA3
				a[ax-4][ay-4] += SA3
				break
			case 78:
				a[ax][ay+0] += EA3
				a[ax-4][ay-4] += EA3
				break
			case 26:
			case 13:
				a[ax][ay+0] += ED3
				a[ax-1][ay-1] += ED3
				break
			case 62:
			case 31:
				a[ax][ay+0] += ED3
				a[ax-2][ay-2] += ED3
				break
			case 74:
			case 37:
				a[ax][ay+0] += ED3
				a[ax-3][ay-3] += ED3
				break
			case 170:
			case 85:
				a[ax-1][ay-1] += ED3
				a[ax-2][ay-2] += ED3
				break
			case 182:
			case 91:
				a[ax-1][ay-1] += ED3
				a[ax-3][ay-3] += ED3
				break
			case 186:
			case 93:
				a[ax-1][ay-1] += ED3
				a[ax-4][ay-4] += ED3
				break
			case 218:
			case 109:
				a[ax-1][ay-1] += ED3
				a[ax-4][ay-4] += ED3
				break
			case 222:
			case 111:
				a[ax-2][ay-2] += ED3
				a[ax-4][ay-4] += ED3
				break
			case 234:
			case 117:
				a[ax-3][ay-3] += ED3
				a[ax-4][ay-4] += ED3
				break
			case 24:
			case 12:
				a[ax][ay] += EA2
				a[ax-1][ay-1] += EA2
				a[ax-4][ay-4] += EA2
				break
			case 60:
			case 30:
				a[ax][ay] += EA2
				a[ax-2][ay-2] += EA2
				a[ax-4][ay-4] += EA2
				break
			case 72:
			case 36:
				a[ax][ay] += EA2
				a[ax-3][ay-3] += EA2
				a[ax-4][ay-4] += EA2
				break
			case 8:
			case 4:
				a[ax][ay] += ED2
				a[ax-1][ay-1] += ED2
				a[ax-2][ay-2] += ED2
				break
			case 20:
			case 10:
				a[ax][ay] += ED2
				a[ax-1][ay-1] += ED2
				a[ax-3][ay-3] += ED2
				break
			case 56:
			case 28:
				a[ax][ay] += ED2
				a[ax-2][ay-2] += ED2
				a[ax-3][ay-3] += ED2
				break
			case 164:
			case 82:
				a[ax-1][ay-1] += ED2
				a[ax-2][ay-2] += ED2
				a[ax-3][ay-3] += ED2
				break
			case 168:
			case 84:
				a[ax-1][ay-1] += ED2
				a[ax-2][ay-2] += ED2
				a[ax-4][ay-4] += ED2
				break
			case 180:
			case 90:
				a[ax-1][ay-1] += ED2
				a[ax-3][ay-3] += ED2
				a[ax-4][ay-4] += ED2
				break
			case 216:
			case 108:
				a[ax-2][ay-2] += ED2
				a[ax-3][ay-3] += ED2
				a[ax-4][ay-4] += ED2
				break
			case 54:
			case 27:
				a[ax+0][ay] += EA1
				a[ax-2][ay-2] += EA1
				a[ax-3][ay-3] += EA1
				a[ax-4][ay-4] += EA1
				break
			case 18:
			case 9:
				a[ax+0][ay] += EA1
				a[ax-1][ay-1] += EA1
				a[ax-3][ay-3] += EA1
				a[ax-4][ay-4] += EA1
				break
			case 6:
			case 3:
				a[ax+0][ay] += EA1
				a[ax-1][ay-1] += EA1
				a[ax-2][ay-2] += EA1
				a[ax-4][ay-4] += EA1
				break
			case 162:
			case 81:
				a[ax-1][ay-1] += 1
				a[ax-2][ay-2] += 1
				a[ax-3][ay-3] += 1
				a[ax-4][ay-4] += 1
				break
			case 2:
			case 1:
				a[ax][ay] += 1
				a[ax-1][ay-1] += 1
				a[ax-2][ay-2] += 1
				a[ax-3][ay-3] += 1
				break
			}
			for k := 0; k < 5; k++ {
				if (ax-k < 0) || (ay+k >= y_max) {
					stat[k] = 3
					continue
				}
				stat[k] = b[ax-k][ay+k]
			}

			switch btov(stat) {
			case 80:
				a[ax][ay] += ED4
				break
			case 188:
				a[ax-1][ay+1] += ED4
				break
			case 224:
				a[ax-2][ay+2] += ED4
				break
			case 236:
				a[ax-3][ay+3] += ED4
				break
			case 240:
				a[ax-4][ay+4] += ED4
				break
			case 40:
				a[ax][ay] += SD4
				break
			case 94:
				a[ax-1][ay+1] += SD4
				break
			case 112:
				a[ax-2][ay+2] += SD4
				break
			case 118:
				a[ax-3][ay+3] += SD4
				break
			case 120:
				a[ax-4][ay+4] += SD4
				break
			case 39:
				a[ax+0][ay] += SA3
				a[ax-4][ay+4] += SA3
				break
			case 78:
				a[ax][ay+0] += EA3
				a[ax-4][ay+4] += EA3
				break
			case 26:
			case 13:
				a[ax][ay+0] += ED3
				a[ax-1][ay+1] += ED3
				break
			case 62:
			case 31:
				a[ax][ay+0] += ED3
				a[ax-2][ay+2] += ED3
				break
			case 74:
			case 37:
				a[ax][ay+0] += ED3
				a[ax-3][ay+3] += ED3
				break
			case 170:
			case 85:
				a[ax-1][ay+1] += ED3
				a[ax-2][ay+2] += ED3
				break
			case 182:
			case 91:
				a[ax-1][ay+1] += ED3
				a[ax-3][ay+3] += ED3
				break
			case 186:
			case 93:
				a[ax-1][ay+1] += ED3
				a[ax-4][ay+4] += ED3
				break
			case 218:
			case 109:
				a[ax-1][ay+1] += ED3
				a[ax-4][ay+4] += ED3
				break
			case 222:
			case 111:
				a[ax-2][ay+2] += ED3
				a[ax-4][ay+4] += ED3
				break
			case 234:
			case 117:
				a[ax-3][ay+3] += ED3
				a[ax-4][ay+4] += ED3
				break
			case 24:
			case 12:
				a[ax][ay] += EA2
				a[ax-1][ay+1] += EA2
				a[ax-4][ay+4] += EA2
				break
			case 60:
			case 30:
				a[ax][ay] += EA2
				a[ax-2][ay+2] += EA2
				a[ax-4][ay+4] += EA2
				break
			case 72:
			case 36:
				a[ax][ay] += EA2
				a[ax-3][ay+3] += EA2
				a[ax-4][ay+4] += EA2
				break
			case 8:
			case 4:
				a[ax][ay] += ED2
				a[ax-1][ay+1] += ED2
				a[ax-2][ay+2] += ED2
				break
			case 20:
			case 10:
				a[ax][ay] += ED2
				a[ax-1][ay+1] += ED2
				a[ax-3][ay+3] += ED2
				break
			case 56:
			case 28:
				a[ax][ay] += ED2
				a[ax-2][ay+2] += ED2
				a[ax-3][ay+3] += ED2
				break
			case 164:
			case 82:
				a[ax-1][ay+1] += ED2
				a[ax-2][ay+2] += ED2
				a[ax-3][ay+3] += ED2
				break
			case 168:
			case 84:
				a[ax-1][ay+1] += ED2
				a[ax-2][ay+2] += ED2
				a[ax-4][ay+4] += ED2
				break
			case 180:
			case 90:
				a[ax-1][ay+1] += ED2
				a[ax-3][ay+3] += ED2
				a[ax-4][ay+4] += ED2
				break
			case 216:
			case 108:
				a[ax-2][ay+2] += ED2
				a[ax-3][ay+3] += ED2
				a[ax-4][ay+4] += ED2
				break
			case 54:
			case 27:
				a[ax+0][ay] += EA1
				a[ax-2][ay+2] += EA1
				a[ax-3][ay+3] += EA1
				a[ax-4][ay+4] += EA1
				break
			case 18:
			case 9:
				a[ax+0][ay] += EA1
				a[ax-1][ay+1] += EA1
				a[ax-3][ay+3] += EA1
				a[ax-4][ay+4] += EA1
				break
			case 6:
			case 3:
				a[ax+0][ay] += EA1
				a[ax-1][ay+1] += EA1
				a[ax-2][ay+2] += EA1
				a[ax-4][ay+4] += EA1
				break
			case 162:
			case 81:
				a[ax-1][ay+1] += 1
				a[ax-2][ay+2] += 1
				a[ax-3][ay+3] += 1
				a[ax-4][ay+4] += 1
				break
			case 2:
			case 1:
				a[ax][ay] += 1
				a[ax-1][ay+1] += 1
				a[ax-2][ay+2] += 1
				a[ax-3][ay+3] += 1
				break
			}
		}
	}
	for ay := 0; ay < y_max; ay++ {
		for ax := 0; ax < x_max; ax++ {
			if (best[0] < a[ax][ay]) && (b[ax][ay] == 0) {
				best[0] = a[ax][ay]
				best[1] = ax
				best[2] = ay
			}
		}
	}

	return []int{best[1], best[2]}
	// return set(best[1], best[2], 1+(now%2))
}
func btov(stat []int) int {
	if (stat[4] == 3) || (stat[3] == 3) || (stat[2] == 3) || (stat[1] == 3) || (stat[0] == 3) {
		return 243
	}
	return (stat[4] + 3*stat[3] + 9*stat[2] + 27*stat[1] + 81*stat[0])
}

func loadPanel(panel [][]int, logs [][]int) {
	fmt.Println("len:", len(logs))
	var player = 1
	if len(logs)%2 == 0 {
		// AI is Player 1
		player = 2
	} else {
		// AI is Player 2
		player = 1
	}
	for _, pos := range logs {
		panel[pos[0]][pos[1]] = player
		player ^= 3
	}
}

func makePanel() [][]int {
	ans := make([][]int, 19, 19)
	for i := 0; i < 19; i++ {
		ans[i] = make([]int, 19, 19)
	}
	return ans
}

