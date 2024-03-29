package httpapp

import (
	"de/internal/core"
	"de/internal/storage/sqlstorage"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func handleIsolation(store *sqlstorage.Store) http.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	type simResult struct {
		TxID  string            `json:"tx"`
		Query string            `json:"query"`
		Sales []sqlstorage.Sale `json:"sales"`
	}
	var message struct {
		SimType  uint64 `json:"type"`
		Quantity uint64 `json:"qty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		defer conn.Close()

		if err := conn.ReadJSON(&message); err != nil {
			log.Println(err)
			return
		}

		var states []core.SaleSimulation
		ctx := r.Context()
		switch message.SimType {
		case 1:
			states, err = core.SimulateDirtyRead(ctx, store)
		case 2:
			states, err = core.SimulateNonRepeatableRead(ctx, store)
		case 3:
			states, err = core.SimulatePhantomRead(ctx, store)
		case 4:
			log.Println(4)
			states, err = core.SimulateLostUpdates(ctx, store)
		default:
			return
		}
		if err != nil {
			log.Println(err)
			return
		}

		if err := store.Refresh(ctx); err != nil {
			log.Println(err)
			return
		}

		for _, st := range states {
			if err := ctx.Err(); err != nil {
				return
			}

			conn.WriteJSON(st)
			time.Sleep(time.Second)
		}
	}
}
