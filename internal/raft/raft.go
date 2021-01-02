package raft

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/CyrusJavan/dovii/internal/dovii"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

func NewRaftServer(store *dovii.KVStore) (*raft.Raft, error) {
	adAddr := getOutboundTCPAddr()
	adAddr.Port = 7654
	conf := raft.DefaultConfig()
	conf.LocalID = raft.ServerID(adAddr.String())
	bs, err := raftboltdb.NewBoltStore(fmt.Sprintf("/tmp/raftboltdbstore-%s", adAddr.String()))
	if err != nil {
		return nil, fmt.Errorf("init boltstore: %w", err)
	}

	ss := raft.DiscardSnapshotStore{}

	log.Println("advertise address:", adAddr.String())
	t, err := raft.NewTCPTransport("0.0.0.0:"+strconv.Itoa(adAddr.Port), &adAddr, 3, time.Second, log.Writer())
	if err != nil {
		return nil, fmt.Errorf("new tcp transport: %w", err)
	}

	fsm := BitcaskFSM{
		Store: store,
	}
	raftServer, err := raft.NewRaft(conf, &fsm, bs, bs, &ss, t)
	if err != nil {
		return nil, fmt.Errorf("starting raft rpc server: %w", err)
	}
	hn, err := os.Hostname()
	if err != nil {
		log.Println("getting hostname:", err)
	}
	log.Println(hn)

	addrs, err := net.LookupHost("dovii")
	if err != nil {
		return nil, fmt.Errorf("dns lookup: %w", err)
	}

	log.Println("addrs:", addrs)

	var servers []raft.Server
	for _, addr := range addrs {
		servers = append(servers, raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(addr + ":" + strconv.Itoa(adAddr.Port)),
			Address:  raft.ServerAddress(addr + ":" + strconv.Itoa(adAddr.Port)),
		})
	}

	nodeConfig := raft.Configuration{Servers: servers}
	future := raftServer.BootstrapCluster(nodeConfig)
	err = future.Error()
	if err != nil {
		log.Println("bootstrap cluster:", err)
	}

	return raftServer, nil
}

func getOutboundTCPAddr() net.TCPAddr {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return net.TCPAddr{
		IP:   localAddr.IP,
		Port: localAddr.Port,
		Zone: localAddr.Zone,
	}
}
