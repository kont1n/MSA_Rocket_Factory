package decoder

import (
	"fmt"
	"time"

	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
	eventsV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewAssemblyRecordedDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.AssemblyRecordedEvent, error) {
	var pb eventsV1.AssemblyRecorded
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.AssemblyRecordedEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	var observedAt *time.Time
	if pb.ObservedAt != nil {
		observedAt = lo.ToPtr(pb.ObservedAt.AsTime())
	}

	return model.AssemblyRecordedEvent{
		UUID:        pb.Uuid,
		ObservedAt:  observedAt,
		Location:    pb.Location,
		Description: pb.Description,
	}, nil
}
