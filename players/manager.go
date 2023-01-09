package players

import "github.com/go-gl/mathgl/mgl32"

var players []*Player

func AddPlayer(
	name string,
	EntityRuntimeID uint64,
	MetaData map[uint32]any,
	position mgl32.Vec3,
	UniqueID int64,
) {
	players = append(players, &Player{
		Name:            name,
		EntityRuntimeID: EntityRuntimeID,
		MetaData:        MetaData,
		Position:        position,
		UniqueID:        UniqueID,
	})
}

func GetPlayerByRuntimeID(id uint64) *Player {
	for _, player := range players {
		if player.EntityRuntimeID == id {
			return player
		}
	}
	return nil
}

func RemovePlayerByEntityRuntimeID(id uint64) {
	var newArray []*Player
	for _, player := range players {
		if player.EntityRuntimeID == id {
			continue
		}
		newArray = append(newArray, player)
	}
	players = newArray
}

func RemovePlayerByUniqueID(id int64) {
	var newArray []*Player
	for _, player := range players {
		if player.UniqueID == id {
			continue
		}
		newArray = append(newArray, player)
	}
	players = newArray
}

func RemovePlayerByName(name string) {
	var newArray []*Player
	for _, player := range players {
		if player.Name == name {
			continue
		}
		newArray = append(newArray, player)
	}
	players = newArray
}
