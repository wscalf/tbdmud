package game

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/wscalf/tbdmud/internal/game/parameters"
)

type Login struct {
	banner string
	//Temporary placeholder for data that would be kept in the DB
	accounts   map[string]*Account
	characters map[string]*Player
}

func NewLogin(banner string) *Login {
	return &Login{
		banner:     banner,
		accounts:   map[string]*Account{},
		characters: map[string]*Player{},
	}
}

func (l *Login) Process(client Client) *Player {
	client.Send(l.banner)
	account := l.loginOrRegister(client)
	if account == nil {
		return nil
	}

	character := l.selectOrCreateCharacter(client, account)
	if character == nil {
		return nil
	}

	return character
}

func (l *Login) loginOrRegister(client Client) *Account {
	const prompt string = "Use `login <username> <password>` to log in, `register <username> <password>` to register, `help` to repeat this message, or `quit` to disconnect."
	paramspec := []parameters.Parameter{parameters.NewName("login", true), parameters.NewName("password", true)}
	client.Send(prompt)
	for input := range client.Recv() {
		cmd, argpart := SplitCommandNameFromArgs(input)

		switch cmd {
		case "login", "register":
			params, err := ExtractParameters(cmd, argpart, paramspec)
			if err != nil {
				client.Send(err.Error())
				continue
			}

			var account *Account
			if cmd == "login" {
				account, err = l.tryLogin(params["login"], params["password"])
			} else {
				account, err = l.tryRegister(params["login"], params["password"])
			}
			if err != nil {
				client.Send(err.Error())
				continue
			}

			return account
		case "help":
			client.Send(prompt)
		case "quit":
			return nil
		default:
			client.Send("unrecognized command")
			client.Send(prompt)
		}
	}

	return nil
}

func (l *Login) tryLogin(login string, password string) (*Account, error) {
	account, ok := l.accounts[login] //Simulates looking up the account in the DB
	if !ok {
		return nil, fmt.Errorf("incorrect login or password")
	}

	if !account.CheckPassword(password) {
		return nil, fmt.Errorf("incorrect login or password")
	}

	return account, nil
}

func (l *Login) tryRegister(login string, password string) (*Account, error) {
	_, exists := l.accounts[login] //Simulates checking for the account in the DB
	if exists {
		return nil, fmt.Errorf("account already exists")
	}

	account := NewAccount(login)
	account.SetPassword(password)
	l.accounts[login] = account //Simulates adding the account to the DB

	return account, nil
}

func (l *Login) selectOrCreateCharacter(client Client, account *Account) *Player {
	characters := account.characters
	client.Send("Select a character:")
	for i, entry := range characters {
		client.Send(fmt.Sprintf("%d: %s", i+1, entry.name))
	}
	client.Send("Or use create <name> to create a new character. (Remember to quote full names.)")

	for input := range client.Recv() {
		i, err := strconv.Atoi(input)
		if err == nil {
			selection := i - 1
			if selection < 0 || selection >= len(characters) {
				client.Send("Please select one of the above options.")
				continue
			}

			entry := account.characters[selection]
			character := l.characters[entry.id] //Simulates looking up the character in the DB
			return character
		}

		//Should be a new character then
		paramspec := []parameters.Parameter{parameters.NewName("name", true)}
		cmd, argpart := SplitCommandNameFromArgs(input)
		if cmd != "create" {
			client.Send("Please select one of the above options.")
			continue
		}
		params, err := ExtractParameters(cmd, argpart, paramspec)
		if err != nil {
			client.Send(err.Error())
			continue
		}

		name := params["name"]
		id := uuid.NewString()
		character := NewPlayer(id, name)
		account.AddCharacter(character) //Changes account which would need to be written back to the DB
		l.characters[id] = character    //Simulates saving the new character to the DB
		return character
	}

	return nil //Client disconnected partway through
}
