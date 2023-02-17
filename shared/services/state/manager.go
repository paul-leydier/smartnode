package state

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/smartnode/shared/services/beacon"
	"github.com/rocket-pool/smartnode/shared/services/config"
	cfgtypes "github.com/rocket-pool/smartnode/shared/types/config"
	"github.com/rocket-pool/smartnode/shared/utils/log"
)

type NetworkStateManager struct {
	cfg          *config.RocketPoolConfig
	rp           *rocketpool.RocketPool
	ec           rocketpool.ExecutionClient
	bc           beacon.Client
	log          *log.ColorLogger
	Config       *config.RocketPoolConfig
	Network      cfgtypes.Network
	ChainID      uint
	BeaconConfig beacon.Eth2Config
	//latestState  *NetworkState
	updateLock sync.Mutex
}

// Create a new manager for the network state
func NewNetworkStateManager(rp *rocketpool.RocketPool, cfg *config.RocketPoolConfig, ec rocketpool.ExecutionClient, bc beacon.Client, log *log.ColorLogger) (*NetworkStateManager, error) {

	// Create the manager
	m := &NetworkStateManager{
		cfg:     cfg,
		rp:      rp,
		ec:      ec,
		bc:      bc,
		log:     log,
		Config:  cfg,
		Network: cfg.Smartnode.Network.Value.(cfgtypes.Network),
		ChainID: cfg.Smartnode.GetChainID(),
	}

	// Get the Beacon config info
	var err error
	m.BeaconConfig, err = m.bc.GetEth2Config()
	if err != nil {
		return nil, err
	}

	return m, nil

}

// Get the state of the network using the latest Execution layer block
func (m *NetworkStateManager) GetHeadState() (*NetworkState, error) {
	// Get the latest EL block
	latestBlockHeader, err := m.ec.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error getting latest EL block: %w", err)
	}

	// Get the corresponding Beacon slot based on the timestamp
	latestBlockTime := time.Unix(int64(latestBlockHeader.Time), 0)
	genesisTime := time.Unix(int64(m.BeaconConfig.GenesisTime), 0)
	secondsSinceGenesis := uint64(latestBlockTime.Sub(genesisTime).Seconds())
	targetSlot := secondsSinceGenesis / m.BeaconConfig.SecondsPerSlot

	// Return
	return m.getState(targetSlot)
}

// Get the state of the network at the provided Beacon slot
func (m *NetworkStateManager) GetStateForSlot(slotNumber uint64) (*NetworkState, error) {
	return m.getState(slotNumber)
}

// Get the state of the network at the provided Beacon slot
func (m *NetworkStateManager) getState(slotNumber uint64) (*NetworkState, error) {
	state, err := CreateNetworkState(m.cfg, m.rp, m.ec, m.bc, m.log, slotNumber, m.BeaconConfig)
	if err != nil {
		return nil, err
	}

	m.updateLock.Lock()
	defer m.updateLock.Unlock()
	return state, nil
}

// Logs a line if the logger is specified
func (m *NetworkStateManager) logLine(format string, v ...interface{}) {
	if m.log != nil {
		m.log.Printlnf(format, v)
	}
}

// Gets the latest valid finalized slot from the given Beacon Chain head info
func (m *NetworkStateManager) GetLatestFinalizedBeaconBlock() (beacon.BeaconBlock, error) {
	head, err := m.bc.GetBeaconHead()
	if err != nil {
		return beacon.BeaconBlock{}, fmt.Errorf("error getting Beacon chain head: %w", err)
	}
	targetSlot := head.FinalizedEpoch*m.BeaconConfig.SlotsPerEpoch + (m.BeaconConfig.SlotsPerEpoch - 1)

	// If that slot is missing, get the latest one that isn't
	for {
		// Try to get the current block
		block, exists, err := m.bc.GetBeaconBlock(fmt.Sprint(targetSlot))
		if err != nil {
			return beacon.BeaconBlock{}, fmt.Errorf("error getting Beacon block %d: %w", targetSlot, err)
		}

		// If the block was missing, try the previous one
		if !exists {
			m.logLine("Slot %d was missing, trying the previous one...", targetSlot)
			targetSlot--
		} else {
			return block, nil
		}
	}
}
