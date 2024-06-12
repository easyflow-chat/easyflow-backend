package database

import (
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseInst struct {
	client   *gorm.DB
	refCount int
	lock     sync.Mutex
}

func NewDatabaseInst(url string, config *gorm.Config) (*DatabaseInst, error) {
	db, err := gorm.Open(mysql.Open(url), config)
	if err != nil {
		return nil, err
	}

	return &DatabaseInst{
		client:   db,
		refCount: 1,
	}, nil
}

func (d *DatabaseInst) Acquire() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.refCount++
}

func (d *DatabaseInst) Release() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.refCount--
	if d.refCount == 0 {
		d.client = nil
	}
}

func (d *DatabaseInst) GetClient() *gorm.DB {
	return d.client
}

func (d *DatabaseInst) Migrate() error {
	return d.client.AutoMigrate(&Message{}, &Chat{}, &User{}, &ChatUserKeys{})
}
