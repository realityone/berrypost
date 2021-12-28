package management

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/realityone/berrypost/pkg/etcd"
	"github.com/sirupsen/logrus"
)

func (m Management) allUserBlueprints(ctx context.Context, userid string) []*ProtoFileMeta {
	//todo: 键设计
	var out []*ProtoFileMeta
	prefix := "/" + userid + "/"
	keys, _, err := etcd.Dao.GetWithPrefix(prefix)
	if err != nil {
		logrus.Error("Failed to get user blueprints: %+v", err)
		return nil
	}
	for _, key := range keys {
		blueprintName := strings.TrimPrefix(key, prefix)
		pm := &ProtoFileMeta{
			Filename: blueprintName,
		}
		pm.Meta.ImportPath = blueprintName
		out = append(out, pm)
	}
	return out
}

func (m Management) blueprintMethods(ctx context.Context, userid string, blueprintIdentifier string) ([]*BlueprintMethodInfo, error) {
	info := []*BlueprintMethodInfo{}
	key := "/" + userid + "/" + blueprintIdentifier
	value, err := etcd.Dao.Get(key)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(value), &info); err != nil {
		return nil, err
	}
	return info, nil
}

func (m Management) newBlueprintMethod(ctx context.Context, userid string, blueprintIdentifier string, newMethods []*BlueprintMethodInfo) ([]*BlueprintMethodInfo, error) {
	info := []*BlueprintMethodInfo{}
	key := "/" + userid + "/" + blueprintIdentifier
	value, err := etcd.Dao.Get(key)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(value), &info); err != nil {
		return nil, err
	}
	info = append(info, newMethods...)
	infoByte, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	if err := etcd.Dao.Put(key, string(infoByte)); err != nil {
		return nil, err
	}
	return info, nil
}

func (m Management) newBlueprint(ctx context.Context, userid string, blueprintIdentifier string, info []*BlueprintMethodInfo) ([]*BlueprintMethodInfo, error) {
	key := "/" + userid + "/" + blueprintIdentifier
	infoByte, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	if err := etcd.Dao.Put(key, string(infoByte)); err != nil {
		return nil, err
	}
	return info, nil
}
