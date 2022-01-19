package etcd

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
)

type etcd struct {
	client *clientv3.Client
}

var Dao *etcd

func init() {
	client, err := NewClient()
	if err != nil {
		log.Println(err)
	}
	Dao = &etcd{
		client: client,
	}
}

func NewClient() (*clientv3.Client, error) {
	port := os.Getenv("ETCD")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{port},
		DialTimeout: time.Second,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (etcd *etcd) Put(key string, value string) error {
	if etcd.client == nil {
		return errors.New("etcd not supported")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := etcd.client.Put(ctx, key, value)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (etcd *etcd) GetWithPrefix(key string) (map[string]string, error) {
	rst := make(map[string]string)
	if etcd.client == nil {
		return nil, errors.New("etcd not supported")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := etcd.client.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, err
	}
	for _, ev := range resp.Kvs {
		k := string(ev.Key)
		v := string(ev.Value)
		rst[k] = v
	}
	return rst, nil
}

func (etcd *etcd) GetIfExist(key string) (string, error) {
	if etcd.client == nil {
		return "", errors.New("etcd not supported")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := etcd.client.Get(ctx, key)
	cancel()
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", errors.New("key does not exist")
	}
	return string(resp.Kvs[0].Value), nil
}

func (etcd *etcd) Exist(key string) (bool, error) {
	if etcd.client == nil {
		return false, errors.New("etcd not supported")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := etcd.client.Get(ctx, key)
	cancel()
	if err != nil {
		return false, err
	}
	if len(resp.Kvs) == 0 {
		return false, nil
	}
	return true, nil
}

func (etcd *etcd) GetKVWithPrefix(prefix string) ([]string, []string, error) {
	var (
		keys   []string
		values []string
	)
	if etcd.client == nil {
		return nil, nil, errors.New("etcd not supported")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := etcd.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	cancel()
	if err != nil {
		return nil, nil, err
	}
	for _, ev := range resp.Kvs {
		keys = append(keys, string(ev.Key))
		values = append(values, string(ev.Value))
	}
	return keys, values, nil
}

func (etcd *etcd) Delete(key string) error {
	if etcd.client == nil {
		return errors.New("etcd not supported")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := etcd.client.Delete(ctx, key)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (etcd *etcd) DeleteWithPrefix(key string) error {
	if etcd.client == nil {
		return errors.New("etcd not supported")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := etcd.client.Delete(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (etcd *etcd) Update(key string, value string) error {
	if etcd.client == nil {
		return errors.New("etcd not supported")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := etcd.client.Delete(ctx, key)
	cancel()
	if err != nil {
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	_, err = etcd.client.Put(ctx, key, value)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (etcd *etcd) Close() error {
	if err := etcd.client.Close(); err != nil {
		return err
	}
	return nil
}
