package gorbac

import (
	"database/sql"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
)

// RoleLevel 角色级别 自定义 要求: >0, 数字越大权限越高
type RoleLevel int

// ActionLevels 动作级别 key 为动作 val 为该动作要求的最低角色级别
type ActionLevels map[string]RoleLevel

// RBAC Role-Base Access Control
type RBAC struct {
	name              string
	roleLevels        []RoleLevel  // 角色集合
	actionLevels      ActionLevels // 动作集合
	db                *sql.DB
	pool              *redigo.Pool
	tableObjectRoles  *modelTable
	tableGroupObjects *modelTable
}

type modelTable struct {
	Name      string
	CreateSQL string
}

// New ...
func New(name string, roleLevels []RoleLevel, actionLevels ActionLevels,
	db *sql.DB, pools ...*redigo.Pool) *RBAC {
	tableObjectRolesName := name + "_object_roles"
	tableGroupObjectsName := name + "_group_objects"

	rbac := &RBAC{
		name:         name,
		roleLevels:   roleLevels,
		actionLevels: actionLevels,
		db:           db,
		tableObjectRoles: &modelTable{
			Name:      tableObjectRolesName,
			CreateSQL: fmt.Sprintf(TableObjectRoles, tableObjectRolesName),
		},
		tableGroupObjects: &modelTable{
			Name:      tableGroupObjectsName,
			CreateSQL: fmt.Sprintf(TableGroupObjects, tableGroupObjectsName),
		},
	}
	if len(pools) > 0 {
		rbac.pool = pools[0]
	}
	return rbac
}

// IsAllow 是否允许uid对object做action
func (rbac *RBAC) IsAllow(uid, object, action string) (bool, error) {
	return false, nil
}

// EnsureTables 确保sql表已建立
func (rbac *RBAC) EnsureTables() error {
	for _, createSQL := range []string{rbac.tableObjectRoles.CreateSQL, rbac.tableGroupObjects.CreateSQL} {
		if _, err := rbac.db.Exec(createSQL); err != nil {
			return err
		}
	}
	return nil
}

// TablesCreateSQL 获得建表语句
func (rbac *RBAC) TablesCreateSQL() []string {
	return []string{rbac.tableObjectRoles.CreateSQL, rbac.tableGroupObjects.CreateSQL}
}

// TablesName 获得表名
func (rbac *RBAC) TablesName() []string {
	return []string{rbac.tableObjectRoles.Name, rbac.tableGroupObjects.Name}
}
