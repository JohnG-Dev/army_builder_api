package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

func LogError(s *state.State, r *http.Request, msg string, err error) {
	s.Logger.Error(
		msg,
		zap.Error(err),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
	)
}

func LogInfo(s *state.State, r *http.Request, msg string, fields ...zap.Field) {
	allFields := append(fields,
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
	)
	s.Logger.Info(
		msg,
		allFields...,
	)
}
