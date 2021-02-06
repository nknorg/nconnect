package admin

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/nknorg/nconnect/config"
	"github.com/nknorg/nkn-sdk-go"
)

const (
	replyTimeout = 10 * time.Second
)

var (
	errReplyTimeout = errors.New("wait for reply timeout")
)

var (
	serverAdminAddr string
)

type Client struct {
	*nkn.MultiClient
	replyTimeout time.Duration
}

func NewClient(account *nkn.Account, clientConfig *nkn.ClientConfig) (*Client, error) {
	m, err := nkn.NewMultiClient(account, config.RandomIdentifier(), 4, false, clientConfig)
	if err != nil {
		return nil, err
	}

	c := &Client{
		MultiClient:  m,
		replyTimeout: replyTimeout,
	}

	<-m.OnConnect.C

	return c, nil
}

func (c *Client) RPCCall(addr, method string, params map[string]interface{}, result interface{}) error {
	req, err := json.Marshal(map[string]interface{}{
		"id":     "nConnect",
		"method": method,
		"params": params,
	})
	if err != nil {
		return err
	}

	onReply, err := c.Send(nkn.NewStringArray(addr), req, nil)
	if err != nil {
		return err
	}

	var reply *nkn.Message
	select {
	case reply = <-onReply.C:
	case <-time.After(c.replyTimeout):
		return errReplyTimeout
	}

	resp := make(map[string]*json.RawMessage)
	err = json.Unmarshal(reply.Data, &resp)
	if err != nil {
		return err
	}
	if resp["error"] != nil {
		var errStr string
		err = json.Unmarshal(*resp["error"], &errStr)
		if err != nil {
			return err
		}
		return errors.New(errStr)
	}

	err = json.Unmarshal(*resp["result"], result)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetInfo(addr string) (*GetInfoJSON, error) {
	res := &GetInfoJSON{}
	err := c.RPCCall(addr, "getInfo", nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
