package sysuserslogic

import (
	"context"
	"testing"

	"entgo.io/ent/dialect"
	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/enttest"
	"github.com/saas-zero/saas-zero-basedata/ent/migrate"
	"github.com/saas-zero/saas-zero-basedata/ent/sysrole"
	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"

	_ "github.com/mattn/go-sqlite3"
)

// testCtx returns a context with test user/tenant info injected
func testCtx() context.Context {
	ctx := context.Background()
	ctx = mixins.SetCurrentUserId(ctx, 1)
	ctx = mixins.SetCurrentUserName(ctx, "test-admin")
	ctx = mixins.SetCurrentTenantId(ctx, 1001)
	return ctx
}

// newTestClient creates an in-memory SQLite ent client for testing
func newTestClient(t *testing.T) *ent.Client {
	t.Helper()
	client := enttest.Open(t, dialect.SQLite,
		"file:ent?mode=memory&cache=shared&_fk=1",
		enttest.WithMigrateOptions(migrate.WithForeignKeys(true)),
	)
	return client
}

func TestCreateUser_Success(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx := testCtx()

	svcCtx := &svc.ServiceContext{DB: client}
	logger := logx.WithContext(ctx)

	logic := &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logger,
	}

	in := &apps.UserReq{
		Username: proto.String("testuser"),
		Password: proto.String("password123"),
		Nickname: proto.String("测试用户"),
		Mobile:   proto.String("13800138000"),
		Email:    proto.String("test@example.com"),
		Status:   proto.String("active"),
		Remark:   proto.String("测试用"),
	}

	resp, err := logic.CreateUser(in)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if resp.GetCode() != 200 {
		t.Fatalf("expected code 200, got %d", resp.GetCode())
	}
	if resp.GetData().GetId() <= 0 {
		t.Fatalf("expected positive id, got %d", resp.GetData().GetId())
	}
	if resp.GetData().GetIdStr() == "" {
		t.Fatal("expected non-empty idStr")
	}
	if resp.GetData().GetStatus() != "active" {
		t.Fatalf("expected status active, got %s", resp.GetData().GetStatus())
	}

	// Verify user was actually saved
	user, err := client.SysUser.Get(ctx, resp.GetData().GetId())
	if err != nil {
		t.Fatalf("failed to fetch created user: %v", err)
	}
	if user.Username != "testuser" {
		t.Fatalf("expected username testuser, got %s", user.Username)
	}
	if user.Nickname != "测试用户" {
		t.Fatalf("expected nickname 测试用户, got %s", user.Nickname)
	}
	if user.Mobile != "13800138000" {
		t.Fatalf("expected mobile 13800138000, got %s", user.Mobile)
	}
	if user.Email != "test@example.com" {
		t.Fatalf("expected email test@example.com, got %s", user.Email)
	}
	if user.Remark != "测试用" {
		t.Fatalf("expected remark 测试用, got %s", user.Remark)
	}
	if user.Status != sysuser.StatusActive {
		t.Fatalf("expected status active, got %s", user.Status)
	}

	// Verify audit fields were set by mixin hooks
	if user.TenantID != 1001 {
		t.Fatalf("expected tenantId 1001, got %d", user.TenantID)
	}
	if user.CreatedBy != "test-admin" {
		t.Fatalf("expected createdBy test-admin, got %s", user.CreatedBy)
	}
}

func TestCreateUser_WithDepartment(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx := testCtx()

	// First create a department
	dept, err := client.SysDept.Create().
		SetName("技术部").
		SetStatus("active").
		SetSort(1).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create dept: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: client}
	logic := &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}

	in := &apps.UserReq{
		Username: proto.String("deptuser"),
		Password: proto.String("pass123"),
		Nickname: proto.String("部门用户"),
		Status:   proto.String("active"),
		DeptId:   proto.Int64(dept.ID),
	}

	resp, err := logic.CreateUser(in)
	if err != nil {
		t.Fatalf("CreateUser with dept failed: %v", err)
	}

	user, err := client.SysUser.Get(ctx, resp.GetData().GetId())
	if err != nil {
		t.Fatalf("failed to fetch user: %v", err)
	}
	if user.DeptID != dept.ID {
		t.Fatalf("expected deptId %d, got %d", dept.ID, user.DeptID)
	}
}

func TestCreateUser_WithRoles(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx := testCtx()

	// First create roles
	role, err := client.SysRole.Create().
		SetName("管理员").
		SetCode("admin").
		SetSort(1).
		SetStatus(sysrole.StatusActive).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create role: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: client}
	logic := &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}

	in := &apps.UserReq{
		Username: proto.String("roleuser"),
		Password: proto.String("pass123"),
		Nickname: proto.String("角色用户"),
		Status:   proto.String("active"),
		RoleIds:  []int64{role.ID},
	}

	resp, err := logic.CreateUser(in)
	if err != nil {
		t.Fatalf("CreateUser with roles failed: %v", err)
	}

	user, err := client.SysUser.Query().
		Where(sysuser.IDEQ(resp.GetData().GetId())).
		WithRoles().
		Only(ctx)
	if err != nil {
		t.Fatalf("failed to fetch user with roles: %v", err)
	}
	if len(user.Edges.Roles) != 1 {
		t.Fatalf("expected 1 role, got %d", len(user.Edges.Roles))
	}
	if user.Edges.Roles[0].Code != "admin" {
		t.Fatalf("expected role code admin, got %s", user.Edges.Roles[0].Code)
	}
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx := testCtx()

	svcCtx := &svc.ServiceContext{DB: client}
	logic := &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}

	in := &apps.UserReq{
		Username: proto.String("duplicate"),
		Password: proto.String("pass123"),
		Nickname: proto.String("重复用户1"),
		Status:   proto.String("active"),
	}

	_, err := logic.CreateUser(in)
	if err != nil {
		t.Fatalf("first create should succeed: %v", err)
	}

	// Second create with same username should fail
	_, err = logic.CreateUser(in)
	if err == nil {
		t.Fatal("expected error for duplicate username, got nil")
	}
}

func TestCreateUser_EmptyUsername(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx := testCtx()
	svcCtx := &svc.ServiceContext{DB: client}
	logic := &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}

	in := &apps.UserReq{
		Username: proto.String(""),
		Password: proto.String("pass123"),
		Status:   proto.String("active"),
	}

	_, err := logic.CreateUser(in)
	if err == nil {
		t.Fatal("expected error for empty username, got nil")
	}
}
