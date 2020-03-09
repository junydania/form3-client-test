package form3

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

var accountId = "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"

const (
	listResponse = `{
		"data": [
		  {
			"type": "accounts",
			"id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
			"version": 0,
			"attributes": {
			  "country": "GB",
			  "base_currency": "GBP"
			}
		  },
		  {
			"type": "accounts",
			"id": "ea6239c1-99e9-42b3-bca1-92f5c068da6b",
			"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
			"version": 0,
			"attributes": {
			  "country": "GB",
			  "base_currency": "GBP"
			}
		  }
		]
	  }`

	postResponse = `{
		"data": {
			"type": "accounts",
			"id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
			"version": 0,
			"attributes": {
			  "country": "GB",
			  "base_currency": "GBP",
			  "account_number": "41426819",
			  "bank_id": 400300,
			  "bank_id_code": "GBDSC"
			}
		  }
	}`

	fetchResponse = `{
		"data": {
		  "type": "accounts",
		  "id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
		  "organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		  "version": 0,
		  "attributes": {
			"country": "GB",
			"base_currency": "GBP",
			"account_number": "41426819",
			"bank_id": 400300,
			"bank_id_code": "GBDSC"
		  }
		}
	}`
)

// func Form3HTTPSTestServer(handler http.Handler) (*httptest.Server, error) {
// 	ts := httptest.NewUnstartedServer(handler)
// 	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
// 	if err != nil {
// 		return nil, err
// 	}
// 	ts.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
// 	ts.StartTLS()
// 	return ts, nil
// }

func Form3APIResponseStub() *httptest.Server {
	var resp string

	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			resp = postResponse
		case "GET":
			if strings.Contains(r.RequestURI, "/organisation/accounts/") {
				resp = fetchResponse
			} else {
				resp = listResponse
			}
		case "DELETE":
			resp = "OK"
		default:
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Write([]byte(resp))

	}))
}

func TestFetchAccounts(t *testing.T) {
	server := Form3APIResponseStub()
	url, _ := url.Parse(server.URL)
	defer server.Close()
	client := NewClient(httpClient, url, "certs/client.pem", "certs/client.key")
	account, err := client.FetchAccount(accountId)
	assert.NoError(t, err)
	assert.Equal(t, accountId, account.Data.OrganisationID)
}

func TestCreateAccount(t *testing.T) {
	server := Form3APIResponseStub()
	url, _ := url.Parse(server.URL)

	defer server.Close()

	client := NewClient(httpClient, url, "certs/client.pem", "certs/client.key")

	postRequest := CreateAccountRequest{
		Data: CreateAccount{
			Type:           "accounts",
			ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
			Attributes: CreateAccountAttributes{
				Country:      "GB",
				BaseCurrency: "GBP",
				BankID:       400300,
				BankIDCode:   "GBDSC",
				Bic:          "NWBKGB22",
			},
		},
	}

	account, _, err := client.CreateAccount(postRequest)
	assert.NoError(t, err)
	assert.Equal(t, 200, account.StatusCode)
}

func TestDeleteAccounts(t *testing.T) {

	server := Form3APIResponseStub()
	url, _ := url.Parse(server.URL)
	defer server.Close()

	client := NewClient(httpClient, url, "certs/client.pem", "certs/client.key")
	account, _, err := client.DeleteAccount(accountId, 0)
	assert.NoError(t, err)
	assert.Equal(t, 200, account.StatusCode)
}

func TestListAccounts(t *testing.T) {
	server := Form3APIResponseStub()
	url, err := url.Parse(server.URL)
	defer server.Close()
	client := NewClient(httpClient, url, "certs/client.pem", "certs/client.key")
	accounts, err := client.ListAccounts(1, 100)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(accounts.Data))
}
