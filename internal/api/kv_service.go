package api

import (
	"context"

	"github.com/Saad7890-web/meridian/internal/storage"
	pb "github.com/Saad7890-web/meridian/proto"
)

type KVService struct {
	pb.UnimplementedKVServer
	store *storage.Store
}

func NewKVService(store *storage.Store) *KVService {
	return &KVService{store: store}
}

func (s *KVService) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	err := s.store.Put(req.Key, req.Value)
	if err != nil {
		return nil, err
	}
	return &pb.PutResponse{Success: true}, nil
}

func (s *KVService) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	entry, ok := s.store.Get(req.Key)
	if !ok {
		return &pb.GetResponse{Found: false}, nil
	}

	return &pb.GetResponse{
		Key:     entry.Key,
		Value:   entry.Value,
		Version: entry.Version,
		Found:   true,
	}, nil
}

func (s *KVService) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := s.store.Delete(req.Key)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteResponse{Success: true}, nil
}