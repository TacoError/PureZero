package cheats

import (
	"PureZero/players"
	"PureZero/utils"
	"github.com/sandertv/gophertunnel/minecraft"
	"strconv"
)

func hitBoxEnable(conn *minecraft.Conn, _ *minecraft.Conn) {
	width, err := strconv.ParseFloat(GetCheatSetting("HitBox", "width").Value, 32)
	if err != nil {
		utils.MessageError("The width setting must be a number.", conn)
		return
	}
	height, err := strconv.ParseFloat(GetCheatSetting("HitBox", "height").Value, 32)
	if err != nil {
		utils.MessageError("The height setting must be a number.", conn)
		return
	}
	players.SetHitBoxAll(float32(width), float32(height), conn)

}

func hitBoxDisable(conn *minecraft.Conn, _ *minecraft.Conn) {
	players.SetHitBoxAll(0.6, 2, conn)
}
