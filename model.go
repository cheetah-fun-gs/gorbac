package gorbac

import (
	"database/sql"
	"time"
)

// sql table
const (
	TableObjectRoles = `CREATE TABLE IF NOT EXISTS %v (
		id int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
		name varchar(45) NOT NULL COMMENT '对象名',
		is_group tinyint(4) unsigned NOT NULL COMMENT '是否组对象',
		uid varchar(128) NOT NULL COMMENT '用户ID',
		level int(10) unsigned NOT NULL COMMENT '角色级别',
		expire_at datetime DEFAULT NULL COMMENT '到期时间',
		created timestamp NOT NULL COMMENT '创建时间',
		updated timestamp NOT NULL COMMENT '更新时间',
		PRIMARY KEY (id),
		UNIQUE KEY uniq_object_uid (object, uid),
		KEY idx_uid (uid),
		KEY idx_created (created),
		KEY idx_updated (updated)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='对象角色表'`
	TableGroupObjects = `CREATE TABLE IF NOT EXISTS %v (
		id int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
		group varchar(45) NOT NULL COMMENT '组对象',
		object varchar(45) NOT NULL COMMENT '对象',
		created timestamp NOT NULL COMMENT '创建时间',
		PRIMARY KEY (id),
		UNIQUE KEY uniq_group_object (group, object),
		KEY idx_object (object),
		KEY idx_created (created)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组对象关系表'`
)

// ModelObjectRole 对象角色
type ModelObjectRole struct {
	ID       int          `json:"id,omitempty"`
	Name     string       `json:"name,omitempty"` // object name Or group name
	IsGroup  bool         `json:"is_group,omitempty"`
	UID      string       `json:"uid,omitempty"`
	Level    int          `json:"level,omitempty"`
	ExpireAt sql.NullTime `json:"expire_at,omitempty"`
	Created  time.Time    `json:"created,omitempty"`
	Updated  time.Time    `json:"updated,omitempty"`
}

// ModelGroupObject 组对象关系
type ModelGroupObject struct {
	ID      int       `json:"id,omitempty"`
	Group   string    `json:"group,omitempty"`
	Object  string    `json:"object,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}
