package shared

import (

	// "crypto/ed25519"

	"log"
	"pocket/consensus"
	"pocket/p2p/pre_p2p"
	"pocket/persistence"
	"pocket/shared/config"
	"pocket/shared/crypto"
	"pocket/utility"

	"pocket/shared/modules"
)

// TODO: SHould we create an interface for this as well?
type Node struct {
	modules.Module

	pocketBusMod modules.Bus

	Address string
}

func Create(config *config.Config) (n *Node, err error) {
	persistenceMod, err := persistence.Create(config)
	if err != nil {
		return nil, err
	}

	// TODO(derrandz): Uncomment this and comment out the pre_p2p.Create once the p2p module is ready for use
	//networkMod, err := p2p.Create(config)
	//if err != nil {
	//	return nil, err
	//}

	networkMod, err := pre_p2p.Create(config)
	if err != nil {
		return nil, err
	}

	utilityMod, err := utility.Create(config)
	if err != nil {
		return nil, err
	}

	consensusMod, err := consensus.Create(config)
	if err != nil {
		return nil, err
	}

	bus, err := CreateBus(nil, persistenceMod, networkMod, utilityMod, consensusMod)
	if err != nil {
		return nil, err
	}
	pk, err := crypto.NewPrivateKey(config.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &Node{
		pocketBusMod: bus,
		Address:      pk.Address().String(),
	}, nil
}

func (node *Node) Start() error {
	log.Println("Starting pocket node...")

	// NOTE: Order of module initialization matters.

	if err := node.GetBus().GetPersistenceModule().Start(); err != nil {
		return err
	}

	if err := node.GetBus().GetNetworkModule().Start(); err != nil {
		return err
	}

	if err := node.GetBus().GetUtilityModule().Start(); err != nil {
		return err
	}

	if err := node.GetBus().GetConsensusModule().Start(); err != nil {
		return err
	}

	for {
		event := node.GetBus().GetBusEvent()
		if err := node.handleEvent(event); err != nil {
			log.Println("Error handling event: ", err)
		}
	}
}

func (m *Node) SetBus(pocketBus modules.Bus) {
	m.pocketBusMod = pocketBus
}

func (m *Node) GetBus() modules.Bus {
	if m.pocketBusMod == nil {
		log.Fatalf("PocketBus is not initialized")
	}
	return m.pocketBusMod
}
