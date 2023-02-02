package models

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
