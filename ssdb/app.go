package ssdb

import (
	"github.com/siddontang/golib/leveldb"
	"net"
	"strings"
)

type App struct {
	cfg *Config

	listener net.Listener

	db *leveldb.DB

	kvTx   *tx
	listTx *tx
	hashTx *tx
	zsetTx *tx
}

func NewApp(cfg *Config) (*App, error) {
	app := new(App)

	app.cfg = cfg

	var err error

	if strings.Contains(cfg.Addr, "/") {
		app.listener, err = net.Listen("unix", cfg.Addr)
	} else {
		app.listener, err = net.Listen("tcp", cfg.Addr)
	}

	if err != nil {
		return nil, err
	}

	app.db, err = leveldb.OpenWithConfig(&cfg.DB)
	if err != nil {
		return nil, err
	}

	app.kvTx = app.newTx()
	app.listTx = app.newTx()
	app.hashTx = app.newTx()
	app.zsetTx = app.newTx()

	return app, nil
}

func (app *App) Close() {
	app.listener.Close()

	app.db.Close()
}

func (app *App) Run() {
	for {
		conn, err := app.listener.Accept()
		if err != nil {
			continue
		}

		newClient(conn, app)
	}
}
