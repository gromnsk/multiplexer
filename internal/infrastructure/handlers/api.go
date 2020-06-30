package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

func (s *Server) MainHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		withTimeout, cancel := context.WithTimeout(r.Context(), s.cfg.Timeout)
		defer cancel()

		request := &Request{}
		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			failedResponse(err, w, http.StatusBadRequest)
			return
		}
		result, err := s.multiplexer.ProcessRequests(withTimeout, request.Urls)
		if err != nil {
			failedResponse(err, w, http.StatusBadRequest)
			return
		}
		response := Response{
			Data: result.Results,
		}
		data, err := json.Marshal(response)
		if err != nil {
			failedResponse(err, w, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(data)
		if err != nil {
			failedResponse(err, w, http.StatusInternalServerError)
			return
		}
	}
}
