package management

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strings"

	"github.com/realityone/berrypost/pkg/etcd"
	log "github.com/sirupsen/logrus"
)

func (m Management) newBlueprint(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	type BlueprintReq struct {
		BlueprintName string `json:"blueprintName"`
	}
	var reqInfo BlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		log.Error(err)
		return
	}
	key := m.fullKey(userid, reqInfo.BlueprintName)
	if err := m.putBlueprint(key, nil); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) copyBlueprintFromFile(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	type BlueprintReq struct {
		BlueprintName string `json:"blueprintName"`
		FileName      string `json:"fileName"`
	}
	var reqInfo BlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		log.Error(err)
		return
	}
	reqInfo.BlueprintName = m.trim(reqInfo.BlueprintName)
	reqInfo.FileName = m.trim(reqInfo.FileName)
	methods, err := m.getMethodsByService(ctx, reqInfo.FileName)
	if err != nil {
		log.Error(err)
		return
	}
	key := m.fullKey(userid, reqInfo.BlueprintName)
	for _, method := range methods {
		info := &BlueprintMethodInfo{
			Filename:   reqInfo.FileName,
			MethodName: method,
		}
		if err := m.appendBlueprintMethod(ctx, key, info); err != nil {
			log.Error(err)
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) copyBlueprint(ctx *gin.Context) {
	fmt.Println(11111)
	userid, _ := ctx.Cookie("userid")
	type BlueprintReq struct {
		Token   string `json:"token"`
		NewName string `json:"newName"`
	}
	var reqInfo BlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		log.Error(err)
		return
	}
	claims, err := m.JwtDecode(ctx, reqInfo.Token)
	if err != nil {
		log.Error(err)
		return
	}
	fromUserid := claims["userid"].(string)
	blueprintName := claims["blueprintName"].(string)
	fromKey := m.fullKey(fromUserid, blueprintName)
	toKey := m.fullKey(userid, reqInfo.NewName)
	ok, err := m.CopyBlueprint(ctx, fromKey, toKey)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, ok)
}

func (m Management) CopyBlueprint(ctx context.Context, fromKey string, toKey string) (bool, error) {
	fromValue, _, err := etcd.Dao.Get(fromKey)
	if err != nil {
		return false, err
	}
	_, ok, err := etcd.Dao.Get(toKey)
	if err != nil {
		return false, err
	}
	if ok {
		return false, nil
	}
	if err := etcd.Dao.Put(toKey, fromValue); err != nil {
		return false, err
	}
	return true, nil
}

func (m Management) renameBlueprint(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	type BlueprintReq struct {
		BlueprintName string `json:"blueprintName"`
	}
	var reqInfo BlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		log.Error(err)
		return
	}
	key := m.fullKey(userid, reqInfo.BlueprintName)
	value, _, err := etcd.Dao.Get(key)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
	}
	//todo 事务
	if err = etcd.Dao.Delete(key); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
	}
	if err = etcd.Dao.Put(key, value); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) savetoBlueprint(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	type BlueprintReq struct {
		BlueprintName string `json:"blueprintName"`
		FileName      string `json:"filename"`
		MethodName    string `json:"methodName"`
	}
	var reqInfo BlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		log.Error(err)
		return
	}
	reqInfo.BlueprintName = m.trim(reqInfo.BlueprintName)
	reqInfo.FileName = m.trim(reqInfo.FileName)
	reqInfo.MethodName = m.trim(reqInfo.MethodName)
	split := strings.Split(reqInfo.MethodName, "/")
	methodRawName := split[len(split)-1]
	key := m.fullKey(userid, reqInfo.BlueprintName)
	methodInfo := &BlueprintMethodInfo{
		Filename:   reqInfo.FileName,
		MethodName: methodRawName,
	}
	if err := m.appendBlueprintMethod(ctx, key, methodInfo); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) appendBlueprint(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	type BlueprintReq struct {
		BlueprintName string   `json:"blueprintName"`
		FileName      string   `json:"filename"`
		MethodName    []string `json:"methodName"`
	}
	var reqInfo BlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	reqInfo.BlueprintName = m.trim(reqInfo.BlueprintName)
	reqInfo.FileName = m.trim(reqInfo.FileName)
	for _, method := range reqInfo.MethodName {
		method = m.trim(method)
		split := strings.Split(method, "/")
		methodRawName := split[len(split)-1]
		key := m.fullKey(userid, reqInfo.BlueprintName)
		methodInfo := &BlueprintMethodInfo{
			Filename:   reqInfo.FileName,
			MethodName: methodRawName,
		}
		if err := m.appendBlueprintMethod(ctx, key, methodInfo); err != nil {
			log.Error(err)
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) deleteBlueprintMethod(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	type BlueprintReq struct {
		BlueprintName string `json:"blueprintName"`
		FileName      string `json:"fileName"`
		MethodName    string `json:"methodName"`
	}
	var reqInfo BlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		log.Error(err)
		return
	}
	reqInfo.BlueprintName = m.trim(reqInfo.BlueprintName)
	reqInfo.MethodName = m.trim(reqInfo.MethodName)
	split := strings.Split(reqInfo.MethodName, "/")
	methodRawName := split[len(split)-1]
	key := m.fullKey(userid, reqInfo.BlueprintName)
	method := &BlueprintMethodInfo{
		Filename:   reqInfo.FileName,
		MethodName: methodRawName,
	}
	if err := m.reduceBlueprintMethod(ctx, key, method); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) shareBlueprint(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	type BlueprintReq struct {
		BlueprintName string `json:"blueprintName"`
	}
	var reqInfo BlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		log.Error(err)
		return
	}
	reqInfo.BlueprintName = m.trim(reqInfo.BlueprintName)
	var hmacSampleSecret []byte
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":        userid,
		"blueprintName": reqInfo.BlueprintName,
	})
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
	}
	url := "/management/public?token=" + tokenString
	ctx.JSON(http.StatusOK, url)
}

func (m Management) delBlueprint(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	type BlueprintReq struct {
		BlueprintName string `json:"blueprintName"`
	}
	var reqInfo BlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		log.Error(err)
		return
	}
	reqInfo.BlueprintName = m.trim(reqInfo.BlueprintName)
	key := m.fullKey(userid, reqInfo.BlueprintName)
	if err := etcd.Dao.Delete(key); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) allUserBlueprints(ctx context.Context, userid string) []string {
	var out []string
	prefix := m.fullKey(userid, "")
	keys, _, err := etcd.Dao.GetWithPrefix(prefix)
	if err != nil {
		log.Error("Failed to get user blueprints: %+v", err)
		return nil
	}
	for _, key := range keys {
		blueprintName := strings.TrimPrefix(key, prefix)
		out = append(out, blueprintName)
	}
	return out
}

func (m Management) allUserBlueprintsMeta(ctx context.Context, userid string) []string {
	var out []string
	prefix := m.fullKey(userid, "")
	keys, _, err := etcd.Dao.GetWithPrefix(prefix)
	if err != nil {
		log.Error("Failed to get user blueprints: %+v", err)
		return nil
	}
	for _, key := range keys {
		blueprintName := strings.TrimPrefix(key, prefix)
		out = append(out, blueprintName)
	}
	return out
}

func (m Management) blueprintMethods(ctx context.Context, userid string, blueprintIdentifier string) ([]*BlueprintMethodInfo, error) {
	info := []*BlueprintMethodInfo{}
	key := m.fullKey(userid, blueprintIdentifier)
	value, _, err := etcd.Dao.Get(key)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(value), &info); err != nil {
		return nil, err
	}
	return info, nil
}

func (m Management) appendBlueprintMethod(ctx context.Context, key string, newMethod *BlueprintMethodInfo) error {
	info := []*BlueprintMethodInfo{}
	value, _, err := etcd.Dao.Get(key)
	if err != nil {
		return err
	}
	if value != "" {
		if err := json.Unmarshal([]byte(value), &info); err != nil {
			return err
		}
		for _, method := range info {
			if *method == *newMethod {
				return nil
			}
		}
	}
	info = append(info, newMethod)
	infoByte, err := json.Marshal(info)
	if err != nil {
		return err
	}
	if err := etcd.Dao.Put(key, string(infoByte)); err != nil {
		return err
	}
	return nil
}

func (m Management) reduceBlueprintMethod(ctx context.Context, key string, deleteMethod *BlueprintMethodInfo) error {
	info := []*BlueprintMethodInfo{}
	value, _, err := etcd.Dao.Get(key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(value), &info); err != nil {
		return err
	}
	for i, method := range info {
		if *method == *deleteMethod {
			info = append(info[:i], info[i+1:]...)
			break
		}
	}
	infoByte, err := json.Marshal(info)
	if err != nil {
		return err
	}
	if err := etcd.Dao.Put(key, string(infoByte)); err != nil {
		return err
	}
	return nil
}

func (m Management) putBlueprint(key string, info []*BlueprintMethodInfo) error {
	infoByte, err := json.Marshal(info)
	if err != nil {
		return err
	}
	if err := etcd.Dao.Put(key, string(infoByte)); err != nil {
		return err
	}
	return nil
}

func (m Management) fullKey(userid string, blueprintIdentifier string) string {
	//todo: 键设计
	return "/blueprint/" + userid + "/" + blueprintIdentifier
}

func (m Management) trim(str string) string {
	return strings.Replace(str, " ", "", -1)
}

func (m Management) getMethodsByService(ctx context.Context, serviceIdentifier string) ([]string, error) {
	var methods []string
	fileProfile, ok := m.findProtoFileByServiceIdentifier(ctx, serviceIdentifier)
	if !ok {
		return nil, errors.Errorf("Failed to find package profile from service identifier: %q", serviceIdentifier)
	}
	for _, s := range fileProfile.ProtoPackage.FileDescriptor.GetServices() {
		for _, m := range s.GetMethods() {
			methods = append(methods, m.GetName())
		}
	}
	return methods, nil
}
