package generate

type IDRecord struct {
	Name     string `xorm:"varchar(32) pk 'Name'" json:"Name"`
	CurValue int64  `xorm:"bigint 'CurValue'" json:"CurValue"`
}

func (p *IDRecord) TableName() string {
	return "idrecord"
}
