package datasource

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Excalibur-1/configuration"
	"github.com/Excalibur-1/gutil"
	"github.com/rs/zerolog/log"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

const (
	app   = "base"
	group = "datasource"
	tag   = ""
	path  = "common"
)

type Config struct {
	Dialect      string                 `json:"dialect"`
	Dsn          string                 `json:"dsn"`
	Debug        bool                   `json:"debug"`
	EnableLog    bool                   `json:"enableLog"`
	Prefix       string                 `json:"prefix"`       // 表名前缀
	MinPoolSize  int                    `json:"minPoolSize"`  // pool最大空闲数
	MaxPoolSize  int                    `json:"maxPoolSize"`  // pool最大连接数
	IdleTimeout  gutil.Duration         `json:"idleTimeout"`  // 连接最长存活时间
	QueryTimeout gutil.Duration         `json:"queryTimeout"` // 查询超时时间
	ExecTimeout  gutil.Duration         `json:"execTimeout"`  // 执行超时时间
	TranTimeout  gutil.Duration         `json:"tranTimeout"`  // 事务超时时间
	Expand       map[string]interface{} `json:"expand"`
}

type dataSource struct {
	namespace  string
	systemId   string
	cfg        configuration.Configuration
	privileges map[string][]string
}

type DataSource interface {
	Config(dsID string) *Config
	Orm(dsID string) *xorm.Engine
}

// Engine 获取数据库引擎的唯一实例。
func Engine(cfg configuration.Configuration, namespace, systemId string) DataSource {
	fmt.Println("Loading Datasource Engine ver:1.0.0")
	return &dataSource{cfg: cfg, namespace: namespace, systemId: systemId, privileges: make(map[string][]string, 0)}
}

func (d *dataSource) Config(dsID string) *Config {
	ds, dsID, err := d.getConfiguration(dsID)
	if err != nil || ds.Dsn == "" {
		panic(fmt.Sprintf("数据源[%s]配置未指定或者读取时发生错误:%+v\n", dsID, err))
	}
	return ds
}

func (d *dataSource) Orm(dsID string) *xorm.Engine {
	c := d.Config(dsID)
	eng, err := xorm.NewEngine(c.Dialect, c.Dsn)
	if err == nil {
		// 这会在控制台打印出生成的SQL语句
		eng.ShowSQL(c.Debug)
		if c.EnableLog {
			// 这会在控制台打印调试及以上的信息
			eng.Logger().SetLevel(xlog.LOG_DEBUG)
		}
		// 设置表名前缀
		eng.SetTableMapper(names.NewPrefixMapper(names.SnakeMapper{}, c.Prefix))
		// 设置连接池的空闲数大小
		eng.SetMaxIdleConns(c.MinPoolSize)
		// 设置最大打开连接数
		eng.SetMaxOpenConns(c.MaxPoolSize)
		// 设置连接的最大生存时间
		eng.SetConnMaxLifetime(time.Duration(c.IdleTimeout))
	} else {
		panic(fmt.Sprintf("初始化datasource引擎出错%+v\n", err))
	}
	return eng
}

// 获取指定标示的数据源的配置信息，返回的配置Config对象
// 需要特别说明：如果给定的数据源标示为空，则表明是要获取当前业务系统的默认数据源配置信息。
// dsID 数据源的标示，如果为空则表明是默认数据源
func (d *dataSource) getConfiguration(dsID string) (*Config, string, error) {
	config := &Config{}
	// 如果是获取默认的数据源，则使用当前系统的标示，否则鉴权
	if dsID == "" {
		dsID = d.systemId
	} else {
		// 数据库的访问权限鉴权
		plist := d.systemPrivileges()
		if len(plist) == 0 || gutil.ContainsString(plist, dsID) == -1 {
			return config, "", fmt.Errorf("系统[%s]无数据源[%s]的访问权限", d.systemId, dsID)
		}
	}
	err := d.readFromConfiguration(dsID, config)
	return config, dsID, err
}

// 加载数据库的访问权限鉴权
func (d *dataSource) systemPrivileges() []string {
	d.cfg.Get(d.namespace, app, group, tag, []string{"privileges"}, d)
	plist := d.privileges[d.systemId]
	fmt.Printf("系统[%s]的数据源权限:%s\n", d.systemId, strings.Join(plist, ","))
	return plist
}

func (d *dataSource) Changed(data map[string]string) {
	for _, v := range data {
		var vl map[string][]string
		if err := json.Unmarshal([]byte(v), &vl); err == nil {
			for k, _v := range vl {
				d.privileges[k] = _v
			}
		}
	}
}

func (d *dataSource) readFromConfiguration(dsID string, config *Config) (err error) {
	if err = d.readCommonProperties(config); err != nil {
		return
	}
	fmt.Printf("从配置中心读取数据源配置:/%s/%s/%s\n", app, group, dsID)
	if err = d.cfg.Clazz(d.namespace, app, group, tag, dsID, config); err != nil {
		log.Error().Err(err).Msgf("数据源[%s]的配置获取失败", dsID)
	}
	return
}

func (d *dataSource) readCommonProperties(config *Config) (err error) {
	fmt.Printf("从配置中心的读取通用数据源配置:/%s/%s/%s\n", app, group, path)
	vl, err := d.cfg.String(d.namespace, app, group, tag, path)
	if err != nil {
		log.Error().Err(err).Msg("配置中心的通用数据源配置获取失败:%v")
		return
	}
	if err = json.Unmarshal([]byte(vl), config); err != nil {
		log.Error().Err(err).Msg("解析数据源的通用配置失败")
	}
	return
}
