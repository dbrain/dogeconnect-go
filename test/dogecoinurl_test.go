package test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	dogeconnectgo "github.com/dogeorg/dogeconnect-go"
)

func TestDogecoinURL(t *testing.T) {
	payTo := "DPD7uK4B1kRmbfGmytBhG1DZjaMWNfbpwY"
	amount := "12.25"
	connectURL := "example.com/dc/1QAB-POvTh2R88nybE8Wwg"
	pubKey, _ := hex.DecodeString("6c52b17752f469c5411b977ba64725d40174d16e780b709b2aff68e0f5abfc50")
	expect := "dogecoin:DPD7uK4B1kRmbfGmytBhG1DZjaMWNfbpwY?amount=12.25&dc=example.com%2Fdc%2F1QAB-POvTh2R88nybE8Wwg&h=72b-LVh5K_mm7zyN9PXO"
	uri, err := dogeconnectgo.DogecoinURI(payTo, amount, "https://"+connectURL, pubKey)
	if err != nil {
		t.Fatalf("failed to build uri: %v", err)
	}
	if uri != expect {
		t.Errorf("incorrect uri:\n%v (found)\n%v (expected)", uri, expect)
	}
	res, err := dogeconnectgo.ParseDogecoinURI(uri)
	if err != nil {
		t.Errorf("failed to parse uri: %v", err)
	}
	if !res.IsConnectURI() {
		t.Errorf("IsConnectURI should return true")
	}
	if res.Address != payTo {
		t.Errorf("wrong address: %v vs %v", res.Address, payTo)
	}
	if res.Amount != amount {
		t.Errorf("wrong amount: %v vs %v", res.Amount, amount)
	}
	if res.ConnectURL != connectURL {
		t.Errorf("wrong connect URL: %v vs %v", res.ConnectURL, connectURL)
	}
	pubSha := sha256.Sum256(pubKey)
	if !bytes.Equal(res.PubKeyHash, pubSha[0:15]) {
		t.Errorf("wrong pubkey hash:\n%x vs\n%x", res.PubKeyHash, pubSha[0:15])
	}
}

func TestPlainDogecoinURI(t *testing.T) {
	res, err := dogeconnectgo.ParseDogecoinURI("dogecoin:DPD7uK4B1kRmbfGmytBhG1DZjaMWNfbpwY?amount=8.25")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsConnectURI() {
		t.Error("plain dogecoin URI should not be a connect URI")
	}
	if res.Address != "DPD7uK4B1kRmbfGmytBhG1DZjaMWNfbpwY" {
		t.Errorf("wrong address: %v", res.Address)
	}
	if res.Amount != "8.25" {
		t.Errorf("wrong amount: %v", res.Amount)
	}
}

func TestParseDogecoinURIErrors(t *testing.T) {
	tests := []struct {
		name string
		uri  string
	}{
		{"wrong scheme", "bitcoin:DPD7uK4B1kRmbfGmytBhG1DZjaMWNfbpwY?amount=1"},
		{"bad base64 h", "dogecoin:DPD7uK4B1kRmbfGmytBhG1DZjaMWNfbpwY?amount=1&dc=example.com&h=!!!invalid!!!"},
		{"dc without h", "dogecoin:DPD7uK4B1kRmbfGmytBhG1DZjaMWNfbpwY?amount=1&dc=example.com/dc/1234"},
		{"h without dc", "dogecoin:DPD7uK4B1kRmbfGmytBhG1DZjaMWNfbpwY?amount=1&h=72b-LVh5K_mm7zyN9PXO"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := dogeconnectgo.ParseDogecoinURI(tc.uri)
			if err == nil {
				t.Errorf("expected error for %q, got nil", tc.uri)
			}
		})
	}
}

func TestDogeConnectURI(t *testing.T) {
	connectURL := "example.com/dc/xyz123"
	pubKey, _ := hex.DecodeString("6c52b17752f469c5411b977ba64725d40174d16e780b709b2aff68e0f5abfc50")
	expect := "dogeconnect:example.com/dc/xyz123?h=72b-LVh5K_mm7zyN9PXO"
	uri, err := dogeconnectgo.DogeConnectURI("https://"+connectURL, pubKey)
	if err != nil {
		t.Fatalf("failed to build uri: %v", err)
	}
	if uri != expect {
		t.Errorf("incorrect uri:\n%v (found)\n%v (expected)", uri, expect)
	}
	res, err := dogeconnectgo.ParseDogeConnectURI(uri)
	if err != nil {
		t.Fatalf("failed to parse uri: %v", err)
	}
	if !res.IsConnectURI() {
		t.Errorf("IsConnectURI should return true")
	}
	if res.Address != "" {
		t.Errorf("address should be empty, got %q", res.Address)
	}
	if res.Amount != "" {
		t.Errorf("amount should be empty, got %q", res.Amount)
	}
	if res.ConnectURL != connectURL {
		t.Errorf("wrong connect URL: %v vs %v", res.ConnectURL, connectURL)
	}
	pubSha := sha256.Sum256(pubKey)
	if !bytes.Equal(res.PubKeyHash, pubSha[0:15]) {
		t.Errorf("wrong pubkey hash:\n%x vs\n%x", res.PubKeyHash, pubSha[0:15])
	}
}

func TestDogeConnectURIStripsHTTPS(t *testing.T) {
	pubKey, _ := hex.DecodeString("6c52b17752f469c5411b977ba64725d40174d16e780b709b2aff68e0f5abfc50")
	uri, err := dogeconnectgo.DogeConnectURI("https://example.com/dc/1", pubKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	res, err := dogeconnectgo.ParseDogeConnectURI(uri)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if res.ConnectURL != "example.com/dc/1" {
		t.Errorf("expected stripped URL, got %q", res.ConnectURL)
	}
}

func TestParseDogeConnectURIErrors(t *testing.T) {
	tests := []struct {
		name string
		uri  string
	}{
		{"wrong scheme", "dogecoin:DPD7uK4B1kRmbfGmytBhG1DZjaMWNfbpwY?amount=1"},
		{"missing connect url", "dogeconnect:?h=72b-LVh5K_mm7zyN9PXO"},
		{"missing h", "dogeconnect:example.com/dc/xyz123"},
		{"bad base64 h", "dogeconnect:example.com/dc/xyz123?h=!!!invalid!!!"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := dogeconnectgo.ParseDogeConnectURI(tc.uri)
			if err == nil {
				t.Errorf("expected error for %q, got nil", tc.uri)
			}
		})
	}
}

func TestDogeConnectURIWithQueryParams(t *testing.T) {
	// connectURL already has a query string — must not produce a double-? URI
	pubKey, _ := hex.DecodeString("6c52b17752f469c5411b977ba64725d40174d16e780b709b2aff68e0f5abfc50")
	uri, err := dogeconnectgo.DogeConnectURI("https://example.com/dc/1?foo=bar", pubKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// must be parseable (validates the URI is well-formed)
	if _, err := dogeconnectgo.ParseDogeConnectURI(uri); err != nil {
		t.Errorf("produced unparseable URI %q: %v", uri, err)
	}
	// must contain exactly one '?'
	count := 0
	for _, c := range uri {
		if c == '?' {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly one '?' in URI, got %d: %q", count, uri)
	}
}

func TestDogeConnectURINotHTTPS(t *testing.T) {
	pubKey, _ := hex.DecodeString("6c52b17752f469c5411b977ba64725d40174d16e780b709b2aff68e0f5abfc50")
	_, err := dogeconnectgo.DogeConnectURI("http://example.com/dc/1", pubKey)
	if err == nil {
		t.Error("expected error for non-https URL, got nil")
	}
}

func TestDogeConnectURIBadPubKey(t *testing.T) {
	_, err := dogeconnectgo.DogeConnectURI("https://example.com/dc/1", []byte{1, 2, 3})
	if err == nil {
		t.Error("expected error for short pubkey, got nil")
	}
}

func TestDogecoinURIBadPubKey(t *testing.T) {
	_, err := dogeconnectgo.DogecoinURI("DPD7uK4B1kRmbfGmytBhG1DZjaMWNfbpwY", "1.0", "https://example.com/dc/1", []byte{1, 2, 3})
	if err == nil {
		t.Error("expected error for short pubkey, got nil")
	}
}

func TestSlashInDC(t *testing.T) {
	connectURL := "example.com/dc/1QAB"
	uri := "dogecoin:DPD7uK4B1kRmbfGmytBhG1DZjaMWNfbpwY?amount=12.25&dc=example.com/dc/1QAB&h=72b-LVh5K_mm7zyN9PXO"
	res, err := dogeconnectgo.ParseDogecoinURI(uri)
	if err != nil {
		t.Errorf("failed to parse uri: %v", err)
	}
	if !res.IsConnectURI() {
		t.Errorf("IsConnectURI should return true")
	}
	if res.ConnectURL != connectURL {
		t.Errorf("wrong connect URL: %v vs %v", res.ConnectURL, connectURL)
	}
}
