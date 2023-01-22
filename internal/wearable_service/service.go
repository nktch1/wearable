package wearable_service

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/nktch1/wearable/pkg/clients/push_sender"
	"github.com/nktch1/wearable/pkg/server/wearable"
)

type Service struct {
	wearable.UnimplementedWearableServiceServer
	sender push_sender.PushSenderClient
}

func NewService() (*Service, error) {
	addr := "push-sender:51126"
	client, err := createClient(addr)
	if err != nil {
		return nil, fmt.Errorf("create push sender client: %w", err)
	}

	return &Service{sender: client}, nil
}

func (p *Service) RegisterGRPC(server *grpc.Server) {
	wearable.RegisterWearableServiceServer(server, p)
}

func createClient(address string) (push_sender.PushSenderClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("create push sender connect: %w", err)
	}

	return push_sender.NewPushSenderClient(conn), nil
}
