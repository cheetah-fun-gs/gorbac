package gorbac

import "time"

// Role 角色
type Role struct {
	UID      string    `json:"uid,omitempty"`
	Level    RoleLevel `json:"level,omitempty"`
	ExpireAt int64     `json:"expire_at,omitempty"`
	Actions  []string  `json:"actions,omitempty"` // 允许的所有动作
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
	return nil
}

// GetObjectRoles 获取object的角色 参数可为 nil 表示全部
func (rbac *RBAC) GetObjectRoles(object string, group *string, isExtend *bool) ([]ObjectRole, error) {
	return nil, nil
}

// AddGroupRole 添加group的角色
func (rbac *RBAC) AddGroupRole(group, uid string, roleLevel RoleLevel, expireAts ...time.Time) error {
	return nil
}

// RemoveGroupRole 移除group的角色
func (rbac *RBAC) RemoveGroupRole(group, uid string) error {
	return nil
}

// GetGroupRoles 获取group的角色
func (rbac *RBAC) GetGroupRoles(group string) ([]GroupRole, error) {
	return nil, nil
}

// AddGroupObject 添加group的object
func (rbac *RBAC) AddGroupObject(group, object string) error {
	return nil
}

// RemoveGroupObject 移除group的object
func (rbac *RBAC) RemoveGroupObject(group, object string) error {
	return nil
}

// GetGroupObjects 获取group的object
func (rbac *RBAC) GetGroupObjects(group string) ([]string, error) {
	return nil, nil
}

// GetUserObjectRoles 获得用户所有object的role 参数可为 nil 表示全部
func (rbac *RBAC) GetUserObjectRoles(uid string, group *string, isExtend *bool) ([]ObjectRole, error) {
	return nil, nil
}

// GetUserGroupRoles 获得用户的所有group的role
func (rbac *RBAC) GetUserGroupRoles(uid string) ([]GroupRole, error) {
	return nil, nil
}
