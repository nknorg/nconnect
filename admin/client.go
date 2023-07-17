package admin

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/nknorg/nconnect/config"
	"github.com/nknorg/nkn-sdk-go"
)

const (
	replyTimeout = 10 * time.Second
)

var (
	ErrReplyTimeout = errors.New("wait for reply timeout")
)

var (
	serverAdminAddr string
)

type Client struct {
	*nkn.MultiClient
	ReplyTimeout time.Duration
}

func NewClient(account *nkn.Account, clientConfig *nkn.ClientConfig, identifier string) (*Client, error) {
	if identifier == "" {
		identifier = config.RandomIdentifier()
	}
	m, err := nkn.NewMultiClient(account, identifier, 4, false, clientConfig)
	if err != nil {
		return nil, err
	}

	c := &Client{
		MultiClient:  m,
		ReplyTimeout: replyTimeout,
	}

	<-m.OnConnect.C

	return c, nil
}

func (c *Client) RPCCall(addr, method string, params interface{}, result interface{}) error {
	req := map[string]interface{}{
		"id":     "nConnect",
		"method": method,
		"params": params,
	}

	reply, err := c.SendMsg(addr, req, true)
	if err != nil {
		return err
	}

	resp := &RpcResp{
		Result: result,
	}
	err = json.Unmarshal(reply.Data, resp)
	if err != nil {
		return err
	}

	if len(resp.Error) > 0 {
		return errors.New(resp.Error)
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

func (c *Client) SendMsg(address string, msg interface{}, waitResponse bool) (reply *nkn.Message, err error) {
	if c.ReplyTimeout == 0 {
		c.ReplyTimeout = replyTimeout
	}

	reqBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	var onReply *nkn.OnMessage
	for i := 0; i < 3; i++ {
		onReply, err = c.Send(nkn.NewStringArray(address), reqBytes, nkn.GetDefaultMessageConfig())
		if err != nil {
			return nil, err
		}

		if !waitResponse {
			return nil, nil
		}

		select {
		case reply = <-onReply.C:
			return reply, nil

		case <-time.After(c.ReplyTimeout):
			err = ErrReplyTimeout
		}
	}

	if err == ErrReplyTimeout {
		log.Printf("Wait for repsone timeout, please make sure the peer is running and reachable")
	}

	return nil, err
}
