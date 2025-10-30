package app

import (
	"sync"
)

// 能力插件注册中心
// 支持能力注册、动态发现、热卸载

type AbilityRegistry struct {
	mu      sync.RWMutex
	plugins map[string]interface{}
}

func NewAbilityRegistry() *AbilityRegistry {
	return &AbilityRegistry{plugins: make(map[string]interface{})}
}

func (a *AbilityRegistry) Register(name string, plugin interface{}) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.plugins[name] = plugin
}

func (a *AbilityRegistry) Get(name string) (interface{}, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	v, ok := a.plugins[name]
	return v, ok
}

func (a *AbilityRegistry) Remove(name string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.plugins, name)
}

func (a *AbilityRegistry) List() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	names := make([]string, 0, len(a.plugins))
	for k := range a.plugins {
		names = append(names, k)
	}
	return names
}
