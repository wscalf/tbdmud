package game

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/wscalf/tbdmud/internal/text"
	"gopkg.in/yaml.v3"
)

type Loader struct {
	dataPath string
}

func (l *Loader) GetRooms(scriptSystem ScriptSystem, defaultRoomType string, defaultLinkType string, defaultObjectType string) (map[string]*Room, error) {
	//Load all roomdata from YAML
	rooms := map[string]*Room{}
	folder := filepath.Join(l.dataPath, "rooms")
	roomsData, err := loadRoomsDataFromAllFiles(folder)
	if err != nil {
		return rooms, err
	}

	//In a first pass, construct room objects from each
	for _, roomData := range roomsData {
		room := NewRoom(roomData.ID, roomData.Name, roomData.Desc, nil)
		script, err := scriptSystem.Wrap(room, selectName(roomData.TypeName, defaultRoomType))
		if err != nil {
			return nil, err
		}

		for key, val := range roomData.ScriptVars {
			script.Set(key, val)
		}

		for _, itemData := range roomData.Objects {
			item := NewObject(itemData.Name, itemData.Desc)

			script, err := scriptSystem.Wrap(item, selectName(itemData.TypeName, defaultObjectType))
			if err != nil {
				return nil, err
			}

			for key, val := range itemData.ScriptVars {
				script.Set(key, val)
			}

			item.AttachScript(script)

			room.AddItem(item)
		}
		room.AttachScript(script)
		rooms[room.ID] = room
	}
	//In a second pass, process the linkdata from each roomdata and create links
	for _, roomData := range roomsData {
		for _, linkData := range roomData.Links {
			from := rooms[roomData.ID]
			to := rooms[linkData.To]

			from.Link(linkData.Command, linkData.Name, linkData.Desc, to, scriptSystem, selectName(linkData.TypeName, defaultLinkType), linkData.ScriptVars)
		}
	}
	//Return
	return rooms, nil
}

func selectName(givenName, defaultName string) string {
	if givenName == "" {
		return defaultName
	}

	return givenName
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
	ID         string         `yaml:"id"`
	Name       string         `yaml:"name"`
	Desc       string         `yaml:"desc"`
	TypeName   string         `yaml:"type"`
	ScriptVars map[string]any `yaml:"vars"`
	Links      []struct {
		Name       string         `yaml:"name"`
		Command    string         `yaml:"cmd"`
		To         string         `yaml:"to"`
		TypeName   string         `yaml:"type"`
		ScriptVars map[string]any `yaml:"vars"`
		Desc       string         `yaml:"desc"`
	} `yaml:"links"`
	Objects []struct {
		Name       string         `yaml:"name"`
		Desc       string         `yaml:"desc"`
		TypeName   string         `yaml:"type"`
		ScriptVars map[string]any `yaml:"vars"`
	} `yaml:"objects"`
}

type Metadata struct {
	Banner            string `yaml:"banner"`
	ChargenRoom       string `yaml:"chargen_room"`
	DefaultRoom       string `yaml:"default_room"`
	DefaultObjectType string `yaml:"object_type"`
	DefaultPlayerType string `yaml:"player_type"`
	DefaultRoomType   string `yaml:"room_type"`
	DefaultLinkType   string `yaml:"link_type"`
}

func (l *Loader) GetLayouts() (map[string]*text.Layout, error) {
	layouts := map[string]*text.Layout{}
	folder := filepath.Join(l.dataPath, "layouts")
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) != ".tmpl" {
			return nil
		}

		name := strings.TrimSuffix(d.Name(), ".tmpl")
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		layout, err := text.NewLayout(name, string(body))
		if err != nil {
			return err
		}

		layouts[name] = layout
		return nil
	})
	if err != nil {
		return nil, err
	}

	return layouts, nil
}

func (l *Loader) ReadModuleTextFile(relativePath string) (string, error) {
	path := filepath.Join(l.dataPath, relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}
