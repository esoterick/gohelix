package gohelix

import (
	"fmt"
	"sync"
)

type HelixManager struct {
	zkAddress string

	conn *Connection
}

type (
	ExternalViewChangeListener   func(externalViews []*Record, context *Context)
	LiveInstanceChangeListener   func(liveInstances []*Record, context *Context)
	CurrentStateChangeListener   func(instance string, currentState []*Record, context *Context)
	IdealStateChangeListener     func(idealState []*Record, context *Context)
	InstanceConfigChangeListener func(configs []*Record, context *Context)
	ControllerMessageListener    func(messages []*Record, context *Context)
	MessageListener              func(instance string, messages []*Record, context *Context)
)

func NewHelixManager(zkAddress string) *HelixManager {
	return &HelixManager{
		zkAddress: zkAddress,
	}
}

func (m *HelixManager) NewSpectator(clusterID string) *Spectator {
	return &Spectator{
		ClusterID:                   clusterID,
		zkConnStr:                   m.zkAddress,
		externalViewListeners:       []ExternalViewChangeListener{},
		liveInstanceChangeListeners: []LiveInstanceChangeListener{},
		currentStateChangeListeners: map[string][]CurrentStateChangeListener{},
		idealStateChangeListeners:   []IdealStateChangeListener{},
		keys: KeyBuilder{clusterID},
		stop: make(chan bool),
		externalViewResourceMap: map[string]bool{},
		idealStateResourceMap:   map[string]bool{},
		externalViewChanged:     make(chan string, 100),
		liveInstanceChanged:     make(chan string, 100),
		currentStateChanged:     make(chan string, 100),
		idealStateChanged:       make(chan string, 100),
		instanceConfigChanged:   make(chan string, 100),

		stopCurrentStateWatch: make(map[string]chan interface{}),

		currentStateChangeListenersLock: sync.Mutex{},
	}
}

func (m *HelixManager) NewParticipant(clusterID string, host string, port string) *Participant {
	return &Participant{
		ClusterID:     clusterID,
		Host:          host,
		Port:          port,
		ParticipantID: fmt.Sprintf("%s_%s", host, port),
		zkConnStr:     m.zkAddress,
		started:       make(chan interface{}),
		stop:          make(chan bool),
		stopWatch:     make(chan bool),
		keys:          KeyBuilder{clusterID},
	}
}
