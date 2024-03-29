package httpapp

import (
	"de/internal/storage/sqlstorage"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func handleTransfer(db *sqlstorage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := parseTransferReq(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch req.Type {
		case atomicTransfer:
			err = db.AtomicTransfer(r.Context(), req.From, req.To, req.Amount)
		case nonAtomicTransfer:
			err = db.NonAtomicTransfer(r.Context(), req.From, req.To, req.Amount)
		case atomicFailedTransfer:
			_ = db.FailedAtomicTransfer(r.Context(), req.From, req.To, req.Amount)
		case nonAtomicFailedTransfer:
			_ = db.FailedNonAtomicTransfer(r.Context(), req.From, req.To, req.Amount)
		}

		if err != nil {
			http.Error(w, "transfer failed", http.StatusUnprocessableEntity)
			log.Printf("transfer failed: %v", err)
			return
		}

		qp := url.Values{
			"from":   []string{strconv.FormatUint(req.From, 10)},
			"to":     []string{strconv.FormatUint(req.To, 10)},
			"amount": []string{strconv.FormatUint(req.Amount, 10)},
		}
		http.Redirect(w, r, "/?"+qp.Encode(), http.StatusFound)
	}
}

type transferReq struct {
	Type             transferType
	From, To, Amount uint64
}

type transferType uint64

const (
	atomicTransfer transferType = iota + 1
	nonAtomicTransfer
	atomicFailedTransfer
	nonAtomicFailedTransfer
	unknownTransfer
)

func parseTransferReq(r *http.Request) (transferReq, error) {
	from, err := strconv.ParseUint(r.FormValue("from"), 10, 64)
	if err != nil {
		return transferReq{}, fmt.Errorf("parse from: %v", err)
	}

	to, err := strconv.ParseUint(r.FormValue("to"), 10, 64)
	if err != nil {
		return transferReq{}, fmt.Errorf("parse to: %v", err)
	}

	amount, err := strconv.ParseUint(r.FormValue("amount"), 10, 64)
	if err != nil {
		return transferReq{}, fmt.Errorf("parse amount: %v", err)
	}

	type_, err := strconv.ParseUint(r.FormValue("type"), 10, 64)
	if err != nil {
		return transferReq{}, fmt.Errorf("parse transfer type: %v", err)
	}

	ttype := transferType(type_)

	if ttype < atomicTransfer || ttype >= unknownTransfer {
		return transferReq{}, fmt.Errorf("invalid transfer type: %d", type_)
	}

	return transferReq{
		From:   from,
		To:     to,
		Amount: amount,
		Type:   ttype,
	}, nil
}
