package players

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Player struct {
	Name            string
	EntityRuntimeID uint64
	// UniqueID is used solely for the purpose of removing the player in the
	// RemoveActor packet
	UniqueID int64
	// MetaData can be used to make some cool changes to players
	// by changing the values and sending a SetActorData packet
	// https://github.com/Sandertv/gophertunnel/blob/master/minecraft/protocol/entity_metadata.go
	MetaData map[uint32]any
	Position mgl32.Vec3
}

// SetHitBox - an Update will need to be called for this to take place
func (p *Player) SetHitBox(horizontal, vertical float32) {
	p.MetaData[protocol.EntityDataKeyWidth] = horizontal
	p.MetaData[protocol.EntityDataKeyHeight] = vertical
}

func (p *Player) Update(conn *minecraft.Conn) {
	err := conn.WritePacket(&packet.SetActorData{
		EntityRuntimeID: p.EntityRuntimeID,
		EntityMetadata:  p.MetaData,
		Tick:            0,
	})
	if err != nil {
		return
	}
}
