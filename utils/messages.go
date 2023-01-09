package utils

import (
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func SendMessage(message string, conn *minecraft.Conn) {
	err := conn.WritePacket(&packet.Text{
		TextType: packet.TextTypeChat,
		Message:  message,
	})
	if err != nil {
		return
	}
}

func ClearInGameScreen(conn *minecraft.Conn) {
	for i := 0; i <= 100; i++ {
		SendMessage("\n§c §d §e\n", conn)
	}
}

func MessageNotice(message string, conn *minecraft.Conn) {
	SendMessage(fmt.Sprintf("§r§7[§cPure§fZero§7] §b%s", message), conn)
}

func MessageWarning(message string, conn *minecraft.Conn) {
	SendMessage(fmt.Sprintf("§r§7[§cPure§fZero§7] §e%s", message), conn)
}

func MessageError(message string, conn *minecraft.Conn) {
	SendMessage(fmt.Sprintf("§r§7[§cPure§fZero§7] §c%s", message), conn)
}
