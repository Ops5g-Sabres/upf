package pfcpiface

import (
	"context"
	"flag"
	"net"
	"time"

	pb "github.com/omec-project/upf/pfcpiface/click_pb/sdcore"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/wmnsk/go-pfcp/ie"
)

var clickIP = flag.String("click", "localhost:10515", "Click IP/port combo")

type click struct {
	client            pb.SDCoreControlClient
	conn              *grpc.ClientConn
	endMarkerSocket   net.Conn
	notifyClickSocket net.Conn
}

func (b *click) IsConnected(accessIP *net.IP) bool {
	if (b.conn == nil) || (b.conn.GetState() != connectivity.Ready) {
		return false
	}

	return true
}

func (b *click) SendEndMarkers(endMarkerList *[][]byte) error {
	return nil
}

func (b *click) AddSliceInfo(sliceInfo *SliceInfo) error {
	_, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	done := make(chan bool)

	rc := b.GRPCJoin(1, Timeout, done)

	if !rc {
		log.Errorln("Unable to make GRPC calls")
	}

	return nil
}

func (b *click) SendMsgToUPF(
	method upfMsgType, rules PacketForwardingRules, updated PacketForwardingRules) uint8 {
	// create context
	var cause uint8 = ie.CauseRequestAccepted

	pdrs := rules.pdrs
	fars := rules.fars
	qers := rules.qers

	if method == upfMsgTypeMod {
		pdrs = updated.pdrs
		fars = updated.fars
		qers = updated.qers
	}

	calls := len(pdrs) + len(fars) + len(qers)
	if calls == 0 {
		return cause
	}

	_, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	done := make(chan bool)

	for _, pdr := range pdrs {
		log.Traceln(method, pdr)
	}

	for _, far := range fars {
		log.Traceln(method, far)
	}

	for _, qer := range qers {
		log.Traceln(method, qer)
	}

	rc := b.GRPCJoin(calls, Timeout, done)
	if !rc {
		log.Info("Unable to make GRPC calls")
	}

	return cause
}

func (b *click) Exit() {
	log.Info("Exit function Click")
	b.conn.Close()
}

func (b *click) getPortStats(ifname string) *pb.GetStatsResponse {
	ctx := context.Background()
	req := &pb.GetStatsRequest{
		Interface: ifname,
	}

	res, err := b.client.GetStats(ctx, req)
	if err != nil {
		log.Info("Error calling GetPortStats", ifname, err)
		return nil
	}

	return res
}

func (b *click) PortStats(uc *upfCollector, ch chan<- prometheus.Metric) {
	return
}

func (b *click) SummaryLatencyJitter(uc *upfCollector, ch chan<- prometheus.Metric) {
	return
}

func (b *click) SessionStats(pc *PfcpNodeCollector, ch chan<- prometheus.Metric) error {
	return nil
}

func (b *click) processQER(ctx context.Context, thing1, thing2 string) error {

	resp, err := b.client.GetStats(ctx, &pb.GetStatsRequest{
		Interface: thing1,
		Rule:      thing2,
	})

	log.Traceln("qer lookup resp : ", resp)
	return err
}

func (b *click) processFAR(ctx context.Context, thing1, thing2 string) error {

	resp, err := b.client.GetStats(ctx, &pb.GetStatsRequest{
		Interface: thing1,
		Rule:      thing2,
	})

	log.Traceln("far lookup resp : ", resp)
	return err
}

func (b *click) processPDR(ctx context.Context, thing1, thing2 string) error {

	resp, err := b.client.GetStats(ctx, &pb.GetStatsRequest{
		Interface: thing1,
		Rule:      thing2,
	})

	log.Traceln("pdr lookup resp : ", resp)

	return err
}

func (b *click) clearState() {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	log.Debug("Clearing all the state in Click")

	b.processPDR(ctx, "1", "2")
	b.processFAR(ctx, "1", "2")
	b.processQER(ctx, "1", "2")
}

// setUpfInfo is only called at pfcp-agent's startup
// it clears all the state in Click
func (b *click) SetUpfInfo(u *upf, conf *Conf) {
	//var err error

	log.Info("SetUpfInfo click")

	// get click grpc client
	log.Infof("clickIP: %v", *clickIP)

	var err error
	b.conn, err = grpc.Dial(*clickIP, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("did not connect:", err)
	}

	b.client = pb.NewSDCoreControlClient(b.conn)

	b.clearState()

	/*
		if conf.EnableNotifyClick {
			notifySockAddr := conf.NotifySockAddr
			if notifySockAddr == "" {
				notifySockAddr = SockAddr
			}

			b.notifyClickSocket, err = net.Dial("unixpacket", notifySockAddr)
			if err != nil {
				log.Info("dial error:", err)
				return
			}
		}

		if conf.EnableEndMarker {
			pfcpCommAddr := conf.EndMarkerSockAddr
			if pfcpCommAddr == "" {
				pfcpCommAddr = PfcpAddr
			}

			b.endMarkerSocket, err = net.Dial("unixpacket", pfcpCommAddr)
			if err != nil {
				log.Info("dial error:", err)
				return
			}

			log.Info("Starting end marker loop")
		}

		if (conf.SliceMeterConfig.N6RateBps > 0) ||
			(conf.SliceMeterConfig.N3RateBps > 0) {
			_, cancel := context.WithTimeout(context.Background(), Timeout)
			defer cancel()

			done := make(chan bool)

			rc := b.GRPCJoin(1, Timeout, done)
			if !rc {
				log.Errorln("Unable to make GRPC calls")
			}
		}
	*/
}

func (b *click) setActionValue(f far) uint8 {
	if (f.applyAction & ActionForward) != 0 {
		if f.dstIntf == ie.DstInterfaceAccess {
			return farForwardD
		} else if (f.dstIntf == ie.DstInterfaceCore) || (f.dstIntf == ie.DstInterfaceSGiLANN6LAN) {
			return farForwardU
		}
	} else if (f.applyAction & ActionDrop) != 0 {
		return farDrop
	} else if (f.applyAction & ActionBuffer) != 0 {
		return farNotify
	} else if (f.applyAction & ActionNotify) != 0 {
		return farNotify
	}

	// default action
	return farDrop
}

func (b *click) GRPCJoin(calls int, timeout time.Duration, done chan bool) bool {
	/*
		boom := time.After(timeout)
			for {
				select {
				case ok := <-done:
					if !ok {
						log.Info("Error making GRPC calls")
						return false
					}

					calls--
					if calls == 0 {
						return true
					}
				case <-boom:
					log.Info("Timed out adding entries")
					return false
				}
			}
	*/
	return true
}
