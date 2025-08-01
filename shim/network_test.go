package shim

import (
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
