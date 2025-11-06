package animals

type Revenue struct {
	gorm.Model
	UserID uint `gorm:"index"`
	Source string `gorm:"not null"`
	SourceID uint
	Type string `gorm:"not null"`
	Quantity float64
	UnitPrice float64
	Total float64 `gorm:"not null"`
	Date time.Time `gorm:"not null"`
	Notes string
}
