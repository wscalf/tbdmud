package game

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	dataPath string
}

func (l *Loader) GetRooms() (map[string]*Room, error) {
	//Load all roomdata from YAML
	rooms := map[string]*Room{}
	roomsPath := filepath.Join(l.dataPath, "rooms")
	roomsData, err := loadRoomsDataFromAllFiles(roomsPath)
	if err != nil {
		return rooms, err
	}

	//In a first pass, construct room objects from each
	for _, roomData := range roomsData {
		room := NewRoom(roomData.ID, roomData.Name, roomData.Desc, nil)
		rooms[room.ID] = room
	}
	//In a second pass, process the linkdata from each roomdata and create links
	for _, roomData := range roomsData {
		for _, linkData := range roomData.Links {
			from := rooms[roomData.ID]
			to := rooms[linkData.To]

			from.Link(linkData.Command, linkData.Name, linkData.Desc, to)
		}
	}
	//Return
	return rooms, nil
}

func NewLoader(dataPath string) *Loader {
	return &Loader{
		dataPath: dataPath,
	}
}

func (l *Loader) GetMeta() (Metadata, error) {
	var meta Metadata
	worldFilePath := filepath.Join(l.dataPath, "world.yaml")
	f, err := os.Open(worldFilePath)
	if err != nil {
		return meta, err
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&meta)
	return meta, err
}

func loadRoomsDataFromAllFiles(roomsPath string) ([]RoomData, error) {
	roomsData := []RoomData{}
	err := filepath.WalkDir(roomsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(d.Name()) == ".yaml" {
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			defer f.Close()

			rooms, err := extractRoomsData(f)
			if err != nil {
				return err
			}

			roomsData = append(roomsData, rooms...)
		}

		return nil
	})

	return roomsData, err
}

func extractRoomsData(file io.Reader) ([]RoomData, error) {
	data := []RoomData{}
	decoder := yaml.NewDecoder(file)

	for {
		var room RoomData
		if err := decoder.Decode(&room); err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return data, err
			}
		}

		data = append(data, room)
	}

	return data, nil
}

type RoomData struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Desc  string `yaml:"desc"`
	Links []struct {
		Name    string `yaml:"name"`
		Command string `yaml:"cmd"`
		To      string `yaml:"to"`
		Desc    string `yaml:"desc"`
	} `yaml:"links"`
}

type Metadata struct {
	Banner      string `yaml:"banner"`
	ChargenRoom string `yaml:"chargen_room"`
	DefaultRoom string `yaml:"default_room"`
}
