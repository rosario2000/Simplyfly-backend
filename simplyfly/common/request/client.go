package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	hc   *HttpClient
	once sync.Once
)

type Response struct {
	*http.Response
	Body []byte
}

type HttpClient struct {
	Client *http.Client
}

type Request struct {
	Method string
	URL    string
	Params url.Values
	Header http.Header
	Body   interface{}
}

func Instance() *HttpClient {
	once.Do(func() {
		MaxIdleConns := 100
		transport := &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConns,
			MaxConnsPerHost:     MaxIdleConns,
			IdleConnTimeout:     90 * time.Second,
			//TLSHandshakeTimeout:   10 * time.Second,
			//ExpectContinueTimeout: 4 * time.Second,
			//ResponseHeaderTimeout: 8 * time.Second,
		}
		hc = &HttpClient{Client: &http.Client{
			Timeout:   120 * time.Second,
			Transport: transport,
		}}
		return
	})
	return hc
}

func (h *HttpClient) Request(ctx context.Context, method string, urlString string, header map[string]string, reqBody interface{}) (*Response, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, urlString, body)

	if header != nil {
		for k, v := range header {
			req.Header[k] = []string{v}
		}
	}
	if req.Header.Get(echo.HeaderContentType) == "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	if err != nil {
		return nil, err
	}
	return h.Do(ctx, req)
}
func (h *HttpClient) GetReqBody(ctx context.Context, urlString string, header map[string]string, reqBody interface{}) (*Response, error) {
	q, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	rBody := strings.NewReader(string(payload))
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, "GET", q.String(), rBody)
	if err != nil {
		return nil, err
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	return h.Do(ctx, req)
}

func (h *HttpClient) Get(ctx context.Context, urlString string, header map[string]string, params url.Values) (*Response, error) {
	q, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	if params != nil {
		q.RawQuery = params.Encode()
	}
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, "GET", q.String(), nil)
	if err != nil {
		return nil, err
	}

	if header != nil {
		for k, v := range header {
			//req.Header.Set(k, v)
			req.Header[k] = []string{v} // used this line to use case-sensitive headers.
		}
	}
	return h.Do(ctx, req)
}

func (h *HttpClient) Post(ctx context.Context, urlString string, header map[string]string, reqBody interface{}) (req *Response, err error) {
	return h.Request(ctx, http.MethodPost, urlString, header, reqBody)
}

func (h *HttpClient) Put(ctx context.Context, urlString string, header map[string]string, reqBody interface{}) (*Response, error) {
	return h.Request(ctx, http.MethodPut, urlString, header, reqBody)
}

func (h *HttpClient) Delete(ctx context.Context, urlString string, header map[string]string, reqBody interface{}) (req *Response, err error) {
	return h.Request(ctx, http.MethodDelete, urlString, header, reqBody)
}

func (h *HttpClient) PostForm(ctx context.Context, urlString string, header map[string]string, reqBody interface{}) (_ *Response, err error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	r, err := json.Marshal(reqBody)
	if err != nil {
		return
	}

	err = writer.WriteField("body", string(r))
	if err != nil {
		return
	}

	err = writer.Close()
	if err != nil {
		return
	}
	//body := strings.NewReader(params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlString, payload)
	if err != nil {
		return
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return h.Do(ctx, req)
}

func (h *HttpClient) PostForm2(ctx context.Context, urlString string, header map[string]string, reqBody map[string]string) (_ *Response, err error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	//r, err := json.Marshal(reqBody)
	//if err != nil {
	//	return
	//}

	for key, value := range reqBody {
		err = writer.WriteField(key, value)
		if err != nil {
			return
		}
	}

	err = writer.Close()
	if err != nil {
		return
	}
	//body := strings.NewReader(params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlString, payload)
	if err != nil {
		return
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return h.Do(ctx, req)
}

func (h *HttpClient) PostFormForVin(ctx context.Context, urlString string, header map[string]string, reqBody string) (_ *Response, err error) {
	//payload := &bytes.Buffer{}
	//writer := multipart.NewWriter(payload)

	payload := strings.NewReader(reqBody)
	//err = writer.WriteField("body", string(r))
	if err != nil {
		return
	}

	//err = writer.Close()
	if err != nil {
		return
	}
	//body := strings.NewReader(params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlString, payload)
	if err != nil {
		return
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	//req.Header.Set("Content-Type", writer.FormDataContentType())
	return h.Do(ctx, req)
}

func (h *HttpClient) Do(ctx context.Context, req *http.Request) (res *Response, err error) {

	// Commented by Amardeep
	//var reqBody []byte
	//if req.Body != nil {
	//	var bErr error
	//	reqBody, bErr = GetBody(req)
	//	if bErr != nil {
	//		log.Errorj(log.JSON{
	//			"message": "couldn't get request body",
	//			"err":     bErr.Error(),
	//		})
	//	}
	//}
	var resp *http.Response

	resp, err = h.Client.Do(req)

	routeStr, _ := url.Parse(req.URL.String())

	if strings.Contains(req.URL.String(), "CallHttpsApi") {
		urlHeaderValue := ""
		if urlSlice, ok := req.Header["url"]; ok && len(urlSlice) > 0 {
			urlHeaderValue = strings.TrimSpace(urlSlice[0])
		}

		newRouteStr := FetchURL(urlHeaderValue)
		if newRouteStr != nil && newRouteStr.Host != "" {
			routeStr = newRouteStr
		}
	}

	newURL := &url.URL{
		Scheme:   routeStr.Scheme,
		Host:     routeStr.Host,
		Path:     routeStr.Path,
		Fragment: routeStr.Fragment,
	}

	// Replace numeric part in the path with "xxxx"
	newPath := HideNumericPart(newURL.Path)

	newPath = HideTataVariablePart(newPath)

	// Create a new URL with the modified path
	newURL.Path = newPath

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// Commented by Amardeep
	//allowedEndpoints := map[string]bool{
	//	constants.CreateOrder:          true,
	//	constants.UpdateOrderStatus:    true,
	//	constants.CreateDeliveryOrder:  true,
	//	constants.UpdateDeliveryStatus: true,
	//	constants.CancelDeliveryOrder:  true,
	//	constants.CheckServiceability:  true,
	//}
	//if strings.EqualFold(req.Method, http.MethodPost) && allowedEndpoints[endpoint] {
	//	orderID, _ := ctx.Value(constants.OrderID).(int)
	//	bq.WriteOrderData(ctx, &bq.OrderRequestData{
	//		CreatedAt: time.Now(),
	//		Tag: strings.ToUpper(fmt.Sprintf(
	//			"oms_partners.%s.%s.%s",
	//			os.Getenv(constants.OMSPartnersENV),
	//			strings.ReplaceAll(strings.Trim(endpoint, "/"), "/", "_"),
	//			partner,
	//		)),
	//		OrderId:          int64(orderID),
	//		Request:          string(reqBody),
	//		Response:         string(body),
	//		HttpResponseCode: resp.StatusCode,
	//		Url:              req.URL.String(),
	//	})
	//}
	return &Response{
		Response: resp,
		Body:     body,
	}, nil
}

func (h *HttpClient) Send(ctx context.Context, r *Request) (res *Response, err error) {
	urlString := r.URL
	reqBody := bytes.NewBuffer(nil)
	if r.Params != nil {
		q, err := url.Parse(r.URL)
		if err != nil {
			return nil, err
		}
		q.RawQuery = r.Params.Encode()
		urlString = q.String()
	}
	if r.Body != nil {
		err := json.NewEncoder(reqBody).Encode(r.Body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, r.Method, urlString, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header
	if req.Header.Get(echo.HeaderContentType) == "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	return h.Do(ctx, req)
}

func GetBody(req *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}
func DecodeBody(req *http.Request, v interface{}) (err error) {
	body, err := GetBody(req)
	if err != nil {
		return
	}
	return json.Unmarshal(body, v)
}

func (h *HttpClient) PostONDC(ctx context.Context, urlString string, header map[string]string, reqBody string) (_ *Response, err error) {
	//payload := &bytes.Buffer{}
	//writer := multipart.NewWriter(payload)

	payload := strings.NewReader(reqBody)
	//err = writer.WriteField("body", string(r))
	if err != nil {
		return
	}

	//err = writer.Close()
	if err != nil {
		return
	}
	//body := strings.NewReader(params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlString, payload)
	if err != nil {
		return
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	//req.Header.Set("Content-Type", writer.FormDataContentType())
	return h.Do(ctx, req)
}

func (h *HttpClient) SendPostRequestWithFile(ctx context.Context, url string, filePath string) (_ *Response, err error) {
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open(filePath)
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("file", filepath.Base(filePath))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println("Error while creating form file: ", errFile1)
		return
	}
	err = writer.Close()
	if err != nil {
		fmt.Println("Error while writing to file to body: ", err)
		return
	}

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println("Error while making request: ", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return h.Do(ctx, req)
}

func HideTataVariablePart(path string) string {

	if !strings.Contains(path, "api/v2/sso/access-token") {
		return path
	}
	// Find the index of "token/"
	tokenIndex := strings.Index(path, "token/")
	if tokenIndex == -1 {
		return path // If "token/" is not found, return original path
	}

	// Construct the replacement string
	replacement := "xxxx"

	// Concatenate the parts of the path before "token/" and the replacement string
	return path[:tokenIndex+len("token/")] + replacement
}

func HideNumericPart(path string) string {
	// Use regular expression to find numeric parts
	re := regexp.MustCompile(`(\d{3,})`) // Match numeric parts with 3 or more characters
	matches := re.FindAllString(path, -1)

	// If there are matching numeric parts, replace them with "xxxx"
	if len(matches) > 0 {
		for _, match := range matches {
			path = strings.Replace(path, match, "xxxx", 1)
		}
	}

	return path
}

func FetchURL(urlStr string) *url.URL {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil
	}
	return parsedURL
}
