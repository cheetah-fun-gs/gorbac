package gorbac

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	sqlplus "github.com/cheetah-fun-gs/goplus/dao/sql"
)

func objectRoleAppend(roles []ObjectRole, appendRoles ...ObjectRole) []ObjectRole {
	for _, appendRole := range appendRoles {
		var isMatch bool
		for index, role := range roles {
			if appendRole.Object == role.Object {
				isMatch = true
				if role.Level < appendRole.Level {
					roles[index] = appendRole // 高level替换低level
				}
			}
		}
		if !isMatch {
			roles = append(roles, appendRole)
		}
	}
	return roles
}

// 获取不含继承的角色列表
func (rbac *RBAC) dbGetObjectRolesExcludeExtend(uid string, objectOrBlank string) ([]ObjectRole, error) {
	now := time.Now()
	result := []ObjectRole{}

	query := fmt.Sprintf("SELECT * FROM %v WHERE uid = ? AND is_group = 0 AND (expire_at is null || expire_at > ?)",
		rbac.tableObjectRoles.Name)
	args := []interface{}{uid, now}
	if objectOrBlank != "" {
		query += " AND name = ?"
		args = append(args, objectOrBlank)
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
		role := ObjectRole{
			Role: Role{
				UID:      uid,
				Level:    RoleLevel(data.Level),
				ExpireAt: expireAt,
			},
			Object:   data.Name,
			Group:    "",
			IsExtend: false,
		}
		role.fillActions(rbac.actionLevels)
		result = append(result, role)
	}
	return result, nil
}

// 获取包含继承的角色列表
func (rbac *RBAC) dbGetObjectRolesIncludeExtend(uid string, objectOrBlank string) ([]ObjectRole, error) {
	now := time.Now()
	result := []ObjectRole{}

	tx, err := rbac.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit() // 只读， 提交就好，无需回滚

	// 获取数据
	query := fmt.Sprintf("SELECT * FROM %v WHERE uid = ? AND (expire_at is null || expire_at > ?)",
		rbac.tableObjectRoles.Name)
	args := []interface{}{uid, now}

	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}

	modelObjectRoles := []ModelObjectRole{}
	if err = sqlplus.Select(rows, &modelObjectRoles); err == sql.ErrNoRows {
		return result, nil
	} else if err != nil {
		return nil, err
	}

	marks := []string{}
	args = []interface{}{}
	for _, data := range modelObjectRoles {
		if data.IsGroup {
			marks = append(marks, "?")
			args = append(args, data.Name)
		}
	}

	query = fmt.Sprintf("SELECT * FROM %v WHERE group in (%v)",
		rbac.tableGroupObjects.Name, strings.Join(marks, ", "))
	if objectOrBlank != "" {
		query += " AND object = ?"
		args = append(args, objectOrBlank)
	}

	rows, err = tx.Query(query, args...)
	if err != nil {
		return nil, err
	}

	modelGroupObjects := []ModelGroupObject{}
	if err = sqlplus.Select(rows, &modelGroupObjects); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 整合数据
	// 非继承的角色
	for _, data := range modelObjectRoles {
		if !data.IsGroup && (objectOrBlank == "" || data.Name == objectOrBlank) {
			var expireAt int64
			if data.ExpireAt.Valid {
				expireAt = data.ExpireAt.Time.Unix()
			}
			role := ObjectRole{
				Role: Role{
					UID:      uid,
					Level:    RoleLevel(data.Level),
					ExpireAt: expireAt,
				},
				Object:   data.Name,
				Group:    "",
				IsExtend: false,
			}
			role.fillActions(rbac.actionLevels)
			result = append(result, role)
		}
	}

	// 继承的角色
	for _, data := range modelObjectRoles {
		if data.IsGroup {
			for _, dd := range modelGroupObjects {
				if dd.Group == data.Name && (objectOrBlank == "" || dd.Object == objectOrBlank) {
					var expireAt int64
					if data.ExpireAt.Valid {
						expireAt = data.ExpireAt.Time.Unix()
					}
					role := ObjectRole{
						Role: Role{
							UID:      uid,
							Level:    RoleLevel(data.Level),
							ExpireAt: expireAt,
						},
						Object:   dd.Object,
						Group:    dd.Group,
						IsExtend: true,
					}
					role.fillActions(rbac.actionLevels)
					result = objectRoleAppend(result, role)
				}
			}
		}
	}
	return result, nil
}
