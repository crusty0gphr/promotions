package promotions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	routRoot       = "/"
	routPromotions = "/promotions"
	reoutUpdate    = "/update"
)

const (
	keyNameID = "id"
	keyID     = "/{" + keyNameID + "}"
)

func MakeHandlers(logger *log.Logger, mux *mux.Router, addr string, service Service) error {
	logger.Printf("server started %s", addr)

	mux.HandleFunc(routRoot, makeRootHandler()).Methods(http.MethodGet)
	mux.HandleFunc(routPromotions+keyID, makeGetPromotionHandler(logger, service)).Methods(http.MethodGet)
	mux.HandleFunc(reoutUpdate, makeUpdateHandler(logger, service)).Methods(http.MethodPost)

	return http.ListenAndServe(addr, mux)
}

func makeRootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("server is ready"))
	}
}

func makeGetPromotionHandler(logger *log.Logger, service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("got request: %s %s", r.Method, r.URL.Path)

		vars := mux.Vars(r)
		idStr := vars[keyNameID]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusInternalServerError)
			return
		}

		res, err := service.getOne(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	}
}

func makeUpdateHandler(logger *log.Logger, service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("got request: %s %s", r.Method, r.URL.Path)

		if err := service.updateData(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("updated"))
	}
}
