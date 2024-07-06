package utils

type AccessTokenResponse struct {
	Code              int    `json:"code"`
	Expire            int    `json:"expire"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	AppAccessToken    string `json:"app_access_token"`
	UserAccessToken   string `json:"user_access_token"`
}
type FlyStatus struct {
	IsActivated bool `json:"is_activated"`
	IsExited    bool `json:"is_exited"`
	IsFrozen    bool `json:"is_frozen"`
	IsResigned  bool `json:"is_resigned"`
	IsUnjoin    bool `json:"is_unjoin"`
}

type FlyUser struct {
	Mobile string    `json:"mobile"`
	Status FlyStatus `json:"status"`
	UserID string    `json:"user_id"`
}

type FlyUsersResponse struct {
	Code int `json:"code"`
	Data struct {
		UserList []FlyUser `json:"user_list"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type FlyUsersInfoResponse struct {
	Code int `json:"code"`
	Data struct {
		User struct {
			Name   string `json:"name"`
			UserID string `json:"user_id"`
		} `json:"user"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type FlyUserListResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		HasMore   bool   `json:"has_more"`
		PageToken string `json:"page_token"`
		Users     []struct {
			Avatar struct {
				Avatar72     string `json:"avatar_72"`
				Avatar240    string `json:"avatar_240"`
				Avatar640    string `json:"avatar_640"`
				AvatarOrigin string `json:"avatar_origin"`
			} `json:"avatar"`
			DepartmentIds []string `json:"department_ids"`
			Name          string   `json:"name"`
			OpenId        string   `json:"open_id"`
			UserId        string   `json:"user_id"`
		} `json:"users"`
	} `json:"data"`
}

type UserData struct {
	AccessToken      string `json:"access_token"`
	AvatarBig        string `json:"avatar_big"`
	AvatarMiddle     string `json:"avatar_middle"`
	AvatarThumb      string `json:"avatar_thumb"`
	AvatarUrl        string `json:"avatar_url"`
	EnName           string `json:"en_name"`
	ExpiresIn        int    `json:"expires_in"`
	Name             string `json:"name"`
	OpenID           string `json:"open_id"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	Sid              string `json:"sid"`
	TenantKey        string `json:"tenant_key"`
	TokenType        string `json:"token_type"`
	UnionID          string `json:"union_id"`
	DepartmentName   string `json:"department_name"`
}

type UserInfoRes struct {
	Code int      `json:"code"`
	Data UserData `json:"data"`
	Msg  string   `json:"msg"`
}

type FlyMessageRequest struct {
	ReceiveID string `json:"receive_id"`
	MsgType   string `json:"msg_type"`
	Content   string `json:"content"`
	UUID      string `json:"uuid"`
}

type UserTokenRes struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data UserTokenData `json:"data"`
}
type UserTokenData struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	Scope            string `json:"scope"`
}

type DepartmentIdRes struct {
	Code int `json:"code"`
	Data struct {
		User struct {
			DepartmentIDs []string `json:"department_ids"`
			LeaderUserID  string   `json:"leader_user_id"`
			MobileVisible bool     `json:"mobile_visible"`
			OpenID        string   `json:"open_id"`
			Orders        []struct {
				DepartmentID    string `json:"department_id"`
				DepartmentOrder int    `json:"department_order"`
				IsPrimaryDept   bool   `json:"is_primary_dept"`
				UserOrder       int    `json:"user_order"`
			} `json:"orders"`
			UnionID string `json:"union_id"`
		} `json:"user"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type DepartmentInfoRes struct {
	Code int `json:"code"`
	Data struct {
		Department struct {
			Name string `json:"name"`
		} `json:"department"`
	} `json:"data"`
	Msg string `json:"msg"`
}
