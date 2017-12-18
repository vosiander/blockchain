package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"bytes"

	"github.com/siklol/blockchain"
)

func Version(p *Peer) (string, error) {
	v := struct {
		Version string `json:"version"`
	}{}

	rsp, err := get(p, "/version")
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(rsp, &v); err != nil {
		return "", err
	}

	return v.Version, nil
}

func GenesisBlock(p *Peer) (*blockchain.Block, error) {
	rsp, err := get(p, "/blocks/genesis")
	if err != nil {
		return nil, err
	}

	return block(rsp)
}

func Tip(p *Peer) (*blockchain.Block, error) {
	rsp, err := get(p, "/blocks/tip")
	if err != nil {
		return nil, err
	}

	return block(rsp)
}

func BlockAtIndex(p *Peer, i int64) (*blockchain.Block, error) {
	rsp, err := get(p, fmt.Sprintf("/blocks/index/%d", i))
	if err != nil {
		return nil, err
	}

	return block(rsp)
}

func AddPeer(h *Peer, newPeer *Peer) error {
	jsonD, _ := json.Marshal(newPeer)
	_, err := post(h, "/peers", jsonD)

	return err
}

func block(rsp []byte) (*blockchain.Block, error) {
	var block *blockchain.Block
	if err := json.Unmarshal(rsp, &block); err != nil {
		return nil, err
	}

	return block, nil
}

func get(p *Peer, url string) ([]byte, error) {
	response, err := http.Get(fmt.Sprintf("http://%s:%s"+url, p.IP, p.Port)) // TODO https?
	if err != nil {
		return []byte(""), err
	}

	rsp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte(""), err
	}

	return rsp, nil
}

func post(p *Peer, url string, jsonStr []byte) ([]byte, error) {
	response, err := http.Post(fmt.Sprintf("http://%s:%s"+url, p.IP, p.Port), "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return []byte(""), err
	}

	rsp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte(""), err
	}

	return rsp, nil
}
