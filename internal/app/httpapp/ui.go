package httpapp

import (
	"de/internal/storage/sqlstorage"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func handleIndexPage(store *sqlstorage.Store) http.HandlerFunc {
	type data struct {
		Error            string
		Accounts         []sqlstorage.Account
		From, To, Amount uint64
	}

	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(
			"templates/base.tmpl.html",
			"templates/transfer.tmpl.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		accs, err := store.ListAccounts(r.Context(), 10, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		qp := r.URL.Query()
		from, _ := strconv.ParseUint(qp.Get("from"), 10, 64)
		to, _ := strconv.ParseUint(qp.Get("to"), 10, 64)
		amount, _ := strconv.ParseUint(qp.Get("amount"), 10, 64)

		w.Header().Add("Content-Type", "text/html")
		if err := t.Execute(w, data{
			Accounts: accs,
			Error:    r.URL.Query().Get("error"),
			From:     from,
			To:       to,
			Amount:   amount,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleIsolationPage(store *sqlstorage.Store) http.HandlerFunc {
	type radioButton struct {
		Value    string
		Text     string
		Checked  bool
		Disabled bool
	}

	var radioButtons = [...]radioButton{
		{Value: "1", Text: "Dirty Read", Checked: true},
		{Value: "2", Text: "Non Repeatable Read"},
		{Value: "3", Text: "Phantom Read"},
		{Value: "4", Text: "Lost Update (Mysql default uses Row Level Locking)", Disabled: true},
	}

	type tdata struct {
		Error string
		RBtns []radioButton
	}

	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(
			"templates/base.tmpl.html",
			"templates/isolation.tmpl.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "text/html")
		if err := t.Execute(w, tdata{
			RBtns: radioButtons[:],
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleRefreshDB(store *sqlstorage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := store.Refresh(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func handleIndexingPage(store *sqlstorage.Store) http.HandlerFunc {
	type q struct {
		A string
		Q string
	}

	type tdata struct {
		Error   string
		Queries []q
		Count   uint64
	}

	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(
			"templates/base.tmpl.html",
			"templates/indexing.tmpl.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()

		all, err := store.AnalyzeAllColumnSelect(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		indexOnly, err := store.AnalyzeIndexedColumnSelect(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		idOnly, err := store.AnalyzePrimaryKeySelect(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		unindexedOnly, err := store.AnalyzeUnindexedColumnSelect(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		allExplicit, err := store.AnalyzeAllExplicitIndexSelect(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		idName, err := store.AnalyzePrimaryKeyPlusIndexSelect(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pagLOID, err := store.AnalyzePaginationLimitOffset(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pagCID, err := store.AnalyzePaginationCursor(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pagLOIDN, err := store.AnalyzePaginationLOName(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pagCIDN, err := store.AnalyzePaginationCName(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pagLOIDNN, err := store.AnalyzePaginationLONameName2(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pagCIDNN, err := store.AnalyzePaginationCNameName2(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pagLOAll, err := store.AnalyzePaginationLOAll(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pagCAll, err := store.AnalyzePaginationCAll(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		count, err := store.CountEmployees(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "text/html")
		if err := t.Execute(w, tdata{
			Count: count,
			Queries: []q{
				{Q: "*", A: all},
				{Q: "name2 (no index)", A: unindexedOnly},
				{Q: "name (index)", A: indexOnly},
				{Q: "id (pk)", A: idOnly},
				{Q: "id + name (f:id)", A: idName},
				{Q: "id + name + name2 (f:id)", A: allExplicit},
				{Q: "paginate id (limit,offset)", A: pagLOID},
				{Q: "paginate id (cursor)", A: pagCID},
				{Q: "paginate id + name (limit,offset)", A: pagLOIDN},
				{Q: "paginate id + name (cursor)", A: pagCIDN},
				{Q: "paginate id + name + name2 (limit,offset)", A: pagLOIDNN},
				{Q: "paginate id + name + name2 (cursor)", A: pagCIDNN},
				{Q: "paginate * (limit,offset)", A: pagLOAll},
				{Q: "paginate * (cursor)", A: pagCAll},
			},
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
}
