package client

import (
	"bytes"
	"crypto/ecdh"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/AYM1607/ccclip/internal/server"
	"github.com/AYM1607/ccclip/pkg/crypto"
)

type Client struct {
	url string
}

func New(url string) *Client {
	return &Client{
		url: url,
	}
}

func (c *Client) Register(email, password string) error {
	req := server.RegisterRequest{
		Email:    email,
		Password: password,
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return err
	}
	res, err := http.Post(c.url+"/register", "application/json", bytes.NewReader(reqJson))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return errors.New("got unexpected response code from server")
	}

	resBody, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	log.Println(string(resBody))

	return nil
}

func (c *Client) RegisterDevice(email, password string, devicePublicKey []byte) (*server.RegisterDeviceResponse, error) {
	req := server.RegisterDeviceRequest{
		Email:     email,
		Password:  password,
		PublicKey: devicePublicKey,
	}
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	hres, err := http.Post(c.url+"/registerDevice", "application/json", bytes.NewReader(reqJson))
	if err != nil {
		return nil, err
	}
	if hres.StatusCode != http.StatusCreated {
		return nil, errors.New("got unexpected response code from server")
	}

	hresBody, err := io.ReadAll(hres.Body)
	defer hres.Body.Close()
	if err != nil {
		return nil, err
	}
	log.Println(string(hresBody))

	var res server.RegisterDeviceResponse
	err = json.Unmarshal(hresBody, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetDevices(deviceId string, pvk *ecdh.PrivateKey) (*server.GetUserDevicesResponse, error) {
	req := server.GetUserDevicesRequest{
		FingerPrint: server.FingerPrint{Timestamp: time.Now().UTC()},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	key := crypto.NewSharedKey(pvk, serverPublicKey, crypto.SendDirection)
	encryptedReq := crypto.Encrypt(key, reqBytes)

	authReq := server.AuthenticatedPayload{
		DeviceID: deviceId,
		Payload:  encryptedReq,
	}
	authReqJson, err := json.Marshal(authReq)
	if err != nil {
		return nil, err
	}

	hres, err := http.Post(c.url+"/userDevices", "application/json", bytes.NewReader(authReqJson))
	if err != nil {
		return nil, err
	}

	hresBody, err := io.ReadAll(hres.Body)
	defer hres.Body.Close()
	if err != nil {
		return nil, err
	}
	log.Println(string(hresBody))

	if hres.StatusCode != http.StatusOK {
		return nil, errors.New("got unexpected response code from server")
	}

	var res server.GetUserDevicesResponse
	err = json.Unmarshal(hresBody, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
