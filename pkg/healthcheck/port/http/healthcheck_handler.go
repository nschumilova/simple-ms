package healthcheck

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type HealthView struct {
	Status string `json:"status"`
}

type HeathCheckHandler struct {
	Log *zap.SugaredLogger
}

func (h *HeathCheckHandler) Health(w http.ResponseWriter, r *http.Request) {
	health := HealthView{
		Status: "OK",
	}
	data, err := json.Marshal(&health)
	if err != nil {
		h.Log.Errorw("cannot marshal health to json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response building problem."))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
