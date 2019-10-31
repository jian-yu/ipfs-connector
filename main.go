package main

import (
	"context"
	"fmt"
	"log"
	"os"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/interface-go-ipfs-core/path"

	ipfs "github.com/ipfs/go-ipfs-http-client"
	ma "github.com/multiformats/go-multiaddr"
)

const (
	ipfsGateway = "/ip4/172.38.8.89/tcp/9090/ipfs/QmcPVRJAKSQA17bH2kiBhCaFfQ5cHGjEBraD71VdGgSBYt"
)

func main() {
	Put()
}

func Get(ipfsPath string) {
	ctx := context.Background()
	multiAddr, err := ma.NewMultiaddr(ipfsGateway)
	if err != nil {
		log.Fatalf("generate multiaddr error by %s", err)
	}
	ipfsApi, err := ipfs.NewApi(multiAddr)
	if err != nil {
		log.Fatalf("ipfs http connect error by %s", err)
	}
	srcPath := path.New(ipfsPath)
	nodeObject, err := ipfsApi.Object().Get(ctx, srcPath)
	if err != nil {
		log.Fatalf("get object node fail by %s", err)
	}
	data, err := GetOne(ctx, nodeObject, ipfsApi)
	if err != nil {
		log.Fatalf("get object fail by %s", err)
	}
	fmt.Println(string(data))
}

func GetOne(ctx context.Context, node ipld.Node, api *ipfs.HttpApi) ([]byte, error) {
	var err error
	data := make([]byte, 0)
	data = append(data, node.RawData()[:]...)
	if len(node.Links()) == 0 {
		return data, err
	}
	for _, link := range node.Links() {
		srcPath := path.New(link.Cid.String())
		node, err := api.Object().Get(ctx, srcPath)
		if err != nil {
			log.Fatalf("get object link fail by %s", err)
		}
		data, err = GetOne(ctx, node, api)
		if err != nil {
			log.Fatalf("get link node fail by %s", err)
			break
		}
		data = append(data, data[:]...)
	}
	return data, err
}

func Put() {
	ctx := context.Background()
	multiAddr, err := ma.NewMultiaddr(ipfsGateway)
	if err != nil {
		log.Fatalf("generate multiaddr error by %s", err)
		return
	}
	ipfsApi, err := ipfs.NewApi(multiAddr)
	if err != nil {
		log.Fatalf("ipfs http connect error by %s", err)
		return
	}
	file, err := os.Open("README.md")
	if err != nil {
		log.Fatalf("open file error by %s", err)
		return
	}
	bStat, err := ipfsApi.Block().Put(ctx, file)
	if err != nil {
		log.Fatalf("put file into ipfs error by %s", err)
		return
	}
	log.Printf("file store path is %s and size is %d", bStat.Path().String(), bStat.Size())
}
