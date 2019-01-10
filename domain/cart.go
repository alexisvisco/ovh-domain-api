package domain

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/alexisvisco/ovh-domain-api/domain/subsidiary"
	"io/ioutil"
	"net/http"
	"strings"
)

type CartError error

var (
	InvalidSubsidiaryFormat CartError = errors.New("parameter ovhSubsidiary isn't formated correctly")
)

type Cart struct {
	CartId string `json:"cartId"`
	Expire string `json:"expire"`
}

func newCart(subsidiary subsidiary.List) (cart *Cart, err error) {
	body := strings.NewReader("{\"ovhSubsidiary\": \"" + string(subsidiary) + "\"}")
	req, err := http.NewRequest("POST", "https://www.ovh.com/engine/apiv6/order/cart", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(by, []byte("Parameter ovhSubsidiary isn't formated correctly")) {
		return nil, InvalidSubsidiaryFormat
	}
	cart = &Cart{}
	err = json.Unmarshal(by, cart)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return cart, nil
}
