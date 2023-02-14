package models

type Room struct {
	Id      int64  `xorm:"pk autoincr BIGINT(20)"`
	Name    string `xorm:"not null default '' comment('房间名称') VARCHAR(255)"`
	Onlines int    `xorm:"not null default 0 comment('在线人数') INT(11)"`
	Img     string `xorm:"not null default '' comment('封面图') VARCHAR(255)"`
}

type RoomUser struct {
	Id     int64 `xorm:"pk autoincr BIGINT(20)"`
	Uid    int64 `xorm:"not null default 0 comment('用户id') unique BIGINT(20)"`
	RoomId int64 `xorm:"not null default 0 comment('房间id') index BIGINT(20)"`
	Ctm    int64 `xorm:"not null default 0 comment('加入时间') BIGINT(20)"`
}

type User struct {
	Id       int64  `xorm:"pk autoincr BIGINT(20)"`
	UserName string `xorm:"not null default '' comment('用户名') unique VARCHAR(20)"`
	Token    string `xorm:"not null default '' comment('token') VARCHAR(32)"`
	Pwd      string `xorm:"not null default '' comment('密码') CHAR(32)"`
	Salt     string `xorm:"not null default '' comment('随机字符串') VARCHAR(20)"`
	NickName string `xorm:"not null default '' comment('昵称') VARCHAR(20)"`
	Avatar   string `xorm:"not null default '' comment('头像') VARCHAR(255)"`
	Ctm      int64  `xorm:"not null default 0 comment('注册时间') BIGINT(20)"`
	Ip       string `xorm:"not null default '' comment('注册ip') VARCHAR(40)"`
	Silent   int64  `xorm:"not null default 0 comment('禁言截止时间') BIGINT(20)"`
	Disable  int64  `xorm:"not null default 0 comment('封号截止时间') BIGINT(20)"`
	Online   int    `xorm:"not null default 0 comment('0 离线 1在线') INT(11)"`
	Ltm      int64  `xorm:"not null default 0 comment('最近登录时间') BIGINT(20)"`
	Lip      string `xorm:"not null default '0' comment('最近登录ip') VARCHAR(40)"`
}
