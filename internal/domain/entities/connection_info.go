package entities

type CertificateConnection struct {
	ID          uint   `gorm:"primaryKey"`
	Certificate []byte `gorm:"column:certificate;type:bytea;not null"`
	Key         []byte `gorm:"column:key;type:bytea;not null"`
	Policy      string `gorm:"column:policy;type:varchar(255)"`
	Mode        string `gorm:"column:mode;type:varchar(255)"`
}

type AnonymousConnection struct {
	ID     uint   `gorm:"primaryKey"`
	Policy string `gorm:"column:policy;type:varchar(255)"`
	Mode   string `gorm:"column:mode;type:varchar(255)"`
}

type PasswordConnection struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"column:username;type:varchar(255);not null"`
	Password string `gorm:"column:password;type:varchar(255);not null"`
	Policy   string `gorm:"column:policy;type:varchar(255)"`
	Mode     string `gorm:"column:mode;type:varchar(255)"`
}
