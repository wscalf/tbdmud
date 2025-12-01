package scripting

import (
	"time"

	"github.com/dop251/goja"
	"github.com/wscalf/tbdmud/internal/game"
)

type PlayersWrapper struct {
	players *game.Players
	system  *GojaScriptSystem
}

func (p *PlayersWrapper) FindById(id string) *game.Player {
	return p.players.FindById(id)
}

func (p *PlayersWrapper) FindByName(name string) *game.Player {
	return p.players.FindByName(name)
}

func (p *PlayersWrapper) All() []*game.Player {
	return p.players.All()

}

func (p *PlayersWrapper) FindByNameIncludingOffline(name string) *goja.Promise {
	promise := p.system.createPromise()
	asyncContext := p.system.getCurrentCommandAsyncContext()
	asyncContext.SetPromise(promise)
	go func() {
		saveData, err := p.players.FindByNameIncludingOffline(name)
		if err != nil {
			asyncContext.SetError(err)
		} else {
			asyncContext.SetResult(saveData)
		}
	}()
	return promise.promise
}

type SystemWrapper struct {
	system *GojaScriptSystem
	sys    *game.System
}

func (s *SystemWrapper) Wait(seconds int) *goja.Promise {
	promise := s.system.createPromise()
	asyncContext := s.system.getCurrentCommandAsyncContext()
	asyncContext.SetPromise(promise)

	go func() {
		s.sys.Wait(time.Duration(seconds) * time.Second)
		asyncContext.SetResult(goja.Null())
	}()

	return promise.promise
}
