package pager

import (
	"github.com/mrinjamul/authenticator-desktop/config"
	"github.com/mrinjamul/authenticator-desktop/constants"
	"github.com/mrinjamul/authenticator-desktop/pages"
)

type Pager interface {
	Start()
}

type pager struct {
	config config.Config
	pages  map[string]pages.Page
}

func (p *pager) Start() {
	p.launch(constants.PAGE_LAUNCHER_KEY)
}

func (p *pager) launch(key string) {
	p.pages[key].Render()
}

func Init(conf config.Config) Pager {
	p := preInitialize(conf)
	pages := []pages.Page{
		pages.NewLauncher(conf),
		pages.NewAddAccountPage(conf),
		pages.NewEditAccountPage(conf),
	}
	return postInitialize(p, pages)
}

func preInitialize(conf config.Config) *pager {
	p := &pager{
		config: conf,
		pages:  make(map[string]pages.Page),
	}
	conf.SetLauncher(p.launch)
	return p
}

func postInitialize(p *pager, pages []pages.Page) Pager {
	for _, page := range pages {
		p.pages[page.HashCode()] = page
	}
	return p
}
