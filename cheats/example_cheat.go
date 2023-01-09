package cheats

import (
	"PureZero/utils"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft"
)

func exampleEnable(conn *minecraft.Conn, _ *minecraft.Conn) {
	value := GetCheatSetting("Example", "cool").Value
	utils.MessageNotice(fmt.Sprintf("(Enabled) value: %s", value), conn)
}

func exampleDisable(conn *minecraft.Conn, _ *minecraft.Conn) {
	value := GetCheatSetting("Example", "cool").Value
	utils.MessageNotice(fmt.Sprintf("(Disabled) value: %s", value), conn)
}
