package game

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/wscalf/tbdmud/internal/game/parameters"
	"github.com/wscalf/tbdmud/internal/text"
)

var ErrIncorrectInput error = errors.New("invalid input")

type Login struct {
	banner              string
	defaultPlayerLayout *text.Layout
	store               Storage
}

func NewLogin(banner string, playerLayout *text.Layout, store Storage) *Login {
	return &Login{
		banner:              banner,
		defaultPlayerLayout: playerLayout,
		store:               store,
	}
}

func (l *Login) Process(client Client) (*PlayerSaveData, error) {
	sendSynchronous(client, l.banner)
	account, err := l.loginOrRegister(client)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, nil
	}

	character, err := l.selectOrCreateCharacter(client, account)
	if err != nil {
		return nil, err
	}
	if character == nil {
		return nil, nil
	}

	return character, nil
}

func (l *Login) loginOrRegister(client Client) (*Account, error) {
	const prompt string = "Use `login <username> <password>` to log in, `register <username> <password>` to register, `help` to repeat this message, or `quit` to disconnect."
	paramspec := []parameters.Parameter{parameters.NewName("login", true), parameters.NewName("password", true)}
	sendSynchronous(client, prompt)
	for input := range client.Recv() {
		cmd, argpart := SplitCommandNameFromArgs(input)

		switch cmd {
		case "login", "register":
			params, err := ExtractParameters(cmd, argpart, paramspec)
			if err != nil {
				sendSynchronous(client, err.Error())
				continue
			}

			var account *Account
			if cmd == "login" {
				account, err = l.tryLogin(params["login"], params["password"])
			} else {
				account, err = l.tryRegister(params["login"], params["password"])
			}
			if err != nil {
				if errors.Is(err, ErrIncorrectInput) {
					sendSynchronous(client, err.Error())
					continue
				} else {
					return nil, err
				}
			}

			return account, nil
		case "help":
			sendSynchronous(client, prompt)
		case "quit":
			return nil, nil
		default:
			sendSynchronous(client, "unrecognized command")
			sendSynchronous(client, prompt)
		}
	}

	return nil, nil
}

func (l *Login) tryLogin(login string, password string) (*Account, error) {
	account, err := l.store.FindAccount(login)
	if err != nil {
		return nil, fmt.Errorf("error loading account %s: %w", account, err)
	}

	if account == nil {
		return nil, fmt.Errorf("%w: incorrect login or password", ErrIncorrectInput)
	}

	if !account.CheckPassword(password) {
		return nil, fmt.Errorf("%w: incorrect login or password", ErrIncorrectInput)
	}

	return account, nil
}

func (l *Login) tryRegister(login string, password string) (*Account, error) {
	existing, err := l.store.FindAccount(login)
	if err != nil {
		return nil, fmt.Errorf("error loading account: %w", err)
	}

	if existing != nil {
		return nil, fmt.Errorf("%w: account already exists", ErrIncorrectInput)
	}

	account := NewAccount(login)
	account.SetPassword(password)
	err = l.store.CreateOrUpdateAccount(account)
	if err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}

	return account, nil
}

func (l *Login) selectOrCreateCharacter(client Client, account *Account) (*PlayerSaveData, error) {
	characters := account.characters
	sendSynchronous(client, "Select a character:")
	for i, entry := range characters {
		sendSynchronous(client, "%d: %s", i+1, entry.name)
	}
	sendSynchronous(client, "Or use create <name> to create a new character. (Remember to quote full names.)")

	for input := range client.Recv() {
		i, err := strconv.Atoi(input)
		if err == nil {
			selection := i - 1
			if selection < 0 || selection >= len(characters) {
				sendSynchronous(client, "Please select one of the above options.")
				continue
			}

			entry := account.characters[selection]
			character, err := l.store.FindPlayer(entry.id)
			if err != nil {
				return nil, fmt.Errorf("error loading character: %w", err)
			}
			return character, nil
		}

		//Should be a new character then
		paramspec := []parameters.Parameter{parameters.NewName("name", true)}
		cmd, argpart := SplitCommandNameFromArgs(input)
		if cmd != "create" {
			sendSynchronous(client, "Please select one of the above options.")
			continue
		}
		params, err := ExtractParameters(cmd, argpart, paramspec)
		if err != nil {
			sendSynchronous(client, err.Error())
			continue
		}

		name := params["name"]
		id := uuid.NewString()
		character := &PlayerSaveData{ObjectSaveData: ObjectSaveData{ID: id, Name: name}}
		account.AddCharacter(character)

		err = l.store.CreateOrUpdatePlayer(character)
		if err != nil {
			return nil, fmt.Errorf("error saving character: %w", err)
		}

		err = l.store.CreateOrUpdateAccount(account)
		if err != nil {
			return nil, fmt.Errorf("error adding character to account: %w", err)
		}

		return character, nil
	}

	return nil, nil
}

func sendSynchronous(client Client, template string, params ...interface{}) error {
	job := text.NewPrintfJob(template+"\n", params...)
	return client.Send(job)
}
