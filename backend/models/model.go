package models

type BaseModel struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	CreatedAt uint64 `gorm:"column:created_at;type:bigint unsigned;NULL;" json:"created_at"`
	UpdatedAt uint64 `gorm:"column:updated_at;type:bigint unsigned;NULL;" json:"updated_at"`
}

type TaskModel struct {
	Status      string `gorm:"column:status;type:text;check:status IN ('queue', 'in_progress', 'completed')" json:"status"`
	StartedAt   uint64 `gorm:"column:started_at;type:bigint unsigned;NULL;" json:"started_at"`
	CompletedAt uint64 `gorm:"column:completed_at;type:bigint unsigned;NULL;" json:"completed_at"`
}
