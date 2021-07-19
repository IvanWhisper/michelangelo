package excode

const (
	UserErr          ExCode = 2000
	UserNoFound      ExCode = 2001
	UserRoleNoFound  ExCode = 2002
	UserLoginFailed  ExCode = 2003
	UserLogoutFailed ExCode = 2004
	UserHasNoFound   ExCode = 2005
	UserPwdError     ExCode = 2006

	UserCreateErr        ExCode = 2100 // 创建用户时，失败
	UserRoleCreateFailed ExCode = 2101 // 创建用户时，角色分配失败
	UserHasExist         ExCode = 2102 // 创建用户时，用户已经存在

	UserModifyFailed     ExCode = 2200
	UserRoleModifyFailed ExCode = 2201
	UserPwdModifyFailed  ExCode = 2202

	UserDeleteFailed     ExCode = 2300
	UserRoleDeleteFailed ExCode = 2301
)
