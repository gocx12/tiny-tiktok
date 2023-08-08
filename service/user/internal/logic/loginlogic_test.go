package logic

import (
	"context"
	"flag"
	"testing"
	"tiny-tiktok/service/user/internal/config"
	"tiny-tiktok/service/user/internal/svc"
	"tiny-tiktok/service/user/pb/user"

	"github.com/zeromicro/go-zero/core/conf"
)

// var db sqlx.SqlConn

var configFile = flag.String("f", "../../etc/user.yaml", "the config file")
var ctx = context.Background()
var svcCtx *svc.ServiceContext

func TestMain(m *testing.M) {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx = svc.NewServiceContext(c)

	m.Run()
}
func TestLogin(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		testName string
		args     args
		wantErr  bool
	}{
		{
			testName: "existing name and right password",
			args: args{
				username: "gao",
				password: "123456",
			},
			wantErr: false,
		},
		{
			testName: "not existing name",
			args: args{
				username: "Alex",
				password: "123456",
			},
			wantErr: true,
		},
		{
			testName: "existing name and wrong password",
			args: args{
				username: "gao",
				password: "wrong_password",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			_, err := NewLoginLogic(ctx, svcCtx).Login(&user.LoginReq{
				Username: test.args.username,
				Password: test.args.password,
			})
			if (err != nil) != test.wantErr {
				t.Errorf("Login() error: %v, wantErr %v", err, test.wantErr)
				return
			}
		})
	}
}
