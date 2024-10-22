package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

type ResponseData[T any] struct {
	Results struct {
		Items map[string]T `json:"-"`
	} `json:"results"`
	More bool `json:"more"`
}

func (rd *ResponseData[T]) DecodeBody(r io.Reader, fieldName string) error {
	var rawResults struct {
		Results map[string]json.RawMessage `json:"results"`
		More    bool                       `json:"more"`
	}
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&rawResults)
	if err != nil {
		return err
	}
	rd.More = rawResults.More

	itemsData, ok := rawResults.Results[fieldName]
	if !ok {
		return fmt.Errorf("expected field '%s' not found in results", fieldName)
	}
	var items map[string]T
	err = json.Unmarshal(itemsData, &items)
	if err != nil {
		return err
	}
	rd.Results.Items = items
	return nil
}

func (rd *ResponseData[T]) ExtractItems() ([]T, bool) {
	var items []T
	for _, item := range rd.Results.Items {
		items = append(items, item)
	}
	return items, rd.More
}

func APIRequest[Params any, Body any, Res any](params *Params, body *Body, method, URL, token, fieldName string) ([]Res, bool, *RequestError) {
	baseURL, err := url.Parse(URL)
	if err != nil {
		return nil, false, NewRequestError(http.StatusInternalServerError, fmt.Errorf("error parsing base URL: %w", err), false)
	}

	if params != nil {
		queryParams, err := query.Values(params)
		if err != nil {
			return nil, false, NewRequestError(http.StatusInternalServerError, fmt.Errorf("error extracting query parameters: %w", err), false)
		}
		baseURL.RawQuery = queryParams.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, false, NewRequestError(http.StatusInternalServerError, fmt.Errorf("error marshaling body to JSON: %w", err), false)
		}
		bodyReader = bytes.NewReader(bodyJSON)
	}

	req, err := http.NewRequest(method, baseURL.String(), bodyReader)
	if err != nil {
		return nil, false, NewRequestError(http.StatusInternalServerError, fmt.Errorf("error creating request: %w", err), false)
	}

	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, false, NewRequestError(http.StatusInternalServerError, fmt.Errorf("error executing request: %w", err), false)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		if res.StatusCode == 429 {
			return nil, false, NewRequestError(res.StatusCode, fmt.Errorf("rate limit reached"), true)
		}
		return nil, false, NewRequestError(res.StatusCode, fmt.Errorf("%s request error", method), false)
	}

	var response ResponseData[Res]
	err = response.DecodeBody(res.Body, fieldName)
	if err != nil {
		return nil, false, NewRequestError(http.StatusInternalServerError, fmt.Errorf("unable to decode response: %w", err), false)
	}

	items, more := response.ExtractItems()
	return items, more, NewRequestError(http.StatusOK, nil, false)
}

func GetData[Req any, Res any](params *Req, URL, token string, fieldName string) ([]Res, bool, *RequestError) {
	baseURL, err := url.Parse(URL)
	if err != nil {
		return nil, false, NewRequestError(http.StatusInternalServerError, fmt.Errorf("error parsing base URL: %w", err), false)
	}

	queryParams, err := query.Values(params)
	if err != nil {
		return nil, false, NewRequestError(http.StatusInternalServerError, fmt.Errorf("error extracting query parameters: %w", err), false)
	}

	baseURL.RawQuery = queryParams.Encode()
	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return nil, false, NewRequestError(http.StatusInternalServerError, fmt.Errorf("error creating request: %w", err), false)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, false, NewRequestError(http.StatusInternalServerError, fmt.Errorf("error executing request: %w", err), false)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		if res.StatusCode == 429 {
			return nil, false, NewRequestError(res.StatusCode, fmt.Errorf("rate limit reached"), true)
		}
		return nil, false, NewRequestError(res.StatusCode, fmt.Errorf("get request error"), false)
	}

	var response ResponseData[Res]
	err = response.DecodeBody(res.Body, fieldName)
	if err != nil {
		return nil, false, NewRequestError(http.StatusInternalServerError, fmt.Errorf("unable to decode response: %w", err), false)
	}

	items, more := response.ExtractItems()
	return items, more, NewRequestError(200, nil, false)
}
