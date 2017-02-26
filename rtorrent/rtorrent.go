package rtorrent

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/kolo/xmlrpc"
)

type RtClient struct {
	Transport http.RoundTripper
	Client    *xmlrpc.Client
}

type Torrent struct {
	Name               string
	Hash               string
	TotalSizeInBytes   int64
	CompletedBytes     int64
	UpRate             int64
	DownRate           int64
	State              int64
	TotalUploadedBytes int64
}

type ViewType string

func (rt *RtClient) Call(method string, args interface{}, reply interface{}) (err error) {
	err = rt.Client.Call(method, args, reply)
	if err != nil {
		log.Println("Error occured:", err)
	}
	return
}

func Client(transport http.RoundTripper) (c *RtClient, err error) {
	c = &RtClient{
		Transport: transport,
	}
	c.Client, err = xmlrpc.NewClient("/RPC2", transport)
	if err != nil {
		log.Println("Error occured:", err)
	}
	return
}

func (rt *RtClient) GetList(viewName string) (items []Torrent, err error) {
	var result []interface{}
	args := []string{viewName,
		"d.name=",
		"d.hash=",
		"d.size_bytes=",
		"d.completed_bytes=",
		"d.up.rate=",
		"d.down.rate=",
		"d.state=",
		"d.up.total="}
	err = rt.Call("d.multicall", args, &result)
	for _, item := range result {
		var torrent Torrent
		array := item.([]interface{})
		torrent.Name = array[0].(string)
		torrent.Hash = array[1].(string)
		torrent.TotalSizeInBytes = array[2].(int64)
		torrent.CompletedBytes = array[3].(int64)
		torrent.UpRate = array[4].(int64)
		torrent.DownRate = array[5].(int64)
		torrent.State = array[6].(int64)
		torrent.TotalUploadedBytes = array[7].(int64)

		items = append(items, torrent)
	}

	return
}

func (rt *RtClient) Erase(hash string) (err error) {
	var result int8
	if err = rt.Stop(hash); err == nil {
		err = rt.Call("d.erase", hash, &result)
	}
	return
}

func (rt *RtClient) Stop(hash string) (err error) {
	var result int8
	err = rt.Call("d.stop", hash, &result)
	return
}

func (rt *RtClient) Start(hash string) (err error) {
	var result int8
	err = rt.Call("d.start", hash, &result)
	return
}

func (rt *RtClient) LoadRaw(data []byte, tag string) (err error) {
	var result int8
	args := []interface{}{
		"",
		xmlrpc.Base64(base64.StdEncoding.EncodeToString(data)),
	}
	if len(tag) > 0 {
		args = append(args, fmt.Sprintf("d.custom1.set=%s/", tag))
	}
	err = rt.Call("load.raw_verbose", args, &result)
	return
}
