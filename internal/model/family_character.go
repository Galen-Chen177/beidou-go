package model

// FamilyCharacter 实体映射
type FamilyCharacter struct {
	Cid             int    `gorm:"primaryKey;column:cid" json:"cid,omitempty"`
	Familyid        *int   `json:"familyid,omitempty"`
	Seniorid        *int   `json:"seniorid,omitempty"`
	Reputation      *int   `json:"reputation,omitempty"`
	Todaysrep       *int   `json:"todaysrep,omitempty"`
	Totalreputation *int   `json:"totalreputation,omitempty"`
	Reptosenior     *int   `json:"reptosenior,omitempty"`
	Precepts        string `gorm:"column:precepts;type:varchar(200)" json:"precepts,omitempty"`
	Lastresettime   int64  `json:"lastresettime,omitempty"`
}

func (FamilyCharacter) TableName() string {
	return "family_character"
}
