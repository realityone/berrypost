package management

import (
	"context"
	"encoding/json"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := m.NewBlueprint(ctx, m.fullKey(userid, reqInfo.BlueprintName)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) NewBlueprint(ctx context.Context, key string) error {
	ok, err := etcd.Dao.Exist(key)
	if err != nil {
		log.Error("%+v", err)
		return err
	}
	if ok {
		return errors.New("blueprint name already exists")
	}
	if err := m.putBlueprint(key, nil); err != nil {
		log.Error("%+v", err)
		return err
	}
	return nil
}

func (m Management) copyBlueprintFromFile(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	var reqInfo *CopyBlueprintFromFileReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := m.CopyBlueprintFromFile(ctx, userid, reqInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) CopyBlueprintFromFile(ctx context.Context, userid string, req *CopyBlueprintFromFileReq) error {
	blueprintName := m.trim(req.BlueprintName)
	fileName := m.trim(req.FileName)
	methods, err := m.getMethodsByService(ctx, fileName)
	if err != nil {
		log.Error("%+v", err)
		return err
	}
	key := m.fullKey(userid, blueprintName)
	for _, method := range methods {
		info := &BlueprintMethodInfo{
			Filename:   fileName,
			MethodName: method,
		}
		if err := m.appendBlueprintMethod(ctx, key, info); err != nil {
			log.Error("%+v", err)
			return err
		}
	}
	return nil
}

func (m Management) copyBlueprint(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	var reqInfo *CopyBlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := m.CopyBlueprint(ctx, userid, reqInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) CopyBlueprint(ctx context.Context, userid string, req *CopyBlueprintReq) error {
	claims, err := m.JwtDecode(ctx, req.Token)
	if err != nil {
		log.Error("%+v", err)
		return err
	}
	fromUserid := claims["userid"].(string)
	blueprintName := claims["blueprintName"].(string)
	fromKey := m.fullKey(fromUserid, blueprintName)
	toKey := m.fullKey(userid, req.NewName)
	ok, err := etcd.Dao.Exist(toKey)
	if err != nil {
		return err
	}
	if ok {
		return errors.New("blueprint name already exists")
	}
	if err := m.CopyObject(ctx, fromKey, toKey); err != nil {
		log.Error("%+v", err)
		return err
	}
	return nil
}

func (m Management) CopyObject(ctx context.Context, fromKey string, toKey string) error {
	fromValue, err := etcd.Dao.GetIfExist(fromKey)
	if err != nil {
		return err
	}
	if err := etcd.Dao.Put(toKey, fromValue); err != nil {
		return err
	}
	return nil
}

func (m Management) renameBlueprint(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	type BlueprintReq struct {
		OldName string `json:"oldName"`
		NewName string `json:"newName"`
	}
	var reqInfo BlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	oldKey := m.fullKey(userid, reqInfo.OldName)
	newKey := m.fullKey(userid, reqInfo.NewName)
	if err := m.RenameObject(ctx, oldKey, newKey); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) RenameObject(ctx context.Context, oldKey string, newKey string) error {
	value, err := etcd.Dao.GetIfExist(oldKey)
	if err != nil {
		log.Error("%+v", err)
		return err
	}
	if err = etcd.Dao.Update(newKey, value); err != nil {
		log.Error("%+v", err)
		return err
	}
	return nil
}

func (m Management) appendBlueprint(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	var reqInfo *BlueprintMethodReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := m.AppendBlueprint(ctx, userid, reqInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) AppendBlueprint(ctx context.Context, userid string, req *BlueprintMethodReq) error {
	blueprintName := m.trim(req.BlueprintName)
	fileName := m.trim(req.FileName)
	methodName := m.trim(req.MethodName)
	split := strings.Split(methodName, "/")
	methodRawName := split[len(split)-1]
	key := m.fullKey(userid, blueprintName)
	methodInfo := &BlueprintMethodInfo{
		Filename:   fileName,
		MethodName: methodRawName,
	}
	if err := m.appendBlueprintMethod(ctx, key, methodInfo); err != nil {
		log.Error("%+v", err)
		return err
	}
	return nil
}

func (m Management) listAppendBlueprint(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	var reqInfo *ListAppendBlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := m.ListAppendBlueprint(ctx, userid, reqInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) ListAppendBlueprint(ctx context.Context, userid string, req *ListAppendBlueprintReq) error {
	appendInfo := &BlueprintMethodReq{
		BlueprintName: req.BlueprintName,
		FileName:      req.FileName,
	}
	for _, method := range req.MethodName {
		appendInfo.MethodName = method
		if err := m.AppendBlueprint(ctx, userid, appendInfo); err != nil {
			log.Error("%+v", err)
			return err
		}
	}
	return nil
}

func (m Management) deleteBlueprintMethod(ctx *gin.Context) {
	var reqInfo *BlueprintMethodReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := m.DeleteBlueprintMethod(ctx, reqInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) DeleteBlueprintMethod(ctx *gin.Context, req *BlueprintMethodReq) error {
	userid, _ := ctx.Cookie("userid")
	blueprintName := m.trim(req.BlueprintName)
	methodName := m.trim(req.MethodName)
	split := strings.Split(methodName, "/")
	methodRawName := split[len(split)-1]
	key := m.fullKey(userid, blueprintName)
	method := &BlueprintMethodInfo{
		Filename:   req.FileName,
		MethodName: methodRawName,
	}
	if err := m.reduceBlueprintMethod(ctx, key, method); err != nil {
		log.Error("%+v", err)
		return err
	}
	historyKey := m.historyKey(userid, blueprintName, methodName)
	ok, err := etcd.Dao.Exist(historyKey)
	if err != nil {
		log.Error("%+v", err)
		return err
	}
	if !ok {
		return nil
	}
	if err = etcd.Dao.Delete(historyKey); err != nil {
		return err
	}
	return nil
}

func (m Management) reduceBlueprintMethod(ctx context.Context, key string, deleteMethod *BlueprintMethodInfo) error {
	info := []*BlueprintMethodInfo{}
	value, err := etcd.Dao.GetIfExist(key)
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

func (m Management) shareBlueprint(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	type BlueprintReq struct {
		BlueprintName string `json:"blueprintName"`
	}
	var reqInfo BlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	url, err := m.ShareURL(ctx, userid, reqInfo.BlueprintName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, url)
}

func (m Management) ShareURL(ctx *gin.Context, userid string, blueprintName string) (string, error) {
	blueprintName = m.trim(blueprintName)
	var hmacSampleSecret []byte
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":        userid,
		"blueprintName": blueprintName,
	})
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		log.Error("%+v", err)
		return "", err
	}
	url := "/management/public?token=" + tokenString
	return url, nil
}

func (m Management) delBlueprint(ctx *gin.Context) {
	type BlueprintReq struct {
		BlueprintName string `json:"blueprintName"`
	}
	var reqInfo BlueprintReq
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := m.DeleteBlueprint(ctx, reqInfo.BlueprintName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) DeleteBlueprint(ctx *gin.Context, blueprintName string) error {
	userid, _ := ctx.Cookie("userid")
	blueprintName = m.trim(blueprintName)
	key := m.fullKey(userid, blueprintName)
	if err := etcd.Dao.Delete(key); err != nil {
		return err
	}
	history := m.historyPrefix(userid, blueprintName)
	if err := etcd.Dao.DeleteWithPrefix(history); err != nil {
		return err
	}
	return nil
}

func (m Management) allUserBlueprints(ctx context.Context, userid string) []string {
	var out []string
	prefix := m.fullKey(userid, "")
	keys, _, err := etcd.Dao.GetKVWithPrefix(prefix)
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
	value, err := etcd.Dao.GetIfExist(key)
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
	value, err := etcd.Dao.GetIfExist(key)
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
