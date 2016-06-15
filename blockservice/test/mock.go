package bstest

import (
	. "github.com/RealImage/go-ipfs/blockservice"
	bitswap "github.com/RealImage/go-ipfs/exchange/bitswap"
	tn "github.com/RealImage/go-ipfs/exchange/bitswap/testnet"
	mockrouting "github.com/RealImage/go-ipfs/routing/mock"
	delay "github.com/RealImage/go-ipfs/thirdparty/delay"
)

// Mocks returns |n| connected mock Blockservices
func Mocks(n int) []*BlockService {
	net := tn.VirtualNetwork(mockrouting.NewServer(), delay.Fixed(0))
	sg := bitswap.NewTestSessionGenerator(net)

	instances := sg.Instances(n)

	var servs []*BlockService
	for _, i := range instances {
		servs = append(servs, New(i.Blockstore(), i.Exchange))
	}
	return servs
}
