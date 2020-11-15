//生成的文件建议不要改动,详见mysql-struct-maker.go ParseColumn方法源码生成格式 
package dbmodels 

type Characters struct {
	CharacterID int64 `sql:"CharacterID"` // 数据库注释:角色id 一个账号下有多个角色 
 	AcountID int64 `sql:"AcountID"` // 数据库注释:玩家账户id 
 	Job int8 `sql:"Job"` // 数据库注释:职业 
 	Gender int8 `sql:"Gender"` // 数据库注释:性别 
 	Lvl int32 `sql:"Lvl"` // 数据库注释:等级 
 	VipLvl int8 `sql:"VipLvl"` // 数据库注释:vip等级 
 }
