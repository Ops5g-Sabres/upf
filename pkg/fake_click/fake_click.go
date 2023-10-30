// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Open Networking Foundation

package fake_click

import (
	"fmt"
	"net"
	"sync"

	//mgmt "github.com/omec-project/upf/pfcpiface/click_pb/mgmt"
	//sabres "github.com/omec-project/upf/pfcpiface/click_pb/sabres"
	click_pb "github.com/omec-project/upf/pfcpiface/click_pb/sdcore"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type fakeClickService struct {
	click_pb.UnimplementedSDCoreControlServer
	mtx sync.Mutex
}

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

	log.Infof("Fake click run called on: %s", address)

	listener, err := net.Listen("tcp", fmt.Sprintf(address))
	if err != nil {
		return err
	}

	b.grpcServer = grpc.NewServer()
	click_pb.RegisterSDCoreControlServer(b.grpcServer, b.service)

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

func (b *FakeClick) GetPdrTableEntries() {
	return
}

func (b *FakeClick) GetFarTableEntries() {
	return
}

// Session QERs are missing a QerID and are therefore returned as a slice, not map.
func (b *FakeClick) GetSessionQerTableEntries() {
	return
}

func (b *FakeClick) GetAppQerTableEntries() {
	return
}
