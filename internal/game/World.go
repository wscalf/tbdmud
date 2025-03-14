package game

import "github.com/wscalf/tbdmud/internal/text"

type World struct {
	chargen           *Room
	defaultRoom       *Room
	rooms             map[string]*Room
	defaultRoomLayout *text.Layout
}

func NewWorld() *World {
	return &World{
		chargen:     nil,
		defaultRoom: nil,
		rooms:       map[string]*Room{},
	}
}

func (w *World) InitializeRooms(rooms map[string]*Room) {
	if w.defaultRoomLayout != nil {
		for _, r := range rooms {
			r.layout = w.defaultRoomLayout
		}
	}
	w.rooms = rooms
}

func (w *World) AddRoom(r *Room) {
	if w.defaultRoomLayout != nil {
		r.layout = w.defaultRoomLayout
	}
	w.rooms[r.ID] = r
}

func (w *World) SetRoomLayout(layout *text.Layout) {
	for _, r := range w.rooms {
		r.layout = w.defaultRoomLayout
	}

	w.defaultRoomLayout = layout
}

func (w *World) SetSystemRooms(chargen string, defaultRoom string) {
	w.chargen = w.rooms[chargen]
	w.defaultRoom = w.rooms[defaultRoom]
}
