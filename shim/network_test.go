package shim

import (
    "io"
    "testing"
)

func TestGetLocalIPs(t *testing.T) {
    localIPs, err := GetLocalIPs()
    if err != nil {
        t.Log("Error getting local IPs:", err)
        return
    }

    t.Log("Local IPs:", localIPs)
}

func TestGetLocalIP(t *testing.T) {
    t.Log(GetLocalIP())
}

func TestNewDefaultHttpClientWithProxyURL(t *testing.T) {
    args := []string{
        "",
        "http://127.0.0.1:8118",
        "socks5h://127.0.0.1:10222",
        "socks5h://127.0.0.1:10223",
    }

    for _, address := range args {
        httpClient, err := NewDefaultHttpClientWithProxyURL(address)
        if err != nil {
            t.Error("Error creating http client:", err)
        }
        resp, err := httpClient.Post("http://ipinfo.io", "application/json", nil)
        if err != nil {
            t.Errorf("Error posting to http client: %v", err)
            return
        }
        defer resp.Body.Close()

        bytes, err := io.ReadAll(resp.Body)
        if err != nil {
            t.Errorf("Error reading response body: %v", err)
            return
        }

        t.Log(string(bytes))
    }

}
