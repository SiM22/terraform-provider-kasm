package client

type SessionPermission struct {
	UserID      string `json:"user_id"`
	Access      string `json:"access"`
	VNCUsername string `json:"vnc_username"`
	Username    string `json:"username"`
}

type SetSessionPermissionsRequest struct {
	TargetSessionPermissions TargetSessionPermissions `json:"target_session_permissions"`
}

type GetSessionPermissionsRequest struct {
	TargetSessionPermissions struct {
		KasmID string `json:"kasm_id"`
	} `json:"target_session_permissions"`
}

type DeleteAllSessionPermissionsRequest struct {
	TargetSessionPermissions struct {
		KasmID string `json:"kasm_id"`
	} `json:"target_session_permissions"`
}

type TargetSessionPermissions struct {
	KasmID             string                    `json:"kasm_id"`
	Access             string                    `json:"access,omitempty"`
	SessionPermissions []SessionPermissionAccess `json:"session_permissions,omitempty"`
}

type SessionPermissionAccess struct {
	UserID string `json:"user_id"`
	Access string `json:"access"`
}

type SessionPermissionResponse struct {
	SessionPermissions []SessionPermission `json:"session_permissions"`
}
