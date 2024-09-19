package models

type RolesModels struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type PermissionModels struct {
	ID        int64  `json:"id"`
	Groups    string `json:"groups"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type RolePermissionModels struct {
	ID           int64  `json:"id"`
	RoleID       int64  `json:"roleID"`
	PermissionID int64  `json:"permissionID"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type UserPermissionModels struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"userID"`
	PermissionID int64  `json:"permissionID"`
	Status       bool   `json:"status"`
	GrantedBy    int64  `json:"grantedBy"`
	GrantedAt    string `json:"grantedAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type UserRolePermissionModels struct {
	ID     int64  `json:"id"`
	Group  string `json:"group"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type AssignRoleToUserRequest struct {
	UserID int64 `json:"userID"`
	RoleID int64 `json:"roleID"`
}

type RoleCreateRequest struct {
	Name      string `json:"name"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type PermissionCreateRequest struct {
	Groups    string `json:"groups"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
