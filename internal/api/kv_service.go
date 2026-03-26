package api

import (
	"context"
	"encoding/json"

	"github.com/Saad7890-web/meridian/internal/raft"
	"github.com/Saad7890-web/meridian/internal/storage"
	pb "github.com/Saad7890-web/meridian/proto"
)

type KVService struct {
	pb.UnimplementedKVServer
	store *storage.Store
	raftNode *raft.Node
}

func NewKVService(store *storage.Store) *KVService {
	return &KVService{store: store}
}

func (s *KVService) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	cmd := map[string]string{
		"op":    "put",
		"key":   req.Key,
		"value": req.Value,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	
	err = s.raftNode.Apply(data)
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
	cmd := map[string]string{
		"op":  "delete",
		"key": req.Key,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	err = s.raftNode.Apply(data)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{Success: true}, nil
}