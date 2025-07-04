package domain

import (
	"context"
	"github.com/dan-lugg/go-commands/commands"
)

type NestedCommandRes struct {
	Result int `json:"result"`
}

type NestedCommandReq struct {
	commands.CommandReq[NestedCommandRes]
	ArgX int `json:"argX"`
	ArgY int `json:"argY"`
}

type NestedHandler struct {
	commands.Handler[NestedCommandReq, NestedCommandRes]
	handlerRegistry *commands.HandlerCatalog
}

func NewNestedHandler(registry *commands.HandlerCatalog) *NestedHandler {
	return &NestedHandler{
		handlerRegistry: registry,
	}
}

func (h *NestedHandler) Handle(req NestedCommandReq, ctx context.Context) (res NestedCommandRes, err error) {
	addRes, err := h.handlerRegistry.Handle(AddCommandReq{
		ArgX: req.ArgX,
		ArgY: req.ArgY,
	}, ctx)
	if err != nil {
		return NestedCommandRes{}, err
	}

	subRes, err := h.handlerRegistry.Handle(SubCommandReq{
		ArgX: req.ArgX,
		ArgY: req.ArgY,
	}, ctx)
	if err != nil {
		return NestedCommandRes{}, err
	}

	return NestedCommandRes{
		Result: (addRes.(AddCommandRes)).Result * (subRes.(SubCommandRes)).Result,
	}, nil
}
