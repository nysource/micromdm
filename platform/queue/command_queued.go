package queue

import (
	"errors"

	"google.golang.org/protobuf/proto"

	commandqueued "github.com/micromdm/micromdm/platform/queue/internal/commandqueuedproto"
)

type QueueCommandQueued struct {
	DeviceUDID  string
	CommandUUID string
}

func MarshalQueuedCommand(cq *QueueCommandQueued) ([]byte, error) {
	if cq == nil {
		return nil, errors.New("marshalling nil QueueCommandQueued")
	}
	return proto.Marshal(&commandqueued.CommandQueued{
		DeviceUdid:  cq.DeviceUDID,
		CommandUuid: cq.CommandUUID,
	})
}

func UnmarshalQueuedCommand(data []byte) (*QueueCommandQueued, error) {
	cmdQueued := commandqueued.CommandQueued{}
	if err := proto.Unmarshal(data, &cmdQueued); err != nil {
		return nil, err
	}
	queueCmdQueued := new(QueueCommandQueued)
	queueCmdQueued.DeviceUDID = cmdQueued.DeviceUdid
	queueCmdQueued.CommandUUID = cmdQueued.CommandUuid
	return queueCmdQueued, nil
}
