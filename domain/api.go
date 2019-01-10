package domain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexisvisco/ovh-domain-api/domain/subsidiary"
	"io/ioutil"
	"net/http"
)

type ErrorKnown error

var (
	UnknownExtension ErrorKnown = errors.New("unknown extension")
	InvalidCartId    ErrorKnown = errors.New("invalid cart id")
)

type InfoDomain interface {
	GetOrderable() bool
	GetAction() string
}

type ResultDomain struct {
	Orderable      bool          `json:"orderable"`
	Offer          string        `json:"offer"`
	Phase          string        `json:"phase"`
	Configurations []interface{} `json:"configurations"`
	QuantityMax    int           `json:"quantityMax"`
	DeliveryTime   string        `json:"deliveryTime"`
	ProductID      string        `json:"productId"`
	Action         string        `json:"action"`
	PricingMode    string        `json:"pricingMode"`
	OfferID        string        `json:"offerId"`
	Duration       []string      `json:"duration"`
	Prices         []struct {
		Label string `json:"label"`
		Price struct {
			CurrencyCode string  `json:"currencyCode"`
			Value        float64 `json:"value"`
			Text         string  `json:"text"`
		} `json:"price"`
	} `json:"prices"`
}

type MinimalResultDomain struct {
	Orderable bool   `json:"orderable"`
	Action    string `json:"action"`
}

type FullResults []ResultDomain

type MinimalResults []MinimalResultDomain

type Client struct {
	Cart Cart
}

func NewClient(subsidiary subsidiary.List) (client *Client, err error) {
	cart, err := newCart(subsidiary)
	if err != nil {
		return nil, err
	}
	return &Client{
		Cart: *cart,
	}, nil
}

// InfoDomain return all information from ovh about a domain
func (c *Client) DomainInfo(domain string) (results FullResults, err error) {
	req, err := http.NewRequest("GET", "https://www.ovh.com/engine/apiv6/order/cart/"+c.Cart.CartId+"/domain?domain="+domain, nil)
	if err != nil {
		return FullResults{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return FullResults{}, err
	}
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return FullResults{}, err
	}
	v := &FullResults{}
	if bytes.Contains(by, []byte("Extension not managed")) {
		return *v, UnknownExtension
	} else if bytes.Contains(by, []byte("invalid cartId")) {
		return *v, InvalidCartId
	}
	err = json.Unmarshal(by, v)
	if err != nil {
		return FullResults{}, fmt.Errorf("%s message=%s", err.Error(), string(by))
	}
	defer resp.Body.Close()
	return *v, nil
}

// MinimalDomainInfo return only minimal information about a domain
func (c *Client) MinimalDomainInfo(domain string) (results MinimalResults, err error) {
	req, err := http.NewRequest("GET", "https://www.ovh.com/engine/apiv6/order/cart/"+c.Cart.CartId+"/domain?domain="+domain, nil)
	if err != nil {
		return MinimalResults{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return MinimalResults{}, err
	}
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return MinimalResults{}, err
	}
	v := &MinimalResults{}
	if bytes.Contains(by, []byte("Extension not managed")) {
		return *v, UnknownExtension
	} else if bytes.Contains(by, []byte("invalid cartId")) {
		return *v, InvalidCartId
	}
	err = json.Unmarshal(by, v)
	if err != nil {
		return MinimalResults{}, fmt.Errorf("%s message=%s", err.Error(), string(by))
	}
	defer resp.Body.Close()
	return *v, nil
}

func (f FullResults) IsTaken() bool {
	if len(f) == 0 {
		return true
	}
	for _, v := range f {
		if v.Action != "create" {
			return true
		}
		if !v.Orderable {
			return true
		}
	}
	return false
}

func (f MinimalResults) IsTaken() bool {
	if len(f) == 0 {
		return true
	}
	for _, v := range f {
		if v.Action != "create" {
			return true
		}
		if !v.Orderable {
			return true
		}
	}
	return false
}
