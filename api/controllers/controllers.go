package controllers

import (
	"myapp/internal/logic"

	"go.uber.org/zap"
)

type HandlerLayer struct {
	*logic.LogicLayer
	Logger *zap.SugaredLogger
}

func NewHandlerLayer(ll *logic.LogicLayer, logger *zap.SugaredLogger) *HandlerLayer {
	return &HandlerLayer{
		LogicLayer: ll,
		Logger:     logger,
	}
}
