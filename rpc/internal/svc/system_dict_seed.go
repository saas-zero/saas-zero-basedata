package svc

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysdict"
	"github.com/saas-zero/saas-zero-basedata/ent/sysdictdata"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
)

type SystemDictSeed struct {
	Name   string
	Key    string
	Remark string
	Items  []SystemDictDataSeed
}

type SystemDictDataSeed struct {
	Name  string
	Key   string
	Value string
}

// SystemDictSeeds contains the system-default dictionaries used by the UI.
// They are stored with tenant_id=0 and are inherited by every tenant.
var SystemDictSeeds = []SystemDictSeed{
	{
		Name:   "状态",
		Key:    "status",
		Remark: "系统状态选项",
		Items: []SystemDictDataSeed{
			{Name: "启用", Key: "active", Value: "active"},
			{Name: "禁用", Key: "inactive", Value: "inactive"},
			{Name: "暂停", Key: "suspended", Value: "suspended"},
			{Name: "冻结", Key: "frozen", Value: "frozen"},
			{Name: "过期", Key: "expired", Value: "expired"},
			{Name: "成功", Key: "success", Value: "success"},
			{Name: "失败", Key: "fail", Value: "fail"},
			{Name: "失败", Key: "failure", Value: "failure"},
		},
	},
	{
		Name:   "菜单类型",
		Key:    "menu_type",
		Remark: "菜单节点类型",
		Items: []SystemDictDataSeed{
			{Name: "目录", Key: "directory", Value: "directory"},
			{Name: "菜单", Key: "menu", Value: "menu"},
			{Name: "按钮", Key: "button", Value: "button"},
		},
	},
	{
		Name:   "API类型",
		Key:    "api_type",
		Remark: "API资源类型",
		Items: []SystemDictDataSeed{
			{Name: "分组", Key: "group", Value: "group"},
			{Name: "接口", Key: "api", Value: "api"},
		},
	},
}

// SeedSystemDicts inserts missing system dictionaries and items in one transaction.
func SeedSystemDicts(ctx context.Context, db *ent.Client) error {
	tx, err := db.Tx(ctx)
	if err != nil {
		return err
	}
	if err := SeedSystemDictsTx(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// SeedSystemDictsTx is used by InitAll so dictionary rows participate in its transaction.
func SeedSystemDictsTx(ctx context.Context, tx *ent.Tx) error {
	systemCtx := mixins.SetCurrentTenantId(ctx, 0)
	if mixins.GetCurrentUserId(systemCtx) <= 0 {
		systemCtx = mixins.SetCurrentUserId(systemCtx, 1)
	}
	if mixins.GetCurrentUserName(systemCtx) == "" {
		systemCtx = mixins.SetCurrentUserName(systemCtx, "system")
	}

	for _, seed := range SystemDictSeeds {
		dict, err := tx.SysDict.Query().
			Where(sysdict.KeyEQ(seed.Key), sysdict.TenantIDEQ(0), sysdict.DeletedAtIsNil()).
			Only(systemCtx)
		if ent.IsNotFound(err) {
			dict, err = tx.SysDict.Create().
				SetName(seed.Name).
				SetKey(seed.Key).
				SetStatus(sysdict.StatusActive).
				SetRemark(seed.Remark).
				Save(systemCtx)
		} else if err != nil {
			return err
		}
		if err != nil {
			return err
		}

		for _, item := range seed.Items {
			_, err := tx.SysDictData.Query().
				Where(
					sysdictdata.DictIDEQ(dict.ID),
					sysdictdata.KeyEQ(item.Key),
					sysdictdata.TenantIDEQ(0),
					sysdictdata.DeletedAtIsNil(),
				).
				Only(systemCtx)
			if ent.IsNotFound(err) {
				_, err = tx.SysDictData.Create().
					SetDictID(dict.ID).
					SetName(item.Name).
					SetKey(item.Key).
					SetValue(item.Value).
					SetStatus(sysdictdata.StatusActive).
					Save(systemCtx)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}
