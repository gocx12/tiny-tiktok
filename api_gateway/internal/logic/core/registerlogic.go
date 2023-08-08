package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"tiny-tiktok/api_gateway/internal/svc"
	"tiny-tiktok/api_gateway/internal/types"
	"tiny-tiktok/service/user/pb/user"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// todo: add your logic here and delete this line
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "user.rpc",
		},
	})
	client := user.NewUserServiceClient(conn.Conn())
	respRpc, err := client.Register(context.Background(), &user.RegisterReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		resp = &types.RegisterResp{
			StatusCode: http.StatusOK,
			StatusMsg:  "Register fail",
			UserID:     respRpc.UserId, // is -1
		}
		log.Fatal(err)
		err = nil
		return
	} else if respRpc.UserId == -1 {
		// the username does not exsit or the password is incorrect
		resp = &types.RegisterResp{
			StatusCode: http.StatusOK,
			StatusMsg:  respRpc.StatusMsg,
			UserID:     respRpc.UserId, // is -1
		}
		err = nil
		return
	}

	secretKey := l.svcCtx.Config.Auth.AccessSecret
	iat := time.Now().Unix() // maybe not word on Windows OS
	seconds := l.svcCtx.Config.Auth.AccessExpire
	payload := fmt.Sprintf("%d,%s", respRpc.UserId, req.Username)

	token, err := getJwtToken(secretKey, iat, seconds, payload)
	if err != nil {
		resp = &types.RegisterResp{
			StatusCode: http.StatusOK,
			StatusMsg:  "Login fail",
			UserID:     respRpc.UserId, // is -1
		}
		log.Fatal(err)
		err = nil
		return
	}

	resp = &types.RegisterResp{
		StatusCode: http.StatusOK,
		StatusMsg:  respRpc.StatusMsg,
		UserID:     respRpc.UserId,
		Token:      token,
	}

	return
}
