package game

import "github.com/wscalf/tbdmud/internal/text"

type World struct {
	chargen           *Room
	defaultRoom       *Room
	rooms             map[string]*Room
	defaultRoomLayout *text.Layout
	scriptSystem      ScriptSystem
	defaultRoomType   string
}

func NewWorld(scriptSystem ScriptSystem, defaultRoomType string) *World {
	return &World{
		chargen:         nil,
		defaultRoom:     nil,
		rooms:           map[string]*Room{},
		scriptSystem:    scriptSystem,
		defaultRoomType: defaultRoomType,
	}
}

func (w *World) InitializeRooms(rooms map[string]*Room) {
	for _, r := range rooms {
		if w.defaultRoomLayout != nil {
			r.layout = w.defaultRoomLayout
		}

		script, _ := w.scriptSystem.Wrap(r, w.defaultRoomType)
		r.AttachScript(script)
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
