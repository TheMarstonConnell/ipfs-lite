package main

// This example launches an IPFS-Lite peer and fetches a hello-world
// hash from the IPFS network.

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	ipfslite "github.com/hsanjuan/ipfs-lite"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ds := ipfslite.NewInMemoryDatastore()
	priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		panic(err)
	}

	listen, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4005")

	h, dht, err := ipfslite.SetupLibp2p(
		ctx,
		priv,
		nil,
		[]multiaddr.Multiaddr{listen},
		ds,
		ipfslite.Libp2pOptionsExtra...,
	)

	if err != nil {
		panic(err)
	}

	lite, err := ipfslite.New(ctx, ds, nil, h, dht, nil)
	if err != nil {
		panic(err)
	}

	lite.Bootstrap(ipfslite.DefaultBootstrapPeers())

	f := []byte("this is a chunk test, there are a few chunks here to test out!")
	buf := bytes.NewReader(f)
	n, err := lite.AddFile(context.Background(), buf, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("CID: %s\n", n.Cid())

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done // Will block here until user hits ctrl+c
}
