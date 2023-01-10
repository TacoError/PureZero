package cheats

import (
	"PureZero/utils"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft"
	"strings"
)

func update(cheat *Cheat, conn *minecraft.Conn) {
	if cheat.Enabled {
		utils.MessageNotice(fmt.Sprintf("Enabled cheat: %s", cheat.Name), conn)
		return
	}
	utils.MessageNotice(fmt.Sprintf("Disabled cheat: %s", cheat.Name), conn)
}

func generateSettingsString(settings []*Setting) string {
	var str []string
	for _, setting := range settings {
		str = append(str, fmt.Sprintf(" §7- §e%s: §f%s", setting.Name, setting.Description))
	}
	return strings.Join(str, "\n")
}

func ParseMessage(message string, conn *minecraft.Conn, serverConn *minecraft.Conn) bool {
	message = strings.ToLower(message)[1:]
	messageSplit := strings.Split(message, " ")
	message = messageSplit[0]
	messageSplit = messageSplit[1:]

	if message == "list" {
		var list []string
		for _, cheat := range cheats {
			list = append(list, cheat.Name)
		}
		utils.MessageNotice(fmt.Sprintf("Cheats list: %s", strings.Join(list, ", ")), conn)
		return true
	}

	if message == "help" {
		if len(messageSplit) < 1 {
			utils.SendMessage("§r§7[§cPure§fZero§7] §r§fPlease use §e.help (cheat) §7or §e.list §ffor a list of cheats.", conn)
			return true
		}
		c := GetCheat(messageSplit[0])
		if c == nil {
			utils.MessageWarning("There is no cheat with that name.", conn)
			return true
		}
		info := []string{
			fmt.Sprintf("§r§7[§cPure§fZero§7] §f(%s)", c.Name),
			fmt.Sprintf("§eDescription§7: §f%s", c.Description),
			"§eSettings§7: §f",
			generateSettingsString(c.Settings),
		}
		utils.SendMessage(strings.Join(info, "\n"), conn)
		return true
	}

	if len(messageSplit) < 1 {
		return true
	}
	cheat := GetCheat(message)
	if cheat == nil {
		return false
	}
	if messageSplit[0] == "toggle" {
		cheat.Enabled = !cheat.Enabled
		if cheat.Enabled {
			cheat.Enable(conn, serverConn)
		} else {
			cheat.Disable(conn, serverConn)
		}
		update(cheat, conn)
		return true
	}
	s := GetCheatSetting(cheat.Name, messageSplit[0])
	if s == nil {
		utils.MessageWarning("Invalid setting name, please run .help (cheat) for a list of settings.", conn)
		return true
	}
	if len(messageSplit) < 2 {
		utils.MessageWarning("Please provide a new value for the setting.", conn)
		return true
	}
	old := s.Value
	s.Value = messageSplit[1]
	utils.MessageNotice(fmt.Sprintf("Value changed (%s -> %s)", old, s.Value), conn)
	return true
}
