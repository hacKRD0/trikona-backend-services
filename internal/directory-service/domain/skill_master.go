// internal/directory-service/domain/skill_master.go
package domain

// SkillMaster holds the canonical list of skills
// +-----------+      +-----------------+      +--------+
// | students  |<---->| student_skill   |<---->| skill_master |
// +-----------+      +-----------------+      +--------+
// student_skill is the join table

type SkillMaster struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100;unique;not null" json:"name"`
}
