package cheats

import "github.com/sandertv/gophertunnel/minecraft"

var cheats []*Cheat

func LoadCheats() {
	newCheat("Example", "Example cheat", []*Setting{
		newSetting("cool", "Example setting", "Default value"),
	}, exampleEnable, exampleDisable)
	makeEffectCheat("Haste", "Give yourself haste", 3)
}

func newCheat(
	name string,
	description string,
	settings []*Setting,
	enable func(conn *minecraft.Conn, serverConn *minecraft.Conn),
	disable func(conn *minecraft.Conn, serverConn *minecraft.Conn),
) {
	cheats = append(cheats, &Cheat{
		Name:        name,
		Description: description,
		Settings:    settings,
		Enable:      enable,
		Disable:     disable,
		Enabled:     false,
	})
}

func newSetting(name string, description string, defaultValue string) *Setting {
	return &Setting{
		Name:        name,
		Description: description,
		Value:       defaultValue,
	}
}
