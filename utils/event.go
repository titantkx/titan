package utils

import (
	"errors"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/cosmos/gogoproto/proto"
)

var ErrTypedEventNotFound = errors.New("typed event not found")

func GetABCIEventAttribute(event *abci.Event, key string) string {
	for _, attr := range event.Attributes {
		if attr.Key == key {
			return attr.Value
		}
	}
	return ""
}

func GetABCIEvent(ctx sdk.Context, eventType string) *abci.Event {
	events := ctx.EventManager().ABCIEvents()
	for i := range events {
		if events[i].Type == eventType {
			return &events[i]
		}
	}
	return nil
}

func GetTypedEvent[T proto.Message](ctx sdk.Context, msg T) (T, error) {
	eventType := proto.MessageName(msg)

	event := GetABCIEvent(ctx, eventType)
	if event == nil {
		return msg, ErrTypedEventNotFound
	}

	protoMsg, err := sdk.ParseTypedEvent(*event)
	if err != nil {
		return msg, err
	}

	actualMsg, ok := protoMsg.(T)
	if !ok {
		return msg, fmt.Errorf("unexpected event type: got %s expected %s", proto.MessageName(protoMsg), eventType)
	}

	return actualMsg, nil
}
