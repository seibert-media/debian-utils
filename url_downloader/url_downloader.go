package url_downloader

import (
	"fmt"
	"io/ioutil"
	"net/http"

	http_requestbuilder "github.com/bborbe/http/requestbuilder"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)
type HTTPRequestBuilderProvider func(url string) http_requestbuilder.HttpRequestBuilder

type URLDownloader interface {
	DownloadURL(url string) (string, error)
}

type urlDownloader struct {
	httpRequestBuilderProvider HTTPRequestBuilderProvider
	executeRequest             ExecuteRequest
}

func New(executeRequest ExecuteRequest, httpRequestBuilderProvider HTTPRequestBuilderProvider) *urlDownloader {
	u := new(urlDownloader)
	u.httpRequestBuilderProvider = httpRequestBuilderProvider
	u.executeRequest = executeRequest
	return u
}

func (u *urlDownloader) DownloadURL(url string) (string, error) {
	rb := u.httpRequestBuilderProvider(url)
	req, err := rb.Build()
	if err != nil {
		return "", err
	}
	resp, err := u.executeRequest(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("get url failed: %s", url)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
