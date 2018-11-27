package net

import (
	"bytes"
	"douyuDm/query"
	"douyuDm/util"
	"encoding/binary"
	"fmt"
	"golang.org/x/net/websocket"
	"strconv"
	"time"
)

func NewDYWebSocket(rid int) (*DYWebSocket, error) {
	obj := &DYWebSocket{
		buffer:     []byte{},
		readLength: 0,
		roomId:     rid,
	}
	cfg, err := websocket.NewConfig("wss://danmuproxy.douyu.com:8505/", "https://www.douyu.com/"+strconv.Itoa(rid))

	if err != nil {
		return nil, err
	}
	con, err := websocket.DialConfig(cfg)
	obj.con = con
	obj.wait = util.NewWaitGroup()
	return obj, nil
}

type DYWebSocket struct {
	buffer     []byte
	readLength uint32
	roomId     int
	OnMessage  func(msg string)
	con        *websocket.Conn
	wait       *util.WaitGroup
}

func (d *DYWebSocket) Run() {
	d.wait.Execute(d.read)
	d.login()
	d.joinGroup()
	d.wait.Execute(d.head)
	d.wait.Wait()
}

func (d *DYWebSocket) decode(data []byte) {
	if len(d.buffer) == 0 {
		d.buffer = data
	} else {
		d.buffer = append(d.buffer, data...)
	}

	for len(d.buffer) > 0 {
		if d.readLength == 0 {
			if len(d.buffer) < 4 {
				return
			}
			d.readLength = binary.LittleEndian.Uint32(d.buffer)
			d.buffer = d.buffer[4:]
		}

		if uint32(len(d.buffer)) < d.readLength {
			return
		}
		var msg = d.buffer[8 : d.readLength-1]
		d.buffer = d.buffer[d.readLength:]
		d.readLength = 0
		d.onMessage(string(msg))
	}
}

func (d *DYWebSocket) encode(msg string) []byte {
	buffer := append([]byte(msg), byte(0))
	length := 8 + len(buffer)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, uint32(length))
	binary.Write(buf, binary.LittleEndian, uint32(length))
	binary.Write(buf, binary.LittleEndian, int16(689))
	binary.Write(buf, binary.LittleEndian, int16(0))
	binary.Write(buf, binary.LittleEndian, buffer)
	return buf.Bytes()
}

func (d *DYWebSocket) onMessage(msg string) {
	kmp := query.Decode(msg)

	switch kmp["type"] {
	case "chatmsg":
		//fmt.Println("弹幕消息->", "["+kmp["nn"]+"]", kmp["txt"])
		break
	case "uenter":
		//大佬近直播间了
		break
	case "spbc":
		//map[bid:1002012_1541420233_5800 sn:百度0小宇 dn:冠磊敲可爱 gs:1 gfid:2050 rid:18927600 gn:盛典火箭 gb:0 eid:542 sid:189413723 drid:5810364 es:1 cl2:0 type:spbc gc:1 bgl:3 ifs:0 gid:0]
		//sn 送给 dn gn
		break
	case "mrkl":
		//心跳包
		break
	case "pingreq":
		//返回的服务器时间
		break
	case "loginres":
		//登录返回的数据
		break
	case "synexp":
		//map[type:synexp o_exp:1955840422 o_lev:60 o_minexp:1804500000 o_upexp:48659578 rid:475252]
		//貌似是有人送主播礼物返回的主播到下一级需要多少经验
		break

	default:
		fmt.Println(kmp)
	}
}

func (d *DYWebSocket) login() {
	var msg = map[string]interface{}{
		"type":   "loginreq",
		"roomid": d.roomId,
	}
	d.send(msg)
}

func (d *DYWebSocket) signVk(timestamp int64, devId string) string {
	str := util.Md5(strconv.FormatInt(timestamp, 10) + "r5*^5;}2#${XF[h+;'./.Q'1;,-]f'p[" + devId)
	return str
}

func (d *DYWebSocket) send(msg map[string]interface{}) (int, error) {
	return d.con.Write(d.encode(query.Encode(msg)))
}

func (d *DYWebSocket) read() {
	for {
		buf := make([]byte, 1024)
		n, err := d.con.Read(buf)
		if err != nil {
			panic(err)
		}
		d.decode(buf[0:n])
	}
}

func (d *DYWebSocket) joinGroup() {
	var join = map[string]interface{}{
		"type": "joingroup",
		"rid":  d.roomId,
		"gid":  -9999,
	}
	d.send(join)
}

func (d *DYWebSocket) head() {
	ticker := time.NewTicker(time.Millisecond * 4500)
	go func() {
		for {
			select {
			case <-ticker.C:
				var msg = map[string]interface{}{
					"type": "mrkl",
				}
				d.send(msg)
			}
		}
	}()
}
