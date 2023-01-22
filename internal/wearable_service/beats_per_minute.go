package wearable_service

import (
	"context"
	"fmt"
	"github.com/nktch1/wearable/pkg/clients/push_sender"
	"math/rand"
	"time"

	"github.com/nktch1/wearable/pkg/server/wearable"
)

func (p *Service) BeatsPerMinute(in *wearable.BeatsPerMinuteRequest, stream wearable.WearableService_BeatsPerMinuteServer) error {
	const batchSize = 100

	for idx := 0; idx < batchSize; idx++ {
		heartRate := newRandInt()

		if somethingIsGoingWrong(heartRate) {
			_, err := p.sender.Notify(context.Background(), &push_sender.NotifyRequest{
				Uuid:    "some_uuid",
				Message: "Something is going wrong!",
			})

			if err != nil {
				return fmt.Errorf("notify: %w", err)
			}
		}

		response := wearable.BeatsPerMinuteResponse{
			Value:  newRandInt(),
			Minute: uint32(idx),
		}

		stream.Send(&response)

		time.Sleep(time.Millisecond * 300)
	}

	return nil
}

func newRandInt() uint32 {
	min := 30
	max := 160
	return uint32(rand.Intn(max-min) + min)
}

func somethingIsGoingWrong(heartRate uint32) bool {
	return heartRate < 40 || heartRate > 140
}
