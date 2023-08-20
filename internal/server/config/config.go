package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"

	"github.com/easyai-io/easyai-platform/pkg/util/json"
	"github.com/easyai-io/easyai-platform/pkg/util/structure"
)

var (
	// c global config
	c    Config
	once sync.Once
	lock sync.RWMutex
)

// Global return config for global usage
func Global() *Config {
	lock.RLock()
	defer lock.RUnlock()
	return c.DeepCopy()
}

// SetWWW set www dir
func SetWWW(www string) {
	lock.Lock()
	defer lock.Unlock()
	c.WWW = www
}

// SetGracePeriodSeconds set grace period seconds
func SetGracePeriodSeconds(gracePeriodSeconds int) {
	lock.Lock()
	defer lock.Unlock()
	c.GracePeriodSeconds = gracePeriodSeconds
}

// MustLoad Load config file (toml/json/yaml)
func MustLoad(fpath string) {
	once.Do(func() {

		reloadConfig(fpath)

		// watch && update
		go func() {
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				fmt.Println("should not happen", err)
				panic(err)
			}
			defer watcher.Close()
			if err = watcher.Add(fpath); err != nil {
				fmt.Println("should not happen", err)
				panic(err)
			}

			for {
				select {
				case event := <-watcher.Events:
					fmt.Println("receive config file event:", event)
					// k8s configmap file use symlinks, for example：link-file(mounted in container) --> cm-file
					// when cm is changed, the kubelet will create new-file, modify link-file to new-file，and delete old-file。
					// since fsnotify will watch the real file refer by link-file, so the watch event is REMOVE with real-file when cm changed,
					// we should re-watch the link file

					// warn: k8s sub-path mount use bind-mount(which will parse symlinks and use real file),
					// so with sub-path mount, file in container will never be updated; we should mount cm as directory, do not use sub-path mount
					if event.Op == fsnotify.Remove {
						// file deleted or k8s configMap changed
						_ = watcher.Remove(event.Name)
						time.Sleep(time.Second)
						// re-watch config files(symlink/file)
						if err = watcher.Add(fpath); err != nil {
							fmt.Println("should not happen", err)
							panic(err)
						}
						reloadConfig(fpath)
					}
					// also allow normal files to be modified and reloaded.
					if event.Op&fsnotify.Write == fsnotify.Write {
						time.Sleep(time.Second)
						reloadConfig(fpath)
					}
				case err := <-watcher.Errors:
					fmt.Println("watch config error:", err)
				}
			}
		}()

		// exec syncFuncMap when Config Changed
		go func() {
			for {
				_, ok := <-syncChannel
				if !ok {
					break
				}
				time.Sleep(time.Second / 2)
				func() {
					syncFuncLock.Lock()
					defer syncFuncLock.Unlock()
					for key, f := range syncFuncMap {
						if err := f(Global()); err != nil {
							fmt.Println("exec sync config func", key, err)
						}
					}
				}()
			}

		}()

	})
}

func reloadConfig(fpath string) {
	func() {
		lock.Lock()
		defer lock.Unlock()
		_, err := toml.DecodeFile(fpath, &c)
		if err != nil {
			panic(err)
		}
		setDefault(&c)
	}()
	PrintWithJSON()
	syncChannel <- struct{}{}
}

func setDefault(c *Config) {

}

// PrintWithJSON print config
func PrintWithJSON() {
	if cc := Global(); cc.PrintConfig {
		b, err := json.MarshalIndent(cc, "", " ")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
			return
		}
		os.Stdout.WriteString(string(b) + "\n")
	}
}

// IsDebugMode is debug
func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

// Config config
type Config struct {
	Environment        string             `toml:"environment"`
	AppName            string             `toml:"app_name"`
	AppGroup           string             `toml:"app_group"`
	RunMode            string             `toml:"run_mode"`
	WWW                string             `toml:"www"`
	BuildInWWWAsset    bool               `toml:"build_in_www_asset"`
	GracePeriodSeconds int                `toml:"grace_period_seconds"`
	Docs               string             `toml:"docs"`
	Swagger            bool               `toml:"swagger"`
	PrintConfig        bool               `toml:"print_config"`
	HTTP               HTTP               `toml:"http"`
	Log                Log                `toml:"log"`
	JWTAuth            JWTAuth            `toml:"jwt_auth"`
	Monitor            Monitor            `toml:"monitor"`
	RateLimiter        RateLimiter        `toml:"rate_limiter"`
	CORS               CORS               `toml:"cors"`
	GZIP               GZIP               `toml:"gzip"`
	Redis              Redis              `toml:"redis"`
	EntOrm             EntOrm             `toml:"ent_orm"`
	MySQL              MySQL              `toml:"mysql"`
	Sqlite3            Sqlite3            `toml:"sqlite3"`
	Postgres           Postgres           `toml:"postgres"`
	Clusters           map[string]Cluster `toml:"clusters"`
	MessageChannel     MessageChannel     `toml:"message_channel"`
}

// Log log
type Log struct {
	Level         int    `toml:"level"`
	Format        string `toml:"format"`
	Output        string `toml:"output"`
	OutputFile    string `toml:"output_file"`
	RotationCount int    `toml:"rotation_count"`
	RotationTime  int    `toml:"rotation_time"`
}

// JWTAuth jwt auth
type JWTAuth struct {
	Enable        bool   `toml:"enable"`
	SigningMethod string `toml:"signing_method"`
	SigningKey    string `toml:"signing_key"`
	Expired       int    `toml:"expired"`
	Store         string `toml:"store"`
	FilePath      string `toml:"file_path"`
	RedisDB       int    `toml:"redis_db"`
	RedisPrefix   string `toml:"redis_prefix"`
}

// HTTP struct
type HTTP struct {
	Host               string `toml:"host"`
	Port               int    `toml:"port"`
	CertFile           string `toml:"cert_file"`
	KeyFile            string `toml:"key_file"`
	ShutdownTimeout    int    `toml:"shutdown_timeout"`
	MaxContentLength   int64  `toml:"max_content_length"`
	MaxReqLoggerLength int    `toml:"max_req_logger_length" default:"1024"`
	MaxResLoggerLength int    `toml:"max_res_logger_length" `
}

// Monitor monitor
type Monitor struct {
	Enable    bool   `toml:"enable"`
	Addr      string `toml:"addr"`
	ConfigDir string `toml:"config_dir"`
}

// RateLimiter rate limiter
type RateLimiter struct {
	Enable  bool  `toml:"enable"`
	Count   int64 `toml:"count"`
	RedisDB int   `toml:"redis_db"`
}

// CORS cors
type CORS struct {
	Enable           bool     `toml:"enable"`
	AllowOrigins     []string `toml:"allow_origins"`
	AllowMethods     []string `toml:"allow_methods"`
	AllowHeaders     []string `toml:"allow_headers"`
	AllowCredentials bool     `toml:"allow_credentials"`
	MaxAge           int      `toml:"max_age"`
}

// GZIP gzip
type GZIP struct {
	Enable             bool     `toml:"enable"`
	ExcludedExtensions []string `toml:"excluded_extensions"`
	ExcludedPaths      []string `toml:"excluded_paths"`
}

// Redis  struct redis
type Redis struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

// EntOrm for sql
type EntOrm struct {
	Debug             bool   `toml:"debug"`
	DBType            string `toml:"db_type"`
	MaxLifetime       int    `toml:"max_life_time"`
	MaxOpenConns      int    `toml:"max_open_conns"`
	MaxIdleConns      int    `toml:"max_idle_conns"`
	EnableAutoMigrate bool   `toml:"enable_auto_migrate"`
}

// MySQL mysql
type MySQL struct {
	ReadDSN  string `toml:"read_dsn"`
	WriteDSN string `toml:"write_dsn"`
}

// Sqlite3 Sqlite3
type Sqlite3 struct {
	InMemoryMode bool   `toml:"in_memory_mode"`
	FilePath     string `toml:"file_path"`
}

// Postgres Postgres
type Postgres struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
	SSLMode  bool   `toml:"ssl_mode"`
}

// DSN for postgres
func (p *Postgres) DSN() string {
	sslMode := "sslmode=disable"
	if p.SSLMode {
		sslMode = "sslmode=enable"
	}
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s %s",
		p.Host, p.Port, p.User, p.DBName, p.Password, sslMode)
}

// MessageChannel to send msg
type MessageChannel struct {
	DefaultReceivers []string    `toml:"default_receivers"`
	WechatRobot      WechatRobot `toml:"wechat_robot"`
	WechatApp        WechatApp   `toml:"wechat_app"`
	Email            Email       `toml:"email"`
	StdPrinter       StdPrinter  `toml:"std_printer"`
}

// WechatRobot for 群bot
type WechatRobot struct {
	Webhook string `toml:"webhook"`
}

// WechatApp for 企微App
type WechatApp struct {
	CorpID     string `toml:"corp_id"`
	CorpSecret string `toml:"corp_secret"`
	AgentID    int64  `toml:"agent_id"`
	// for msg callback
	Token     string `toml:"token"`
	EncodeKey string `toml:"encode_key"`
}

// Email for mail
type Email struct {
	Suffix   string `toml:"suffix"`
	SendFrom string `toml:"send_from"`
	Password string `toml:"password"`
	Server   string `toml:"server"`
}

// StdPrinter for std
type StdPrinter struct {
	Stdout bool `toml:"stdout"`
	Stderr bool `toml:"stderr"`
}

// GetCluster return cluster config
func (c *Config) GetCluster(clusterName string) *Cluster {
	v, ok := c.Clusters[clusterName]
	if !ok {
		panic(fmt.Errorf("cluster(%s) not found in config file", clusterName))
	}
	if v.Parallelism <= 0 {
		v.Parallelism = 4
	}
	return &v
}

// CheckClusterNamespace 检查cluster和namespace是否可用
func (c *Config) CheckClusterNamespace(clusterName, nsName string, servingWorkload ...bool) error {
	v, ok := c.Clusters[clusterName]
	if !ok {
		return fmt.Errorf("cluster<%s> not exist", clusterName)
	}
	if v.Disabled {
		return fmt.Errorf("cluster<%s> is disabled by admin", clusterName)
	}
	ns, ok := v.TrainingNS[nsName]
	if len(servingWorkload) > 0 && servingWorkload[0] {
		ns, ok = v.ServingNS[nsName]
	}
	if !ok {
		return fmt.Errorf("namespace<%s> for cluster<%s> not exist", nsName, clusterName)
	}
	if ns.Disabled {
		return fmt.Errorf("namespace<%s> for cluster<%s> disabled", nsName, clusterName)
	}
	return nil
}

// DSN for sqlite3
func (s *Sqlite3) DSN() string {
	/*
		https://github.com/mattn/go-sqlite3/blob/master/README.md#faq
		Each connection to ":memory:" opens a brand new in-memory sql database,
		so if only specified ":memory:", that connection will see a brand new database.
		A workaround is to use "file::memory:?cache=shared" (or "file:foobar?mode=memory&cache=shared")
	*/
	dsn := "file:ent?mode=memory&cache=shared&_fk=1" // in memory mod
	if !s.InMemoryMode && s.FilePath != "" {
		_ = os.MkdirAll(filepath.Dir(s.FilePath), 0777)
		dsn = fmt.Sprintf("file:%s?_auth&_auth_user=root&_auth_pass=123456&cache=shared&_fk=1", s.FilePath)
	}
	return dsn
}

// DeepCopy for config
func (c *Config) DeepCopy() *Config {
	var config Config
	_ = structure.Copy(c, &config)
	return &config
}

// SetFakeConfig for test
func SetFakeConfig(cfg Config) {
	lock.Lock()
	defer lock.Unlock()
	_ = structure.Copy(&cfg, &c)
	setDefault(&c)
}
