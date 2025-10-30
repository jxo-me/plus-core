package config

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ConsulKVBridge 基于HTTP的轻量KV watcher + 配置提供器（JSON）
// 用法：
//   br := NewConsulKVBridge("http://127.0.0.1:8500", "config/app.json", "", 60*time.Second)
//   b.WithConfigProvider(br.Provider).Watchers(br.Watcher)

type ConsulKVBridge struct {
	addr   string
	key    string
	token  string
	wait   time.Duration
	client *http.Client
	// 最新配置缓存（JSON -> map[string]any）
	lastCfg map[string]any
}

func NewConsulKVBridge(addr, key, token string, wait time.Duration) *ConsulKVBridge {
	if wait <= 0 {
		wait = 60 * time.Second
	}
	return &ConsulKVBridge{
		addr:   strings.TrimRight(addr, "/"),
		key:    strings.TrimLeft(key, "/"),
		token:  token,
		wait:   wait,
		client: &http.Client{Timeout: wait + 10*time.Second},
	}
}

// Provider 供 Bootstrap.WithConfigProvider 使用
func (c *ConsulKVBridge) Provider(ctx context.Context) interface{} {
	return c.lastCfg
}

// Watcher 供 Bootstrap.Watchers 使用：监听KV变化并在变化时调用 onChange
func (c *ConsulKVBridge) Watcher(ctx context.Context, onChange func()) error {
	var lastIndex uint64
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		index, cfg, err := c.blockingGet(ctx, lastIndex)
		if err != nil {
			// 退避后重试
			time.Sleep(2 * time.Second)
			continue
		}
		if index > lastIndex {
			lastIndex = index
			c.lastCfg = cfg
			onChange()
		}
	}
}

func (c *ConsulKVBridge) blockingGet(ctx context.Context, lastIndex uint64) (uint64, map[string]any, error) {
	q := url.Values{}
	if lastIndex > 0 {
		q.Set("index", strconv.FormatUint(lastIndex, 10))
	}
	q.Set("wait", fmt.Sprintf("%ds", int(c.wait.Seconds())))
	endpoint := fmt.Sprintf("%s/v1/kv/%s?%s", c.addr, url.PathEscape(c.key), q.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, nil, err
	}
	if c.token != "" {
		req.Header.Set("X-Consul-Token", c.token)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return headerIndex(resp.Header.Get("X-Consul-Index")), nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return 0, nil, errors.New(string(b))
	}
	idx := headerIndex(resp.Header.Get("X-Consul-Index"))
	// Consul KV API返回数组，元素含Value通常为base64
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	// 兼容返回为数组或对象的情况
	var obj any
	if err := json.Unmarshal(b, &obj); err != nil {
		return idx, nil, nil
	}
	switch v := obj.(type) {
	case []any:
		if len(v) == 0 {
			return idx, nil, nil
		}
		m, _ := v[0].(map[string]any)
		if raw, ok := m["Value"].(string); ok && raw != "" {
			if cfg, ok := parseConsulValue(raw); ok {
				return idx, cfg, nil
			}
		}
		if val, ok := m["Value"].(map[string]any); ok {
			return idx, val, nil
		}
		return idx, nil, nil
	case map[string]any:
		return idx, map[string]any{"data": v}, nil
	default:
		return idx, nil, nil
	}
}

// parseConsulValue 同时支持base64编码或明文JSON两种形式
func parseConsulValue(raw string) (map[string]any, bool) {
	// 尝试作为base64解码
	if dec, err := base64.StdEncoding.DecodeString(raw); err == nil {
		var m map[string]any
		if json.Unmarshal(dec, &m) == nil {
			return m, true
		}
	}
	// 尝试作为JSON明文
	var mm map[string]any
	if json.Unmarshal([]byte(raw), &mm) == nil {
		return mm, true
	}
	return nil, false
}

func headerIndex(s string) uint64 {
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}
