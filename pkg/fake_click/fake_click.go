// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Open Networking Foundation

package fake_click

import (
	"fmt"
	"net"

	"github.com/omec-project/upf/pfcpiface/click_pb"
	"google.golang.org/grpc"
)

type FakeClick struct {
	grpcServer *grpc.Server
	service    *fakeClickService
}

// NewFakeClick creates a new fake Click gRPC server. Its modules can be programmed in the same way
// as the real Click and keep track of their state.
func NewFakeClick() *FakeClick {
	return &FakeClick{
		service: newFakeClickService(),
	}
}

// Run starts and runs the Click gRPC server on the given address. Blocking until Stop is called.
func (b *FakeClick) Run(address string) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(address))
	if err != nil {
		return err
	}

	b.grpcServer = grpc.NewServer()
	click_pb.RegisterClickControlServer(b.grpcServer, b.service)

	// Blocking
	err = b.grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

// Stop the Click gRPC server.
func (b *FakeClick) Stop() {
	b.grpcServer.Stop()
}

func (b *FakeClick) GetPdrTableEntries() (entries map[uint32][]FakePdr) {
	entries = make(map[uint32][]FakePdr)
	msgs := b.service.GetOrAddModule(pdrLookupModuleName).GetState()
	for _, m := range msgs {
		e, ok := m.(*click_pb.WildcardMatchCommandAddArg)
		if !ok {
			panic("unexpected message type")
		}
		pdr := UnmarshalPdr(e)
		entries[pdr.PdrID] = append(entries[pdr.PdrID], pdr)
	}

	return
}

func (b *FakeClick) GetFarTableEntries() (entries map[uint32]FakeFar) {
	entries = make(map[uint32]FakeFar)
	msgs := b.service.GetOrAddModule(farLookupModuleName).GetState()
	for _, m := range msgs {
		e, ok := m.(*click_pb.ExactMatchCommandAddArg)
		if !ok {
			panic("unexpected message type")
		}
		far := UnmarshalFar(e)
		entries[far.FarID] = far
	}
	return
}

// Session QERs are missing a QerID and are therefore returned as a slice, not map.
func (b *FakeClick) GetSessionQerTableEntries() (entries []FakeQer) {
	msgs := b.service.GetOrAddModule(sessionQerModuleName).GetState()
	for _, m := range msgs {
		e, ok := m.(*click_pb.QosCommandAddArg)
		if !ok {
			panic("unexpected message type")
		}
		entries = append(entries, UnmarshalSessionQer(e))
	}
	return
}

func (b *FakeClick) GetAppQerTableEntries() (entries []FakeQer) {
	msgs := b.service.GetOrAddModule(appQerModuleName).GetState()
	for _, m := range msgs {
		e, ok := m.(*click_pb.QosCommandAddArg)
		if !ok {
			panic("unexpected message type")
		}
		entries = append(entries, UnmarshalAppQer(e))
	}
	return
}
