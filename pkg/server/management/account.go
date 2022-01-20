package management

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/realityone/berrypost/pkg/etcd"
	log "github.com/sirupsen/logrus"
)

func (m Management) login(ctx *gin.Context) {
	_, err := ctx.Cookie("userid")
	if err == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/management/rediect-to-example")
		return
	}
	ctx.HTML(http.StatusOK, "login.html", &LoginPage{
		Meta: m.server.Meta(),
	})
}

func (m Management) register(ctx *gin.Context) {
	_, err := ctx.Cookie("userid")
	if err == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/management/rediect-to-example")
		return
	}
	ctx.HTML(http.StatusOK, "register.html", &LoginPage{
		Meta: m.server.Meta(),
	})
}

func (m Management) signOut(ctx *gin.Context) {
	ctx.SetCookie("session", "", -1, "", "", true, true)
	ctx.SetCookie("userid", "", -1, "", "", true, true)
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) signIn(ctx *gin.Context) {
	type RespBody struct {
		Userid   string `json:"userid"`
		Password string `json:"password"`
	}
	var reqInfo RespBody
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, err := m.userSignIn(ctx, reqInfo.Userid, reqInfo.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		log.Error("%+v", err)
		return
	}
	ctx.JSON(http.StatusOK, ok)
}

func (m Management) signUp(ctx *gin.Context) {
	type RespBody struct {
		Userid   string `json:"userid"`
		Password string `json:"password"`
	}
	var reqInfo RespBody
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, err := m.userSignUp(ctx, reqInfo.Userid, reqInfo.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		log.Error("%+v", err)
		return
	}
	ctx.JSON(http.StatusOK, ok)
}

func (m Management) userSignIn(ctx *gin.Context, userid string, password string) (bool, error) {
	ok, err := m.userVerify(ctx, userid, password)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	var hmacSampleSecret []byte
	ok, err = m.isAdmin(ctx, userid)
	if err != nil {
		return false, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": userid,
		"admin":  ok,
	})
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return false, err
	}
	ctx.SetCookie("session", tokenString, 3600, "", "", true, true)
	ctx.SetCookie("userid", userid, 3600, "", "", true, true)
	return true, nil
}

func (m Management) userVerify(ctx context.Context, userid string, password string) (bool, error) {
	userKey := m.userKey(userid)
	ok, err := etcd.Dao.Exist(userKey)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	pw, err := etcd.Dao.GetIfExist(userKey)
	if err != nil {
		return false, err
	}
	if pw != password {
		return false, nil
	}
	return true, nil
}

func (m Management) isAdmin(ctx context.Context, userid string) (bool, error) {
	userKey := m.adminKey(userid)
	ok, err := etcd.Dao.Exist(userKey)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (m Management) userSignUp(ctx *gin.Context, userid string, password string) (bool, error) {
	userKey := m.userKey(userid)
	ok, err := etcd.Dao.Exist(userKey)
	if err != nil {
		return false, err
	}
	if ok {
		return false, nil
	}
	if err = etcd.Dao.Put(userKey, password); err != nil {
		return false, err
	}
	return true, nil
}
