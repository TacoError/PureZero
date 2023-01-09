package cheats

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"strings"
)

type Cheat struct {
	Name        string
	Description string
	Enable      func(conn *minecraft.Conn, serverConn *minecraft.Conn)
	Disable     func(conn *minecraft.Conn, serverConn *minecraft.Conn)
	Settings    []*Setting
	Enabled     bool
}

func GetCheat(name string) *Cheat {
	name = strings.ToLower(name)
	for _, cheat := range cheats {
		if strings.ToLower(cheat.Name) == name {
			return cheat
		}
	}
	return nil
}
