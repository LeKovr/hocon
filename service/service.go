package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	gen "SELF/zgen/go/proto"
)

const (
	port   = ":38899"
	CmdSet = "setPilot"
	CmdGet = "getPilot"
	OffSceneID = 0
)

// TODO: set this in config
var scene = map[gen.LampScene]int{
	gen.LampScene_OFF:   OffSceneID,
	gen.LampScene_DAY:   12,
	gen.LampScene_NIGHT: 14,
}

type Params struct {
	State    bool   `json:"state"`
	SceneID  int    `json:"sceneId,omitempty"`
	SceneKey string `json:"sceneKey,omitempty"`
}

type Command struct {
	Method string `json:"method"`
	Params Params `json:"params"`
}

type Result struct {
	Params
	Success bool `json:"success,omitempty"`
}

type ResultEnvelope struct {
	Method string `json:"method"`
	Result Result `json:"result"`
}

type HoconServiceImpl struct {
	gen.UnimplementedHoconServiceServer
}

func (srv HoconServiceImpl) LampControl(ctx context.Context, cmd *gen.LampStatus) (rv *gen.LampStatus, err error) {
	fmt.Printf("CALL: %s / %d\n", cmd.Id, cmd.Scene)
	var scene *gen.LampScene
	scene, err = run(cmd.Id, cmd.Scene)
	if err != nil {
		return nil, err
	}
	cmd.Scene = *scene
	return cmd, nil
}

func run(lamp string, key gen.LampScene) (*gen.LampScene, error) {
	var cmd Command
	if key == gen.LampScene_UNKNOWN {
		cmd = Command{Method: CmdGet}
	} else {
		id, ok := scene[key]
		if !ok {
			fmt.Println("Unknown scene:", key)
			return nil, errors.New("Unknown scene:")
		}
		cmd = Command{Method: CmdSet}
		if id != OffSceneID {
			cmd.Params.State = true
			cmd.Params.SceneID = id
		}
	}
	out, err := json.Marshal(cmd)
	if err != nil {
		fmt.Println("json error:", err)
		return nil, errors.New("json error")
	}
	fmt.Println("Send:", string(out))

	rv, err := send(lamp+port, out)
	if err != nil {
		fmt.Println("Resolve:", err)
		return nil, errors.New("resolve")
	}
	var res ResultEnvelope
	err = json.Unmarshal(rv, &res)
	if err != nil {
		fmt.Println("Result:", err)
		return nil, errors.New("unmarshal")
	}
	if cmd.Method == CmdSet {
		// fill empty response fields
		res.Result.Params.SceneID = cmd.Params.SceneID
		res.Result.Params.State = cmd.Params.State
	}
	//		fmt.Printf("OK: %v\n", res.Result.Success)
	//	} else {
	if !res.Result.Params.State {
		res.Result.Params.SceneKey = "off"
	} else {
		for k, v := range scene {
			if res.Result.Params.SceneID == v {
				res.Result.Params.SceneKey = gen.LampScene(k).String()
				key = gen.LampScene(k)
				break
			}
		}

	}
	fmt.Printf("Status: %+v\n", res.Result.Params)
	//	}
	return &key, nil
}

func send(host string, cmd []byte) (rv []byte, err error) {
	var s *net.UDPAddr
	s, err = net.ResolveUDPAddr("udp4", host)
	if err != nil {
		return
	}
	var conn *net.UDPConn
	conn, err = net.DialUDP("udp4", nil, s)
	if err != nil {
		return
	}
	//fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
	defer conn.Close()

	_, err = conn.Write(append(cmd, []byte("\n")...))
	if err != nil {
		return
	}

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return
	}
	// return string(buffer[0:n]), nil
	return buffer[0:n], nil
}
