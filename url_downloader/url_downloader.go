package url_downloader

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bborbe/http/requestbuilder"
)

type UrlDownloader interface {
	DownloadUrl(url string) (string, error)
}

type urlDownloader struct {
	httpRequestBuilderProvider HttpRequestBuilderProvider
	client                     *http.Client
}

type HttpRequestBuilderProvider func(url string) requestbuilder.HttpRequestBuilder

func New(client *http.Client, httpRequestBuilderProvider HttpRequestBuilderProvider) *urlDownloader {
	u := new(urlDownloader)
	u.httpRequestBuilderProvider = httpRequestBuilderProvider
	u.client = client
	return u
}

func (u *urlDownloader) DownloadUrl(url string) (string, error) {
	rb := u.httpRequestBuilderProvider(url)
	req, err := rb.GetRequest()
	if err != nil {
		return "", err
	}
	resp, err := u.client.Do(req)
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
