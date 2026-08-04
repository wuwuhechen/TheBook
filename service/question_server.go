package service

import "TheBook/model"

// QuestionServer 协调题库存储、练习状态和 HTTP 处理器。
type QuestionServer struct {
	DB model.QuestionManager

	PM map[int]*model.Practice

	RS map[int]*model.RandomSession
}
