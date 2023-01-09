package cheats

import "strings"

type Setting struct {
	Name        string
	Description string
	Value       string
}

func GetCheatSetting(cheat string, name string) *Setting {
	cheat = strings.ToLower(cheat)
	name = strings.ToLower(name)
	c := GetCheat(cheat)
	if c == nil {
		return nil
	}
	for _, setting := range c.Settings {
		if strings.ToLower(setting.Name) == name {
			return setting
		}
	}
	return nil
}
