package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
)

type SysTenant struct {
	ent.Schema
}

func (SysTenant) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty().MaxLen(128).Comment("名称 | Name"),
		field.String("code").Unique().NotEmpty().MaxLen(128).Comment("编码 | Code"),
		field.Int64("admin_id").Positive().Comment("管理员ID | Admin ID"),
		field.Int64("parent_id").Optional().Default(0).Comment("父级ID | Parent ID"),
		field.Int64("package_id").Optional().Default(0).Comment("套餐ID | Package ID"),
		field.Time("expired_at").Optional().Comment("到期时间 | Expired At"),
	}
}

func (SysTenant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("sys_users", SysUser.Type),
		edge.To("sys_depts", SysDept.Type),
		edge.From("sys_package", SysPackage.Type).Ref("tenants").Unique().Field("package_id"),
	}
}

func (SysTenant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.BaseMixin{},
		mixins.CreatedMixin{},
		mixins.UpdatedMixin{},
		mixins.DeletedMixin{},
		mixins.StatusMixin{},
		mixins.RemarkMixin{},
	}
}

func (SysTenant) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("Tenant Table | 租户表"),
		entsql.Annotation{Table: "sys_tenants"},
	}
}

func (SysTenant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id", "status"),
		index.Fields("parent_id", "status"),
	}
}
