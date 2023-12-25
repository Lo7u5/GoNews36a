package api

import (
	"GoNews36a/pkg/dbase/postgresql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type API struct {
	R  *mux.Router
	db *postgresql.Store
}

func New(db *postgresql.Store) *API {
	api := API{
		db: db,
	}
	api.R = mux.NewRouter()
	api.endpoints()
	return &api
}

func (api *API) endpoints() {
	api.R.HandleFunc("/news/{limit}", api.postsHandler).Methods(http.MethodGet)
	//подключение веб-интерфейса приложения
	api.R.PathPrefix("/").Handler(http.StripPrefix("/",
		http.FileServer(http.Dir("/Users/kv/GolandProjects/GoNews36a/cmd/gonews/webapp"))))
}

func (api *API) Router() *mux.Router {
	return api.R
}

// postsHandler получение массива новостей заданного размера
func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	s := mux.Vars(r)["limit"]
	limit, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	posts, err := api.db.Posts(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(posts)
}
