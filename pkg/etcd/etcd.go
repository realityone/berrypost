package etcd

import (
	"context"
	"fmt"
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

func (etcd *etcd) Get(key string) (string, error) {
	if etcd.client == nil {
		return "", errors.New("etcd not supported")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := etcd.client.Get(ctx, key)
	cancel()
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) != 1 {
		errStr := "etcd key not found"
		return "", errors.New(errStr)
	}
	return string(resp.Kvs[0].Value), nil
}

func (etcd *etcd) GetWithPrefix(prefix string) ([]string, []string, error) {
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
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
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

func (etcd *etcd) Close() error {
	if err := etcd.client.Close(); err != nil {
		return err
	}
	return nil
}