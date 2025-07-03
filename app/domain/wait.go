package domain

import (
	"context"
	"github.com/dan-lugg/go-commands/commands"
	"time"
)

type WaitCommandRes struct{}

type WaitCommandReq struct {
	commands.CommandReq[WaitCommandRes]
	Duration int `json:"duration"`
}

type WaitHandler struct {
	commands.Handler[WaitCommandReq, WaitCommandRes]
}

func (h *WaitHandler) Handle(req WaitCommandReq, ctx context.Context) (res WaitCommandRes, err error) {
	for duration := req.Duration; duration > 0; duration-- {
		select {
		case <-ctx.Done():
			println("Context cancelled, stopping wait")
			return WaitCommandRes{}, ctx.Err()
		default:
			time.Sleep(1 * time.Second)
			println("Waiting for", duration, "seconds")
		}
	}
	return WaitCommandRes{}, nil
}
