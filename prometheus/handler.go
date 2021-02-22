package prometheus

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

type handler struct {
	prometheusLite *Lite
	mux            *http.ServeMux
}

// Handler returns an http handler exposing Prometheus APIs
func (p *Lite) Handler() http.Handler {
	h := handler{prometheusLite: p}
	mux := http.NewServeMux()
	mux.HandleFunc("/query_range", h.QueryRange)
	h.mux = mux
	return &h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// QueryRange query promql between x and y time
func (h *handler) QueryRange(w http.ResponseWriter, r *http.Request) {

	start, err := parseTime(r.FormValue("start"))
	if err != nil {
		respond(w, http.StatusBadRequest, err.Error())
		return
	}

	end, err := parseTime(r.FormValue("end"))
	if err != nil {
		respond(w, http.StatusBadRequest, err.Error())
		return
	}

	step, err := time.ParseDuration(r.FormValue("step"))
	if err != nil {
		respond(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.prometheusLite.QueryRange(r.Context(), r.FormValue("query"), time.Time(start), time.Time(end), step)
	if err != nil {
		respond(w, http.StatusBadRequest, fmt.Sprintf("failed to execute query: %v", err))
		return
	}

	d, err := json.Marshal(res)
	if err != nil {
		respond(w, http.StatusInternalServerError, fmt.Sprintf("failed to serialize result: %v", err))
		return
	}

	respond(w, http.StatusOK, string(d))
}

func respond(resp http.ResponseWriter, status int, body string) {
	resp.WriteHeader(status)
	resp.Write([]byte(body))
}

func parseTime(s string) (time.Time, error) {
	if t, err := strconv.ParseFloat(s, 64); err == nil {
		s, ns := math.Modf(t)
		ns = math.Round(ns*1000) / 1000
		return time.Unix(int64(s), int64(ns*float64(time.Second))).UTC(), nil
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("cannot parse %q to a valid timestmap", s)
}
