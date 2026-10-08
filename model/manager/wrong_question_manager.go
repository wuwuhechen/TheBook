package manager

import (
	"TheBook/model/structs"
	"time"
)

// WrongQuestionManager 管理用户错题本的累计记录与分页查询。
type WrongQuestionManager interface {
	// RecordWrongQuestion 记录一次答错；同一用户再次答错同一题时累计次数。
	RecordWrongQuestion(userID, questionID uint, wrongAt time.Time) error

	ListWrongQuestionsByUser(userID uint, limit, offset int) ([]*structs.WrongQuestion, error)

	// ListWrongQuestionsByIDs 查询用户错题本中指定的题目，用于重新组卷前的归属校验。
	ListWrongQuestionsByIDs(userID uint, questionIDs []uint) ([]*structs.WrongQuestion, error)
}
