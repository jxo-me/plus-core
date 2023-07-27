package ws

import (
	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	ErrConnectionLoss = gerror.New("ERR_CONNECTION_LOSS")

	ErrSendMessageFull = gerror.New("ERR_SEND_MESSAGE_FULL")

	ErrJoinRoomTwice = gerror.New("ERR_JOIN_ROOM_TWICE")

	ErrNotInRoom = gerror.New("ERR_NOT_IN_ROOM")

	ErrRoomIdInvalid = gerror.New("ERR_ROOM_ID_INVALID")

	ErrDispatchChannelFull = gerror.New("ERR_DISPATCH_CHANNEL_FULL")

	ErrMergeChannelFull = gerror.New("ERR_MERGE_CHANNEL_FULL")
)
