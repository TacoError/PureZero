package main

import (
	"PureZero/cheats"
	"PureZero/players"
	"PureZero/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

func doError(happened string) {
	fmt.Println(happened)
	time.Sleep(3 * time.Second)
	os.Exit(1)
}

type IP struct {
	Ip string `json:"ip"`
}

func main() {
	cmd := exec.Command("title", "PureZero")
	_ = cmd.Run()

	ip := new(IP)
	ipCached, err := ioutil.ReadFile("ip.cache")
	skip := false
	useIP := ""
	if err == nil {
		err := json.Unmarshal(ipCached, ip)
		if err != nil {
			return
		}
		useIP = ip.Ip
		fmt.Print(fmt.Sprintf("Found cached IP (%s) should we use it? (y/n) > ", useIP))
		var yes string
		_, err = fmt.Scan(&yes)
		if err != nil {
			return
		}
		if yes == "n" {
			skip = true
		}
	}
	if err != nil || (err == nil && skip) {
		fmt.Print("Server IP (in form of ip:port) > ")
		_, err := fmt.Scan(&useIP)
		if err != nil {
			doError("An error has occurred reading input. Please try again.")
			return
		}
		store, _ := json.Marshal(IP{Ip: useIP})
		err = ioutil.WriteFile("ip.cache", store, 0644)
		if err != nil {
			return
		}
	}
	proxy(useIP)
}

func proxy(ip string) {
	src := utils.TokenSource()
	p, err := minecraft.NewForeignStatusProvider(ip)
	if err != nil {
		doError("An error has occurred starting foreign status provider. Please try again.")
		return
	}

	listener, err := minecraft.ListenConfig{
		StatusProvider: p,
	}.Listen("raknet", "0.0.0.0:19132")
	if err != nil {
		panic(err)
	}
	defer func(listener *minecraft.Listener) {
		err := listener.Close()
		if err != nil {
			panic(err)
		}
	}(listener)
	fmt.Println("Please join 127.0.0.1:19132")
	for {
		c, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(c.(*minecraft.Conn), listener, ip, src)
	}
}

func handleConn(conn *minecraft.Conn, listener *minecraft.Listener, ip string, src oauth2.TokenSource) {
	serverConn, err := minecraft.Dialer{
		TokenSource: src,
		ClientData:  conn.ClientData(),
	}.Dial("raknet", ip)
	if err != nil {
		panic(err)
	}
	var g sync.WaitGroup
	g.Add(2)
	go func() {
		if err := conn.StartGame(serverConn.GameData()); err != nil {
			panic(err)
		}
		fmt.Println("Started game")
		g.Done()
	}()
	go func() {
		if err := serverConn.DoSpawn(); err != nil {
			panic(err)
		}
		fmt.Println("Spawned")
		g.Done()
	}()
	g.Wait()

	cheats.LoadCheats()
	fmt.Println("Loaded cheats...")
	utils.MessageNotice("Use .help for a list of commands!", conn)
	err = conn.WritePacket(&packet.SetTitle{
		ActionType:      packet.TitleActionSetTitle,
		Text:            "§r§l§fPure§cZero\n§r§o§7Welcome...",
		FadeInDuration:  25,
		RemainDuration:  60,
		FadeOutDuration: 25,
	})
	if err != nil {
		return
	}
	go func() {
		for i := 0; i <= 4; i++ {
			err := conn.WritePacket(&packet.Text{
				TextType: packet.TextTypeJukeboxPopup,
				Message:  "§r§eLoading...",
			})
			if err != nil {
				return
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}()
	time.Sleep(6 * time.Second)
	utils.ClearInGameScreen(conn)

	go func() {
		defer func(listener *minecraft.Listener, conn *minecraft.Conn, message string) {
			err := listener.Disconnect(conn, message)
			if err != nil {
				panic(err)
			}
		}(listener, conn, "connection lost")
		defer func(serverConn *minecraft.Conn) {
			err := serverConn.Close()
			if err != nil {
				panic(err)
			}
		}(serverConn)
		for {
			pk, err := conn.ReadPacket()
			if err != nil {
				return
			}

			switch p := pk.(type) {
			case *packet.Text:
				if p.Message[0:1] == "." {
					if cheats.ParseMessage(p.Message, conn, serverConn) {
						continue
					}
				}
				break
			}

			if err := serverConn.WritePacket(pk); err != nil {
				if disconnect, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
					_ = listener.Disconnect(conn, disconnect.Error())
				}
				return
			}
		}
	}()
	go func() {
		defer func(serverConn *minecraft.Conn) {
			err := serverConn.Close()
			if err != nil {
				panic(err)
			}
		}(serverConn)
		defer func(listener *minecraft.Listener, conn *minecraft.Conn, message string) {
			err := listener.Disconnect(conn, message)
			if err != nil {
				panic(err)
			}
		}(listener, conn, "connection lost")
		for {
			pk, err := serverConn.ReadPacket()
			if err != nil {
				if disconnect, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
					_ = listener.Disconnect(conn, disconnect.Error())
				}
				return
			}

			switch p := pk.(type) {
			case *packet.AddPlayer:
				players.AddPlayer(
					p.Username,
					p.EntityRuntimeID,
					p.EntityMetadata,
					p.Position,
					p.AbilityData.EntityUniqueID,
				)
				break
			case *packet.SetActorData:
				player := players.GetPlayerByRuntimeID(p.EntityRuntimeID)
				if player != nil {
					player.MetaData = p.EntityMetadata
				}
				break
			case *packet.RemoveActor:
				players.RemovePlayerByUniqueID(p.EntityUniqueID)
				break
			case *packet.Transfer:
				err := conn.WritePacket(&packet.Transfer{
					Address: "127.0.0.1",
					Port:    19132,
				})
				if err != nil {
					return
				}
				proxy(fmt.Sprintf("%s:%s", p.Address, strconv.Itoa(int(p.Port))))
				continue
			}

			if err := conn.WritePacket(pk); err != nil {
				return
			}
		}
	}()
}
