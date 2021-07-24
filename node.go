/*
Package goroslib is a library in pure Go that allows to build clients (nodes)
for the Robot Operating System (ROS).

Basic example (more are available at https://github.com/aler9/goroslib/tree/master/examples):

  package main

  import (
      "fmt"
      "github.com/aler9/goroslib"
      "github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
  )

  func onMessage(msg *sensor_msgs.Imu) {
      fmt.Printf("Incoming: %+v\n", msg)
  }

  func main() {
      n, err := goroslib.NewNode(goroslib.NodeConf{
          Name:          "goroslib",
          MasterAddress: "127.0.0.1:11311",
      })
      if err != nil {
          panic(err)
      }
      defer n.Close()

      sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
          Node:     n,
          Topic:    "test_topic",
          Callback: onMessage,
      })
      if err != nil {
          panic(err)
      }
      defer sub.Close()

      select {}
  }

*/
package goroslib

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aler9/goroslib/pkg/apimaster"
	"github.com/aler9/goroslib/pkg/apiparam"
	"github.com/aler9/goroslib/pkg/apislave"
	"github.com/aler9/goroslib/pkg/msgs/rosgraph_msgs"
	"github.com/aler9/goroslib/pkg/prototcp"
	"github.com/aler9/goroslib/pkg/protoudp"
	"github.com/aler9/goroslib/pkg/xmlrpc"
)

type getPublicationsReq struct {
	res chan [][]string
}

type getBusInfoReq struct {
	res chan apislave.ResponseGetBusInfo
}

type getBusInfoSubReq struct {
	pbusInfo *[][]interface{}
	done     chan struct{}
}

type tcpConnSubscriberReq struct {
	conn   *prototcp.Conn
	header *prototcp.HeaderSubscriber
}

type tcpConnServiceClientReq struct {
	conn   *prototcp.Conn
	header *prototcp.HeaderServiceClient
}

type udpSubPublisherCloseReq struct {
	sp   *subscriberPublisher
	done chan struct{}
}

type udpFrameReq struct {
	frame  *protoudp.Frame
	source *net.UDPAddr
}

type subscriberPubUpdateReq struct {
	topic string
	urls  []string
}

type subscriberRequestTopicReq struct {
	req *apislave.RequestRequestTopic
	res chan apislave.ResponseRequestTopic
}

type subscriberNewReq struct {
	sub *Subscriber
	err chan error
}

type publisherNewReq struct {
	pub *Publisher
	err chan error
}

type serviceProviderNewReq struct {
	sp  *ServiceProvider
	err chan error
}

type simtimeSleep struct {
	value time.Time
	done  chan struct{}
}

// NodeConf is the configuration of a Node.
type NodeConf struct {
	// (optional) hostname (or ip) and port of the master node.
	// It defaults to 127.0.0.1:11311
	MasterAddress string

	// (optional) namespace of this node.
	// It defaults to '/' (global namespace).
	Namespace string

	// name of this node.
	Name string

	// (optional) hostname or ip of this node, needed by other nodes
	// in order to communicate with it.
	// if not provided, it will be set automatically.
	Host string

	// (optional) port of the Slave API server of this node.
	// if not provided, it will be chosen automatically.
	ApislavePort int

	// (optional) port of the TCPROS server of this node.
	// if not provided, it will be chosen automatically.
	TcprosPort int

	// (optional) port of the UDPROS server of this node.
	// if not provided, it will be chosen automatically.
	UdprosPort int
}

// Node is a ROS Node, an entity that can create subscribers, publishers, service providers
// and service clients.
type Node struct {
	conf NodeConf

	ctx                 context.Context
	ctxCancel           func()
	masterAddr          *net.TCPAddr
	nodeAddr            *net.TCPAddr
	apiMasterClient     *apimaster.Client
	apiParamClient      *apiparam.Client
	apiSlaveServer      *apislave.Server
	apiSlaveServerURL   string
	tcprosServer        *prototcp.Server
	tcprosServerURL     string
	udprosServer        *protoudp.Server
	tcprosConns         map[*prototcp.Conn]struct{}
	udprosSubPublishers map[*subscriberPublisher]struct{}
	subscribers         map[string]*Subscriber
	publishers          map[string]*Publisher
	serviceProviders    map[string]*ServiceProvider
	publisherLastID     int
	rosoutPublisher     *Publisher
	simtimeEnabled      bool
	simtimeSubscriber   *Subscriber
	simtimeMutex        sync.RWMutex
	simtimeInitialized  bool
	simtimeValue        time.Time
	simtimeSleeps       []*simtimeSleep

	// in
	getPublications        chan getPublicationsReq
	getBusInfo             chan getBusInfoReq
	tcpConnNew             chan *prototcp.Conn
	tcpConnClose           chan *prototcp.Conn
	tcpConnSubscriber      chan tcpConnSubscriberReq
	tcpConnServiceClient   chan tcpConnServiceClientReq
	udpSubPublisherNew     chan *subscriberPublisher
	udpSubPublisherClose   chan udpSubPublisherCloseReq
	udpFrame               chan udpFrameReq
	subscriberRequestTopic chan subscriberRequestTopicReq
	subscriberNew          chan subscriberNewReq
	subscriberClose        chan *Subscriber
	subscriberPubUpdate    chan subscriberPubUpdateReq
	publisherNew           chan publisherNewReq
	publisherClose         chan *Publisher
	serviceProviderNew     chan serviceProviderNewReq
	serviceProviderClose   chan *ServiceProvider

	// out
	done chan struct{}
}

// NewNode allocates a Node. See NodeConf for the options.
func NewNode(conf NodeConf) (*Node, error) {
	if os.Getenv("ROS_NAMESPACE") != "" {
		conf.Namespace = os.Getenv("ROS_NAMESPACE")
	}

	if conf.Namespace == "" {
		conf.Namespace = "/"
	}
	if conf.Namespace[0] != '/' {
		return nil, fmt.Errorf("Namespace must begin with a slash (/)")
	}
	if conf.Namespace != "/" && conf.Namespace[len(conf.Namespace)-1] == '/' {
		return nil, fmt.Errorf("Namespace can't end with a slash (/)")
	}

	if conf.Name == "" {
		return nil, fmt.Errorf("Name not provided")
	}
	if strings.ContainsRune(conf.Name, '/') {
		return nil, fmt.Errorf("Name cannot contain slashes (/), use Namespace to set a namespace")
	}

	if len(conf.MasterAddress) == 0 {
		conf.MasterAddress = "127.0.0.1:11311"
	}

	// support ROS-style master address, in order to increase interoperability
	conf.MasterAddress = strings.TrimPrefix(conf.MasterAddress, "http://")

	// solve master address once
	masterAddr, err := net.ResolveTCPAddr("tcp", conf.MasterAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to solve master address: %s", err)
	}
	if masterAddr.Zone != "" {
		return nil, fmt.Errorf("the master address has a stateless IPv6, which is not supported")
	}

	// find an ip in the same subnet of the master
	if conf.Host == "" {
		conf.Host = func() string {
			ifaces, err := net.Interfaces()
			if err != nil {
				return ""
			}

			for _, i := range ifaces {
				addrs, err := i.Addrs()
				if err != nil {
					continue
				}

				for _, addr := range addrs {
					if v, ok := addr.(*net.IPNet); ok {
						if v.Contains(masterAddr.IP) {
							return v.IP.String()
						}
					}
				}
			}
			return ""
		}()
		if conf.Host == "" {
			return nil, fmt.Errorf("unable to set Host automatically")
		}
	}

	// solve node host once
	nodeAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(conf.Host, "0"))
	if err != nil {
		return nil, fmt.Errorf("unable to solve node host: %s", err)
	}
	if nodeAddr.Zone != "" {
		return nil, fmt.Errorf("the node IP is a stateless IPv6, which is not supported")
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	n := &Node{
		conf:                   conf,
		ctx:                    ctx,
		ctxCancel:              ctxCancel,
		masterAddr:             masterAddr,
		nodeAddr:               nodeAddr,
		tcprosConns:            make(map[*prototcp.Conn]struct{}),
		udprosSubPublishers:    make(map[*subscriberPublisher]struct{}),
		subscribers:            make(map[string]*Subscriber),
		publishers:             make(map[string]*Publisher),
		serviceProviders:       make(map[string]*ServiceProvider),
		simtimeValue:           time.Unix(0, 0),
		getPublications:        make(chan getPublicationsReq),
		getBusInfo:             make(chan getBusInfoReq),
		tcpConnNew:             make(chan *prototcp.Conn),
		tcpConnClose:           make(chan *prototcp.Conn),
		tcpConnSubscriber:      make(chan tcpConnSubscriberReq),
		tcpConnServiceClient:   make(chan tcpConnServiceClientReq),
		udpSubPublisherNew:     make(chan *subscriberPublisher),
		udpSubPublisherClose:   make(chan udpSubPublisherCloseReq),
		udpFrame:               make(chan udpFrameReq),
		subscriberRequestTopic: make(chan subscriberRequestTopicReq),
		subscriberNew:          make(chan subscriberNewReq),
		subscriberClose:        make(chan *Subscriber),
		subscriberPubUpdate:    make(chan subscriberPubUpdateReq),
		publisherNew:           make(chan publisherNewReq),
		publisherClose:         make(chan *Publisher),
		serviceProviderNew:     make(chan serviceProviderNewReq),
		serviceProviderClose:   make(chan *ServiceProvider),
		done:                   make(chan struct{}),
	}

	n.apiMasterClient = apimaster.NewClient(masterAddr.String(), n.absoluteName())

	n.apiParamClient = apiparam.NewClient(masterAddr.String(), n.absoluteName())

	n.apiSlaveServer, err = apislave.NewServer(":" + strconv.FormatInt(int64(conf.ApislavePort), 10))
	if err != nil {
		return nil, err
	}
	n.apiSlaveServerURL = xmlrpc.ServerURL(nodeAddr, n.apiSlaveServer.Port())

	n.tcprosServer, err = prototcp.NewServer(":" + strconv.FormatInt(int64(conf.TcprosPort), 10))
	if err != nil {
		n.apiSlaveServer.Close()
		return nil, err
	}
	n.tcprosServerURL = prototcp.ServerURL(nodeAddr, n.tcprosServer.Port())

	n.udprosServer, err = protoudp.NewServer(":" + strconv.FormatInt(int64(conf.UdprosPort), 10))
	if err != nil {
		n.tcprosServer.Close()
		n.apiSlaveServer.Close()
		return nil, err
	}

	go n.run()

	n.rosoutPublisher, err = NewPublisher(PublisherConf{
		Node:  n,
		Topic: "/rosout",
		Msg:   &rosgraph_msgs.Log{},
	})
	if err != nil {
		n.Close()
		return nil, err
	}

	isSet, err := n.ParamIsSet("/use_sim_time")
	if err != nil {
		n.Close()
		return nil, err
	}

	if isSet {
		n.simtimeEnabled, err = n.ParamGetBool("/use_sim_time")
		if err != nil {
			n.Close()
			return nil, err
		}
	}

	if n.simtimeEnabled {
		n.simtimeSubscriber, err = NewSubscriber(SubscriberConf{
			Node:  n,
			Topic: "/clock",
			Callback: func(msg *rosgraph_msgs.Clock) {
				n.simtimeMutex.Lock()
				defer n.simtimeMutex.Unlock()

				// reinitialize sleeps if simulation time was not initialized before
				if !n.simtimeInitialized {
					n.simtimeInitialized = true
					zero := time.Unix(0, 0)
					for _, s := range n.simtimeSleeps {
						s.value = msg.Clock.Add(s.value.Sub(zero))
					}
				}

				n.simtimeValue = msg.Clock

				for i := 0; i < len(n.simtimeSleeps); {
					s := n.simtimeSleeps[i]
					if !s.value.After(n.simtimeValue) {
						close(s.done)
						n.simtimeSleeps = append(n.simtimeSleeps[:i], n.simtimeSleeps[i+1:]...)
					} else {
						i++
					}
				}
			},
		})
		if err != nil {
			n.Close()
			return nil, err
		}
	}

	return n, nil
}

// Close closes a Node and all its resources.
func (n *Node) Close() error {
	n.ctxCancel()
	<-n.done
	return nil
}

func (n *Node) absoluteTopicName(topic string) string {
	// topic is absolute
	if topic[0] == '/' {
		return topic
	}

	// topic is relative
	if n.conf.Namespace == "/" {
		return "/" + topic
	}
	return n.conf.Namespace + "/" + topic
}

func (n *Node) absoluteName() string {
	if n.conf.Namespace == "/" {
		return "/" + n.conf.Name
	}
	return n.conf.Namespace + "/" + n.conf.Name
}

func (n *Node) run() {
	defer close(n.done)

	var serversWg sync.WaitGroup

	serversWg.Add(3)
	go n.runAPISlaveServer(&serversWg)
	go n.runTcprosServer(&serversWg)
	go n.runUdprosServer(&serversWg)

	var clientsWg sync.WaitGroup

outer:
	for {
		select {
		case req := <-n.getPublications:
			res := [][]string{}
			for _, pub := range n.publishers {
				res = append(res, []string{pub.conf.Topic, pub.msgType})
			}
			req.res <- res

		case req := <-n.getBusInfo:
			var busInfo [][]interface{}

			for _, pub := range n.publishers {
				done := make(chan struct{})
				select {
				case pub.getBusInfo <- getBusInfoSubReq{&busInfo, done}:
					<-done
				case <-pub.ctx.Done():
				}
			}

			for _, sub := range n.subscribers {
				done := make(chan struct{})
				select {
				case sub.getBusInfo <- getBusInfoSubReq{&busInfo, done}:
					<-done
				case <-sub.ctx.Done():
				}
			}

			req.res <- apislave.ResponseGetBusInfo{
				Code:          1,
				StatusMessage: "",
				BusInfo:       busInfo,
			}

		case <-n.ctx.Done():
			break outer

		case conn := <-n.tcpConnNew:
			n.tcprosConns[conn] = struct{}{}
			clientsWg.Add(1)
			go n.runTcprosServerConn(&clientsWg, conn)

		case conn := <-n.tcpConnClose:
			delete(n.tcprosConns, conn)

		case req := <-n.tcpConnSubscriber:
			// pass conn ownership to publisher, if exists
			delete(n.tcprosConns, req.conn)

			pub, ok := n.publishers[req.header.Topic]
			if !ok {
				req.conn.Close()
				continue
			}

			select {
			case pub.subscriberTCPNew <- tcpConnSubscriberReq{
				conn:   req.conn,
				header: req.header,
			}:
			case <-pub.ctx.Done():
			}

		case req := <-n.tcpConnServiceClient:
			// pass conn ownership to service provider, if exists
			delete(n.tcprosConns, req.conn)

			sp, ok := n.serviceProviders[req.header.Service]
			if !ok {
				req.conn.Close()
				continue
			}

			select {
			case sp.clientNew <- req:
			case <-sp.ctx.Done():
			}

		case sp := <-n.udpSubPublisherNew:
			n.udprosSubPublishers[sp] = struct{}{}

		case req := <-n.udpSubPublisherClose:
			delete(n.udprosSubPublishers, req.sp)
			close(req.done)

		case req := <-n.udpFrame:
			for sp := range n.udprosSubPublishers {
				if req.frame.ConnectionID == sp.udpID &&
					req.source.IP.Equal(sp.udpAddr.IP) {

					select {
					case sp.udpFrame <- req.frame:
					case <-sp.ctx.Done():
					}
					break
				}
			}

		case req := <-n.subscriberRequestTopic:
			pub, ok := n.publishers[req.req.Topic]
			if !ok {
				req.res <- apislave.ResponseRequestTopic{
					Code:          0,
					StatusMessage: "topic not found",
				}
				continue
			}

			select {
			case pub.requestTopic <- subscriberRequestTopicReq{req.req, req.res}:
			case <-pub.ctx.Done():
				req.res <- apislave.ResponseRequestTopic{
					Code:          0,
					StatusMessage: "terminating",
				}
			}

		case req := <-n.subscriberNew:
			_, ok := n.subscribers[n.absoluteTopicName(req.sub.conf.Topic)]
			if ok {
				req.err <- fmt.Errorf("Topic %s already subscribed", req.sub.conf.Topic)
				continue
			}

			uris, err := n.apiMasterClient.RegisterSubscriber(
				n.absoluteTopicName(req.sub.conf.Topic),
				req.sub.msgType,
				n.apiSlaveServerURL)
			if err != nil {
				req.err <- err
				continue
			}

			n.subscribers[n.absoluteTopicName(req.sub.conf.Topic)] = req.sub
			req.err <- nil

			// send initial publishers list to subscriber
			select {
			case req.sub.subscriberPubUpdate <- uris:
			case <-req.sub.ctx.Done():
			}

		case sub := <-n.subscriberClose:
			delete(n.subscribers, n.absoluteTopicName(sub.conf.Topic))

		case req := <-n.subscriberPubUpdate:
			sub, ok := n.subscribers[req.topic]
			if !ok {
				continue
			}

			select {
			case sub.subscriberPubUpdate <- req.urls:
			case <-sub.ctx.Done():
			}

		case req := <-n.publisherNew:
			_, ok := n.publishers[n.absoluteTopicName(req.pub.conf.Topic)]
			if ok {
				req.err <- fmt.Errorf("Topic %s already published", req.pub.conf.Topic)
				continue
			}

			_, err := n.apiMasterClient.RegisterPublisher(
				n.absoluteTopicName(req.pub.conf.Topic),
				req.pub.msgType,
				n.apiSlaveServerURL)
			if err != nil {
				req.err <- err
				continue
			}

			n.publisherLastID++
			req.pub.id = n.publisherLastID
			n.publishers[n.absoluteTopicName(req.pub.conf.Topic)] = req.pub
			req.err <- nil

		case pub := <-n.publisherClose:
			delete(n.publishers, n.absoluteTopicName(pub.conf.Topic))

		case req := <-n.serviceProviderNew:
			_, ok := n.serviceProviders[n.absoluteTopicName(req.sp.conf.Name)]
			if ok {
				req.err <- fmt.Errorf("Service %s already provided", req.sp.conf.Name)
				continue
			}

			err := n.apiMasterClient.RegisterService(
				n.absoluteTopicName(req.sp.conf.Name),
				n.tcprosServerURL,
				n.apiSlaveServerURL)
			if err != nil {
				req.err <- err
				continue
			}

			n.serviceProviders[n.absoluteTopicName(req.sp.conf.Name)] = req.sp
			req.err <- nil

		case sp := <-n.serviceProviderClose:
			delete(n.serviceProviders, n.absoluteTopicName(sp.conf.Name))
		}
	}

	n.ctxCancel()

	n.apiSlaveServer.Close()
	n.tcprosServer.Close()
	n.udprosServer.Close()
	serversWg.Wait()

	for c := range n.tcprosConns {
		c.Close()
	}
	clientsWg.Wait()

	if n.simtimeSubscriber != nil {
		n.simtimeSubscriber.Close()
	}

	if n.rosoutPublisher != nil {
		n.rosoutPublisher.Close()
	}
}
