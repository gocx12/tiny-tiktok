// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	messageFieldNames          = builder.RawFieldNames(&Message{})
	messageRows                = strings.Join(messageFieldNames, ",")
	messageRowsExpectAutoSet   = strings.Join(stringx.Remove(messageFieldNames, "`id`","`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	messageRowsWithPlaceHolder = strings.Join(stringx.Remove(messageFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	messageModel interface {
		Insert(ctx context.Context, data *Message) (sql.Result, error)
		FindList(ctx context.Context, id int64,time string) ([]Message, error)
	}

	defaultMessageModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Message struct {
		Id         int64     `db:"id"`
		FromUserId int64     `db:"from_user_id"`
		ToUserId   int64     `db:"to_user_id"`
		Content    string    `db:"content"`
		CreateAt   time.Time `db:"created_at"`
	}
)

func newMessageModel(conn sqlx.SqlConn) *defaultMessageModel {
	return &defaultMessageModel{
		conn:  conn,
		table: "`message`",
	}
}


func (m *defaultMessageModel) FindList(ctx context.Context, id int64,time string) ([]Message, error) {
	query := "select `id`,`from_user_id`,`to_user_id`,`content`,`created_at` from message where `to_user_id` = ? AND DATE_FORMAT(created_at,'%Y-%m-%d %H:%i:%s') > ?  ORDER BY created_at"
	var resp []Message
	err := m.conn.QueryRows(&resp, query, id,time)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMessageModel) Insert(ctx context.Context, data *Message) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, messageRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.FromUserId, data.ToUserId, data.Content)
	return ret, err
}



func (m *defaultMessageModel) tableName() string {
	return m.table
}
