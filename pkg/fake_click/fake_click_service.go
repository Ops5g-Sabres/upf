// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Open Networking Foundation

package fake_click

import (
	"context"

	emptypb "github.com/golang/protobuf/ptypes/empty"
	click_pb "github.com/omec-project/upf/pfcpiface/click_pb/sdcore"
)

func newFakeClickService() *fakeClickService {
	return &fakeClickService{}
}

func (b *fakeClickService) CreateRule(ctx context.Context, request *click_pb.CreateRuleRequest) (*emptypb.Empty, error) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return &emptypb.Empty{}, nil
}

func (b *fakeClickService) ModifyRule(ctx context.Context, request *click_pb.ModifyRuleRequest) (*emptypb.Empty, error) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return &emptypb.Empty{}, nil
}

func (b *fakeClickService) DeleteRule(ctx context.Context, request *click_pb.DeleteRuleRequest) (*emptypb.Empty, error) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return &emptypb.Empty{}, nil
}

func (b *fakeClickService) GetStats(ctx context.Context, request *click_pb.GetStatsRequest) (*click_pb.GetStatsResponse, error) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return &click_pb.GetStatsResponse{}, nil
}

func (b *fakeClickService) ResetStats(ctx context.Context, request *click_pb.ResetStatsRequest) (*emptypb.Empty, error) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return &emptypb.Empty{}, nil
}

func (b *fakeClickService) SendCommand(ctx context.Context, request *click_pb.SendCommandRequest) (*click_pb.SendCommandResponse, error) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return &click_pb.SendCommandResponse{}, nil
}
