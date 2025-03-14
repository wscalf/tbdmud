package game

type World struct {
	chargen     *Room
	defaultRoom *Room
	rooms       map[string]*Room
}

func NewWorld() *World {
	return &World{
		chargen:     nil,
		defaultRoom: nil,
		rooms:       map[string]*Room{},
	}
}

func (w *World) InitializeRooms(rooms map[string]*Room) {
	w.rooms = rooms
}

func (w *World) AddRoom(r *Room) {
	w.rooms[r.ID] = r
}

func (w *World) SetSystemRooms(chargen string, defaultRoom string) {
	w.chargen = w.rooms[chargen]
	w.defaultRoom = w.rooms[defaultRoom]
}
