package ws

import "context"

func (i *Instance) PushAll(msg *MessageReq) error {
	return i.ConnManager.PushAll(msg)
}

func (i *Instance) PushSession(connId uint64, message *MessageReq) error {
	return i.ConnManager.PushByCID(connId, message)
}

func (i *Instance) CloseSession(ctx context.Context, connId uint64) error {
	i.ConnManager.DelConnByCID(ctx, connId)
	return nil
}

func (i *Instance) PushRoom(roomId string, msg *MessageReq) error {
	return i.ConnManager.PushRoom(roomId, msg)
}

func (i *Instance) GetStats() *Statistics {
	return i.Stat.GetStats()
}
