package shim

import (
    "fmt"
    "net"
    "net/http"
    "net/url"

    "github.com/pkg/errors"
)

// GetLocalIPs 获取所有非环回、非链路本地的IPv4地址
func GetLocalIPs() ([]net.IP, error) {
    var ips []net.IP

    // 获取所有网络接口
    interfaces, err := net.Interfaces()
    if err != nil {
        return nil, fmt.Errorf("failed to get network interfaces: %w", err)
    }

    for _, iface := range interfaces {
        // 过滤掉环回接口和关闭的接口
        if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
            continue
        }

        // 获取接口的所有地址
        addrs, err := iface.Addrs()
        if err != nil {
            fmt.Printf("Warning: failed to get addresses for interface %s: %v\n", iface.Name, err)
            continue
        }

        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
            case *net.IPNet:
                ip = v.IP
            case *net.IPAddr:
                ip = v.IP
            }

            // 过滤掉环回地址
            if ip == nil || ip.IsLoopback() {
                continue
            }

            // 确保是IPv4地址
            ip = ip.To4()
            if ip == nil { // 不是IPv4地址
                continue
            }

            // 过滤掉链路本地地址 (169.254.x.x)
            // 链路本地地址通常用于无DHCP环境下的自动配置，不是我们通常意义上的局域网IP
            if ip.IsLinkLocalUnicast() {
                continue
            }

            // 过滤掉多播地址
            if ip.IsMulticast() {
                continue
            }

            // 过滤掉未指定的地址 (0.0.0.0)
            if ip.IsUnspecified() {
                continue
            }

            ips = append(ips, ip)
        }
    }
    return ips, nil
}

// GetLocalIP 获取本地IP
func GetLocalIP() string {
    ips, err := GetLocalIPs()
    if err != nil {
        return ""
    }
    return ips[0].String()
}

// NewDefaultHttpClientWithProxyURL 基于Proxy Address初始化一个走Proxy的http.Client实例，支持HTTP Proxy和Socks Proxy
func NewDefaultHttpClientWithProxyURL(address string) (*http.Client, error) {
    if address == "" {
        return http.DefaultClient, nil
    }

    // 有配置代理
    proxyURL, err := url.Parse(address)
    if err != nil {
        return nil, errors.Wrapf(err, "invalid proxy url: %s", address)
    }

    httpClient := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxyURL),
        },
        CheckRedirect: nil,
        Jar:           nil,
        Timeout:       0,
    }
    return httpClient, nil
}
