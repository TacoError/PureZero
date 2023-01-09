package cheats

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"strconv"
)

func makeEffectCheat(name string, description string, id int) {
	cheats = append(cheats, &Cheat{
		Name:        name,
		Description: description,
		Enabled:     false,
		Settings: []*Setting{
			newSetting("amplifier", "Amplifier of effect", "1"),
		},
		Enable: func(conn *minecraft.Conn, _ *minecraft.Conn) {
			amplifier := GetCheatSetting(name, "amplifier")
			amp := 1
			if amplifier != nil {
				amp, _ = strconv.Atoi(amplifier.Value)
			}
			err := conn.WritePacket(&packet.MobEffect{
				EntityRuntimeID: conn.GameData().EntityRuntimeID,
				Operation:       packet.MobEffectAdd,
				EffectType:      int32(id),
				Amplifier:       int32(amp),
				Duration:        2147483647,
			})
			if err != nil {
				return
			}
		},
		Disable: func(conn *minecraft.Conn, _ *minecraft.Conn) {
			err := conn.WritePacket(&packet.MobEffect{
				EntityRuntimeID: conn.GameData().EntityRuntimeID,
				Operation:       packet.MobEffectRemove,
				EffectType:      int32(id),
			})
			if err != nil {
				return
			}
		},
	})
}
