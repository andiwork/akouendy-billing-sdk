package billing

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

// RequestBeforeFn  is the function signature for the RequestBefore callback function
type RequestBeforeFn func(ctx context.Context, req *http.Request) error

// ResponseAfterFn  is the function signature for the ResponseAfter callback function
type ResponseAfterFn func(ctx context.Context, rsp *http.Response) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	Endpoint string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A callback for modifying requests which are generated before sending over
	// the network.
	RequestBefore RequestBeforeFn

	// A callback for modifying response which are generated before sending over
	// the network.
	ResponseAfter ResponseAfterFn

	// The user agent header identifies your application, its version number, and the platform and programming language you are using.
	// You must include a user agent header in each request submitted to the sales partner API.
	UserAgent string
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Endpoint: billingConfig.billingUrl,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the endpoint URL always has a trailing slash
	if !strings.HasSuffix(client.Endpoint, "/") {
		client.Endpoint += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = http.DefaultClient
	}
	// setting the default useragent
	if client.UserAgent == "" {
		client.UserAgent = fmt.Sprintf("akouendy-billing-api-sdk/v1.0")
	}

	return &client, nil
}

// WithRequestBefore allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestBefore(fn RequestBeforeFn) ClientOption {
	return func(c *Client) error {
		c.RequestBefore = fn
		return nil
	}
}

// WithResponseAfter allows setting up a callback function, which will be
// called right after get response the request. This can be used to log.
func WithResponseAfter(fn ResponseAfterFn) ClientOption {
	return func(c *Client) error {
		c.ResponseAfter = fn
		return nil
	}
}

// WithUserAgent set up useragent
// add user agent to every request automatically
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	CreateOrder(ctx context.Context, transactionId string, body OrderRequest) (orderResponse OrderResponse, billingTrx BillingTransaction, err error)
	GetOrderStatus(ctx context.Context, body OrderSubsRequest) (orderSubsResponse OrderSubsReponse, err error)
	CreatePayment(ctx context.Context, transactionId string, body PaymentRequest) (paymentResponse PaymentResponse, billingTrx BillingTransaction, err error)
	GetPaymentStatus(ctx context.Context, paymentToken string) (paymentStatusResponse PaymentStatusResponse, err error)
}

func (c *Client) CreateOrder(ctx context.Context, transactionId string, body OrderRequest) (orderResponse OrderResponse, billingTrx BillingTransaction, err error) {
	// webhook url
	webUrl, err := url.Parse(billingConfig.AppBaseUrl + "/2021-10-01/billing-webhook/" + transactionId)
	body.Webhook = webUrl.String()

	queryUrl, err := url.Parse(c.Endpoint)
	if err != nil {
		return
	}
	basePath := fmt.Sprintf("/order/create")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return
	}

	buf, err := json.Marshal(body)
	if err != nil {
		return
	}
	bodyReader := bytes.NewReader(buf)

	req, err := http.NewRequest("POST", queryUrl.String(), bodyReader)
	if err != nil {
		log.Println("NewRequest Error", err)
	}

	req.Header.Add("Content-Type", "application/json")

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", c.UserAgent)
	if c.RequestBefore != nil {
		err = c.RequestBefore(ctx, req)
		if err != nil {
			return
		}
	}
	// Debug request
	if billingConfig.Debug {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Println("DumpRequest Error", err)
		}
		log.Printf("DumpRequest = %s", dump)
	}
	rsp, err := c.Client.Do(req)
	if err != nil {
		return
	}

	if billingConfig.Debug {
		dump, err := httputil.DumpResponse(rsp, true)
		if err != nil {
			log.Println(err, "DumpResponse Error")
		}
		log.Printf("DumpResponse = %s", dump)
	}

	if c.ResponseAfter != nil {
		err = c.ResponseAfter(ctx, rsp)
		if err != nil {
			return
		}
	}

	defer rsp.Body.Close()
	bodyByte, err := ioutil.ReadAll(rsp.Body) // response body is []byte

	if rsp.StatusCode >= 200 && rsp.StatusCode <= 299 {
		if err := json.Unmarshal(bodyByte, &orderResponse); err != nil { // Parse []byte to the go struct pointer
			log.Println("Can not unmarshal JSON")
		} else {
			billingTrx.OrderID = orderResponse.OrderID
			billingTrx.PaymentToken = orderResponse.PaymentToken
			billingTrx.AppTrxId = transactionId
		}

		return
	} else {
		err = errors.New(fmt.Sprintf("%s", bodyByte))
	}

	return
}

func (c *Client) GetOrderStatus(ctx context.Context, body OrderSubsRequest) (orderSubsResponse OrderSubsReponse, err error) {

	queryUrl, err := url.Parse(c.Endpoint)
	if err != nil {
		return
	}
	basePath := fmt.Sprintf("/order/check")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return
	}

	buf, err := json.Marshal(body)
	if err != nil {
		return
	}
	bodyReader := bytes.NewReader(buf)

	req, err := http.NewRequest("POST", queryUrl.String(), bodyReader)
	if err != nil {
		log.Println("NewRequest Error", err)
	}

	req.Header.Add("Content-Type", "application/json")

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", c.UserAgent)
	if c.RequestBefore != nil {
		err = c.RequestBefore(ctx, req)
		if err != nil {
			return
		}
	}
	// Debug request
	if billingConfig.Debug {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Println("DumpRequest Error", err)
		}
		log.Printf("DumpRequest = %s", dump)
	}
	rsp, err := c.Client.Do(req)
	if err != nil {
		return
	}

	if billingConfig.Debug {
		dump, err := httputil.DumpResponse(rsp, true)
		if err != nil {
			log.Println(err, "DumpResponse Error")
		}
		log.Printf("DumpResponse = %s", dump)
	}

	if c.ResponseAfter != nil {
		err = c.ResponseAfter(ctx, rsp)
		if err != nil {
			return
		}
	}

	defer rsp.Body.Close()
	bodyByte, err := ioutil.ReadAll(rsp.Body) // response body is []byte

	if rsp.StatusCode >= 200 && rsp.StatusCode <= 299 {
		if err := json.Unmarshal(bodyByte, &orderSubsResponse); err != nil { // Parse []byte to the go struct pointer
			log.Println("Can not unmarshal JSON")
		}
		return
	} else {
		err = errors.New(fmt.Sprintf("%s", bodyByte))
	}

	return
}

func (c *Client) CreatePayment(ctx context.Context, body PaymentRequest) (paymentResponse PaymentResponse, billingTrx BillingTransaction, err error) {
	// webhook url
	webUrl, err := url.Parse(billingConfig.AppBaseUrl + "/2023-05-03/payment-webhook/")
	body.Webhook = webUrl.String()
	str := body.AppId + "|" + body.TransactionId + "|" + strconv.Itoa(body.TotalAmount) + "|akouna_matata"
	body.Hash = hash512(str)
	queryUrl, err := url.Parse(c.Endpoint)
	if err != nil {
		return
	}
	basePath := fmt.Sprintf("billing/payment/init")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return
	}

	buf, err := json.Marshal(body)
	if err != nil {
		return
	}
	bodyReader := bytes.NewReader(buf)

	req, err := http.NewRequest("POST", queryUrl.String(), bodyReader)
	if err != nil {
		log.Println("NewRequest Error", err)
	}

	req.Header.Add("Content-Type", "application/json")

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", c.UserAgent)
	if c.RequestBefore != nil {
		err = c.RequestBefore(ctx, req)
		if err != nil {
			return
		}
	}
	// Debug request
	if billingConfig.Debug {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Println("DumpRequest Error", err)
		}
		log.Printf("DumpRequest = %s", dump)
	}
	rsp, err := c.Client.Do(req)
	if err != nil {
		return
	}

	if billingConfig.Debug {
		dump, err := httputil.DumpResponse(rsp, true)
		if err != nil {
			log.Println(err, "DumpResponse Error")
		}
		log.Printf("DumpResponse = %s", dump)
	}

	if c.ResponseAfter != nil {
		err = c.ResponseAfter(ctx, rsp)
		if err != nil {
			return
		}
	}

	defer rsp.Body.Close()
	bodyByte, err := ioutil.ReadAll(rsp.Body) // response body is []byte

	if rsp.StatusCode >= 200 && rsp.StatusCode <= 299 {
		if err := json.Unmarshal(bodyByte, &paymentResponse); err != nil { // Parse []byte to the go struct pointer
			log.Println("Can not unmarshal JSON")
		} else {
			billingTrx.PaymentToken = paymentResponse.Token
			billingTrx.AppTrxId = body.TransactionId
		}

		return
	} else {
		err = errors.New(fmt.Sprintf("%s", bodyByte))
	}

	return
}

func (c *Client) GetPaymentStatus(ctx context.Context, paymentToken string) (paymentStatusResponse PaymentStatusResponse, err error) {

	queryUrl, err := url.Parse(c.Endpoint)
	if err != nil {
		return
	}
	basePath := fmt.Sprintf("/payment/%s",paymentToken)
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return
	}



	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		log.Println("NewRequest Error", err)
	}

	req.Header.Add("Content-Type", "application/json")

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", c.UserAgent)
	if c.RequestBefore != nil {
		err = c.RequestBefore(ctx, req)
		if err != nil {
			return
		}
	}
	// Debug request
	if billingConfig.Debug {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Println("DumpRequest Error", err)
		}
		log.Printf("DumpRequest = %s", dump)
	}
	rsp, err := c.Client.Do(req)
	if err != nil {
		return
	}

	if billingConfig.Debug {
		dump, err := httputil.DumpResponse(rsp, true)
		if err != nil {
			log.Println(err, "DumpResponse Error")
		}
		log.Printf("DumpResponse = %s", dump)
	}

	if c.ResponseAfter != nil {
		err = c.ResponseAfter(ctx, rsp)
		if err != nil {
			return
		}
	}

	defer rsp.Body.Close()
	bodyByte, err := ioutil.ReadAll(rsp.Body) // response body is []byte

	if rsp.StatusCode >= 200 && rsp.StatusCode <= 299 {
		if err := json.Unmarshal(bodyByte, &paymentStatusResponse); err != nil { // Parse []byte to the go struct pointer
			log.Println("Can not unmarshal JSON")
		}
		return
	} else {
		err = errors.New(fmt.Sprintf("%s", bodyByte))
	}

	return
}
