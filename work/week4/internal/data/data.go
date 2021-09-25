package data

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.mypaas.com.cn/dmp/gopkg/db/mysql"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/conf"
	"sync"
)

// 业务数据访问，包含 cache、db 等封装，实现了 biz 的 repo 接口。我们可能会把 data 与 dao 混淆在一起，
// data 偏重业务的含义，它所要做的是将领域对象重新拿出来，我们去掉了 DDD 的 infra层。
var ProviderSet = wire.NewSet(NewData, NewIndicator)

type Data struct {
	log *logrus.Logger
	// 配置文件
	c *conf.Bootstrap
	//map的访问不是一个线程安全的数据结构，需要加互斥锁来保护
	mutex *sync.Mutex
	// 使用映射存储数据库链接
	clients map[string]mysql.SQLDatabase
}

type CrontabInfo struct {
	Name     string `json:"name" db:"name"`
	Schedule string `json:"schedule" db:"schedule"`
	Expire   int    `json:"expire" db:"expire"`
	IsEnable bool   `json:"is_enable" db:"is_enable"`
}

func NewData(c *conf.Bootstrap, log *logrus.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.Infoln("closing data resources")
	}
	return &Data{
		c:       c,
		mutex:   &sync.Mutex{},
		log:     log,
		clients: map[string]mysql.SQLDatabase{},
	}, cleanup, nil
}

func (d *Data) AllCrontab(ctx context.Context) ([]*CrontabInfo, error) {
	db, err := d.conn(ctx)
	if err != nil {
		return nil, err
	}
	sess := db.GetSession()
	querySql := "select name,schedule,expire,is_enable from t_crontab"
	var crons []*CrontabInfo
	err = sess.SelectContext(ctx, &crons, querySql)
	if err != nil {
		return nil, err
	}
	for i, item := range crons {
		d.log.Debugf("crontab(%d): %v\n", i, item)
	}
	return crons, nil
}

func (d *Data) Crontab(ctx context.Context, name string) (*CrontabInfo, error) {
	db, err := d.conn(ctx)
	if err != nil {
		return nil, err
	}
	sess := db.GetSession()
	querySql := "select name,schedule,expire,is_enable from t_crontab where name=?"
	var cron CrontabInfo
	err = sess.SelectContext(ctx, &cron, querySql)
	if err != nil {
		return nil, err
	}
	return &cron, nil
}

func (d *Data) conn(ctx context.Context) (mysql.SQLDatabase, error) {
	dbName := d.c.Db.Database
	// 加互斥锁保护对map的访问
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// 如果数据库链接已经存在则直接返回
	if database, ok := d.clients[dbName]; ok {
		return database, nil
	}
	defaultCfg := mysql.DefaultConfiguration()
	configuration := mysql.Configuration{
		Driver: "mysql",
		ConnectionURL: fmt.Sprintf("%s:%s@(%s:%d)/%s?%s",
			d.c.Db.User,
			d.c.Db.Password,
			d.c.Db.Host,
			d.c.Db.Port,
			dbName,
			mysql.MysqlConnParams),

		MaxIdleConns: defaultCfg.MaxIdleConns,
		MaxOpenConns: defaultCfg.MaxOpenConns,
	}
	db, err := mysql.NewSQLDatabase(configuration)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open project database: %s", dbName)
	}
	//缓存数据库
	d.clients[dbName] = db
	return db, err
}
