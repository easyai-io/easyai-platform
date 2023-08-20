package config

import (
	"fmt"
	"strings"
	"sync"
)

var (
	// c global config
	syncFuncMap   = make(map[string]func(config *Config) error)
	syncFuncLock  sync.Mutex
	syncChannel   = make(chan struct{}, 64)
	_platformHost = "easyai.io"
)

// RegisterConfigSyncFunc sync handler for k8s config change
func RegisterConfigSyncFunc(key string, f func(config *Config) error) {
	syncFuncLock.Lock()
	defer syncFuncLock.Unlock()
	if _, ok := syncFuncMap[key]; !ok {
		syncFuncMap[key] = f
	}
}

// Cluster  configs for a k8s cluster
type Cluster struct {
	Disabled         bool                 `toml:"disabled"`
	Name             string               `toml:"name"`
	Config           string               `toml:"config"`
	Parallelism      int                  `toml:"parallelism"`
	Socks5Addr       string               `toml:"socks5_addr"`
	StorageClass     StorageClass         `toml:"storage_class"`
	NginxIngress     NginxIngress         `toml:"nginx_ingress"`
	ImagePullSecrets []ImagePullSecret    `toml:"image_pull_secrets"`
	ImageRepoReplace []ImageReplaceRule   `toml:"image_repo_replace"`
	TrainingNS       map[string]Namespace `toml:"training_namespaces"`
	ServingNS        map[string]Namespace `toml:"serving_namespaces"`
	ValidTaints      []string             `toml:"valid_taints"`
	GrafanaDashboard GrafanaDashboard     `toml:"grafana_dashboard"`
}

// Namespace config
type Namespace struct {
	Name            string   `toml:"name"`
	Disabled        bool     `toml:"disabled"`
	ShareNamespaces []string `toml:"share_namespaces"`
	Type            string   `toml:"type"`
	Comment         string   `toml:"comment"`
	Exclusive       bool     `toml:"exclusive"`
	_uid            string
	_clusterUID     string
}

// StorageClass k8s sc
type StorageClass struct {
	Data             string `toml:"data"`
	Notebook         string `toml:"notebook"`
	Job              string `toml:"job"`
	DataFBPrefix     string `toml:"data_fb_prefix"`
	NotebookFBPrefix string `toml:"notebook_fb_prefix"`
	JobFBPrefix      string `toml:"job_fb_prefix"`
}

// NginxIngress nginx ingress config
type NginxIngress struct {
	PlatformHost      string `toml:"platform_host"`
	PlatformProxyHost string `toml:"platform_proxy_host"`
	IngressClass      string `toml:"ingress_class"`
	HTTPS             bool   `toml:"https"`
}

// ImagePullSecret for hub secret
type ImagePullSecret struct {
	Name     string `toml:"name"`
	Registry string `toml:"registry"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// ImageReplaceRule for multi hub with auto sync
type ImageReplaceRule struct {
	Source string `toml:"source"`
	Target string `toml:"target"`
	Match  string `toml:"match"` // prefix | contain
}

// GrafanaDashboard ...
type GrafanaDashboard struct {
	Job  GrafanaURL `toml:"job"`
	Task GrafanaURL `toml:"task"`
	Node GrafanaURL `toml:"node"`
	GPU  GrafanaURL `toml:"gpu"`
}

// GrafanaURL ...
type GrafanaURL struct {
	BaseURL   string `toml:"base_url"`
	QueryArgs string `toml:"query_args"`
}

// ArgMap ...
func (gu *GrafanaURL) ArgMap() map[string]string {
	res := map[string]string{}
	for _, v := range strings.Split(gu.QueryArgs, ";") {
		kv := strings.SplitN(v, "=", 2)
		if len(kv) == 2 {
			res[kv[0]] = kv[1]
		}
	}
	return res
}

// GetStorageClassData fot training data
func (cluster *Cluster) GetStorageClassData() *string {
	sc := "ml-nfs-storage-data"
	if cluster.StorageClass.Data != "" {
		sc = cluster.StorageClass.Data
	}
	return &sc
}

// GetStorageClassJob for job output
func (cluster *Cluster) GetStorageClassJob() *string {
	sc := "ml-nfs-storage-job"
	if cluster.StorageClass.Job != "" {
		sc = cluster.StorageClass.Job
	}
	return &sc
}

// GetStorageClassNotebook for jupyter
func (cluster *Cluster) GetStorageClassNotebook() *string {
	sc := "ml-nfs-storage-notebook"
	if cluster.StorageClass.Notebook != "" {
		sc = cluster.StorageClass.Notebook
	}
	return &sc
}

// GetPlatformHost for nginx ingress
func (cluster *Cluster) GetPlatformHost() string {
	host := _platformHost
	if cluster.NginxIngress.PlatformHost != "" {
		host = cluster.NginxIngress.PlatformHost
		if vv := strings.Split(host, "://"); len(vv) == 2 {
			host = vv[1]
		}
	}
	return host
}

// GetGatewayURL for access
func (cluster *Cluster) GetGatewayURL() string {
	host := _platformHost
	if cluster.NginxIngress.PlatformHost != "" {
		host = cluster.NginxIngress.PlatformHost
		if vv := strings.Split(host, "://"); len(vv) == 2 {
			host = vv[1]
		}
	}
	if cluster.NginxIngress.HTTPS {
		return fmt.Sprintf("https://%s", host)
	}
	return fmt.Sprintf("http://%s", host)
}

// GetPlatformProxyHost for nginx ingress
func (cluster *Cluster) GetPlatformProxyHost() string {
	host := _platformHost
	if cluster.NginxIngress.PlatformProxyHost != "" {
		host = cluster.NginxIngress.PlatformProxyHost
		if vv := strings.Split(host, "://"); len(vv) == 2 {
			host = vv[1]
		}
	}
	return host
}

// GetIngressClass for nginx ingress
func (cluster *Cluster) GetIngressClass() string {
	ic := "nginx"
	if cluster.NginxIngress.IngressClass != "" {
		ic = cluster.NginxIngress.IngressClass
	}
	return ic
}

// GetPlatformProxyURL return easyai proxy url for a k8s cluster
func (cluster *Cluster) GetPlatformProxyURL() string {
	p := "http://"
	if cluster.NginxIngress.HTTPS {
		p = "https://"
	}
	return p + cluster.GetPlatformProxyHost()
}

// ValidTaint  check a taint of node, if taint key in white list
func (cluster *Cluster) ValidTaint(key string) bool {
	ts := append([]string{"gpu-pod", "allow-namespace", "easyai-node"}, cluster.ValidTaints...)
	for _, v := range ts {
		if v == key || strings.HasPrefix(key, "easyai-") {
			return true
		}
	}
	return false
}

// UID for cluster ns
func (ns *Namespace) UID(clusterUID string) string {
	if ns._uid == "" {
		ns._clusterUID = clusterUID
		ns._uid = fmt.Sprintf("%s-%s-%s", clusterUID, ns.Type, ns.Name)
	}
	return ns._uid
}

// Valid if a namespace valid in config define
func (ns *Namespace) Valid(cls Cluster) error {
	for _, name := range ns.ShareNamespaces {
		_, ok1 := cls.TrainingNS[name]
		_, ok2 := cls.ServingNS[name]
		if !ok1 && !ok2 {
			return fmt.Errorf("[config.go] namespace<%s> want share namespace<%s>, which is not exist in cluster<%s>", ns.Name, name, cls.Name)
		}
	}
	ns._clusterUID = cls.Name
	ns._uid = fmt.Sprintf("%s-%s-%s", cls.Name, ns.Type, ns.Name)
	return nil
}

// ClusterUID for a cluster
func (ns *Namespace) ClusterUID() string {
	return ns._clusterUID
}
