//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package grpc provides grpc server logic
package grpc

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/meta"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/pkg/meta/cassandra/service"
)

type server struct {
	cassandra service.Cassandra
}

func New(opts ...Option) meta.MetaServer {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) GetMeta(ctx context.Context, key *payload.Meta_Key) (*payload.Meta_Val, error) {
	val, err := s.cassandra.Get(key.GetKey())
	if err != nil {
		return nil, err
	}
	return &payload.Meta_Val{
		Val: val,
	}, nil
}

func (s *server) GetMetas(ctx context.Context, keys *payload.Meta_Keys) (*payload.Meta_Vals, error) {
	vals, err := s.cassandra.GetMultiple(keys.GetKeys()...)
	if err != nil {
		return nil, err
	}
	return &payload.Meta_Vals{
		Vals: vals,
	}, nil
}

func (s *server) GetMetaInverse(ctx context.Context, val *payload.Meta_Val) (*payload.Meta_Key, error) {
	key, err := s.cassandra.GetInverse(val.GetVal())
	if err != nil {
		return nil, err
	}
	return &payload.Meta_Key{
		Key: key,
	}, nil
}

func (s *server) GetMetasInverse(ctx context.Context, vals *payload.Meta_Vals) (*payload.Meta_Keys, error) {
	keys, err := s.cassandra.GetInverseMultiple(vals.GetVals()...)
	if err != nil {
		return nil, err
	}
	return &payload.Meta_Keys{
		Keys: keys,
	}, nil
}

func (s *server) SetMeta(ctx context.Context, kv *payload.Meta_KeyVal) (_ *payload.Empty, err error) {
	err = s.cassandra.Set(kv.GetKey(), kv.GetVal())
	if err != nil {
		return nil, err
	}
	return new(payload.Empty), nil
}

func (s *server) SetMetas(ctx context.Context, kvs *payload.Meta_KeyVals) (_ *payload.Empty, err error) {
	query := make(map[string]string, len(kvs.GetKvs())/2)
	for _, kv := range kvs.GetKvs() {
		query[kv.GetKey()] = kv.GetVal()
	}
	err = s.cassandra.SetMultiple(query)
	if err != nil {
		return nil, err
	}
	return new(payload.Empty), nil

}

func (s *server) DeleteMeta(ctx context.Context, key *payload.Meta_Key) (*payload.Meta_Val, error) {
	val, err := s.cassandra.Delete(key.GetKey())
	if err != nil {
		return nil, err
	}
	return &payload.Meta_Val{
		Val: val,
	}, nil
}

func (s *server) DeleteMetas(ctx context.Context, keys *payload.Meta_Keys) (*payload.Meta_Vals, error) {
	vals, err := s.cassandra.DeleteMultiple(keys.GetKeys()...)
	if err != nil {
		return nil, err
	}
	return &payload.Meta_Vals{
		Vals: vals,
	}, nil
}

func (s *server) DeleteMetaInverse(ctx context.Context, val *payload.Meta_Val) (*payload.Meta_Key, error) {
	key, err := s.cassandra.DeleteInverse(val.GetVal())
	if err != nil {
		return nil, err
	}
	return &payload.Meta_Key{
		Key: key,
	}, nil
}

func (s *server) DeleteMetasInverse(ctx context.Context, vals *payload.Meta_Vals) (*payload.Meta_Keys, error) {
	keys, err := s.cassandra.DeleteInverseMultiple(vals.GetVals()...)
	if err != nil {
		return nil, err
	}
	return &payload.Meta_Keys{
		Keys: keys,
	}, nil
}