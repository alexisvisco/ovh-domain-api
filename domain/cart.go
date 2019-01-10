package domain

import (
	"encoding/json"
	"github.com/alexisvisco/ovh-domain-api/domain/subsidiary"
	"io/ioutil"
	"net/http"
	"strings"
)

type Cart struct {
	CartId string `json:"cartId"`
	Expire string `json:"expire"`
}

func newCart(subsidiary subsidiary.List) (cart *Cart, err error) {
	body := strings.NewReader("{\"ovhSubsidiary\": \"" + string(subsidiary) + "\"}")
	req, err := http.NewRequest("POST", "https://eu.api.ovh.com/1.0/order/cart", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	cart = &Cart{}
	err = json.Unmarshal(bytes, cart)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return cart, nil
}
