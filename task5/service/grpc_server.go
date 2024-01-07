package service

import (
	"context"
	"errors"
	"lesson15_16/task5/api"
	"lesson15_16/task5/storage"
	"net"
	"unicode/utf8"

	"google.golang.org/grpc"
)

var (
	ErrLessOrEqualZero = errors.New("less than or equal to 0")
	ErrLenMessage      = errors.New("message less than three")
)

type GRPCServer struct {
	storage storage.Storage
	api.UnimplementedStorageServer
}

func NewGRPCStorage(storage storage.Storage) (*GRPCServer, error) {
	return &GRPCServer{
		storage: storage,
	}, nil
}

func (s *GRPCServer) Save(ctx context.Context, req *api.TaskRequest) (*api.TaskResponse, error) {
	if utf8.RuneCountInString(req.GetMessage()) < 3 {
		return nil, ErrLenMessage
	}

	id, err := s.storage.Save(ctx, req.GetMessage())
	if err != nil {
		return nil, err
	}

	return &api.TaskResponse{
		Id: id,
	}, nil
}

func (s *GRPCServer) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	if req.GetId() <= 0 {
		return nil, ErrLessOrEqualZero
	}

	task, err := s.storage.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &api.GetResponse{
		Id:      task.ID,
		Message: task.Message,
	}, nil
}

func (s *GRPCServer) Start() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	api.RegisterStorageServer(server, s)

	err = server.Serve(l)
	if err != nil {
		panic(err)
	}
}
