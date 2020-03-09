package gorbac

import (
	"database/sql"
	"fmt"
	"time"

	sqlplus "github.com/cheetah-fun-gs/goplus/dao/sql"
)

// Role 角色
type Role struct {
	UID      string    `json:"uid,omitempty"`
	Level    RoleLevel `json:"level,omitempty"`
	ExpireAt int64     `json:"expire_at,omitempty"`
	Actions  []string  `json:"actions,omitempty"` // 允许的所有动作
}

func (role *Role) fillActions(actionLevels ActionLevels) {
	for action, minLevel := range actionLevels {
		if minLevel <= role.Level {
			role.Actions = append(role.Actions, action)
		}
	}
}

// ObjectRole ...
type ObjectRole struct {
	Role
	Object   string
	Group    string
	IsExtend bool // 是否继承自Group
}

// GroupRole ...
type GroupRole struct {
	Role
	Group string
}

// AddObjectRole 添加object的角色
func (rbac *RBAC) AddObjectRole(object, uid string, roleLevel RoleLevel, expireAts ...time.Time) error {
	return nil
}

// RemoveObjectRole 移除object的角色
func (rbac *RBAC) RemoveObjectRole(object, uid string) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE name = ? AND is_group = 0 AND uid = ?",
		rbac.tableObjectRoles.Name)
	args := []interface{}{object, uid}

	_, err := rbac.db.Exec(query, args...)
	return err
}

// GetObjectRoles 获取object的角色 isExcludeExtend 是否排除继承的角色, groupOrBlank 非空指定组, 为空不指定组
func (rbac *RBAC) GetObjectRoles(object string, isExcludeExtend bool, groupOrBlank string) ([]ObjectRole, error) {
	return nil, nil
}

// AddGroupRole 添加group的角色
func (rbac *RBAC) AddGroupRole(group, uid string, roleLevel RoleLevel, expireAts ...time.Time) error {
	return nil
}

// RemoveGroupRole 移除group的角色
func (rbac *RBAC) RemoveGroupRole(group, uid string) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE name = ? AND is_group = 1 AND uid = ?",
		rbac.tableObjectRoles.Name)
	args := []interface{}{group, uid}

	_, err := rbac.db.Exec(query, args...)
	return err
}

// GetGroupRoles 获取group的角色
func (rbac *RBAC) GetGroupRoles(group string) ([]GroupRole, error) {
	now := time.Now()
	result := []GroupRole{}

	query := fmt.Sprintf("SELECT * FROM %v WHERE name = ? AND is_group = 1 AND (expire_at is null || expire_at > ?)",
		rbac.tableObjectRoles.Name)
	args := []interface{}{group, now}

	rows, err := rbac.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	modelObjectRoles := []ModelObjectRole{}
	if err = sqlplus.Select(rows, &modelObjectRoles); err == sql.ErrNoRows {
		return result, nil
	} else if err != nil {
		return nil, err
	}

	for _, data := range modelObjectRoles {
		var expireAt int64
		if data.ExpireAt.Valid {
			expireAt = data.ExpireAt.Time.Unix()
		}
		role := GroupRole{
			Role: Role{
				UID:      data.UID,
				Level:    RoleLevel(data.Level),
				ExpireAt: expireAt,
			},
			Group: data.Name,
		}
		role.fillActions(rbac.actionLevels)
		result = append(result, role)
	}
	return result, nil
}

// AddGroupObject 添加group的object
func (rbac *RBAC) AddGroupObject(group, object string) error {
	return nil
}

// RemoveGroupObject 移除group的object
func (rbac *RBAC) RemoveGroupObject(group, object string) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE group = ? AND object = ?", rbac.tableGroupObjects.Name)
	args := []interface{}{group, object}

	_, err := rbac.db.Exec(query, args...)
	return err
}

// GetGroupObjects 获取group的object
func (rbac *RBAC) GetGroupObjects(group string) ([]string, error) {
	query := fmt.Sprintf("SELECT * FROM %v WHERE group = ?", rbac.tableGroupObjects.Name)
	args := []interface{}{group}

	rows, err := rbac.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	modelGroupObjects := []ModelGroupObject{}
	if err = sqlplus.Select(rows, &modelGroupObjects); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	result := []string{}
	for _, data := range modelGroupObjects {
		result = append(result, data.Object)
	}
	return result, nil
}

// GetUserObjectRoles 获得用户所有object的role: objectOrBlank 非空指定对象, 为空不指定对象, isExcludeExtend 是否排除继承的角色
func (rbac *RBAC) GetUserObjectRoles(uid string, objectOrBlank string, isExcludeExtend bool) ([]ObjectRole, error) {
	if isExcludeExtend {
		return rbac.dbGetObjectRolesExcludeExtend(uid, objectOrBlank)
	}

	return rbac.dbGetObjectRolesIncludeExtend(uid, objectOrBlank)
}

// GetUserGroupRoles 获得用户的所有group的role: groupOrBlank 非空指定组, 为空不指定组
func (rbac *RBAC) GetUserGroupRoles(uid string, groupOrBlank string) ([]GroupRole, error) {
	now := time.Now()
	result := []GroupRole{}

	query := fmt.Sprintf("SELECT * FROM %v WHERE uid = ? AND is_group = 1 AND (expire_at is null || expire_at > ?)",
		rbac.tableObjectRoles.Name)
	args := []interface{}{uid, now}
	if groupOrBlank != "" {
		query += " AND name = ?"
		args = append(args, groupOrBlank)
	}

	rows, err := rbac.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	modelObjectRoles := []ModelObjectRole{}
	if err = sqlplus.Select(rows, &modelObjectRoles); err == sql.ErrNoRows {
		return result, nil
	} else if err != nil {
		return nil, err
	}

	for _, data := range modelObjectRoles {
		var expireAt int64
		if data.ExpireAt.Valid {
			expireAt = data.ExpireAt.Time.Unix()
		}
		role := GroupRole{
			Role: Role{
				UID:      uid,
				Level:    RoleLevel(data.Level),
				ExpireAt: expireAt,
			},
			Group: data.Name,
		}
		role.fillActions(rbac.actionLevels)
		result = append(result, role)
	}
	return result, nil
}
