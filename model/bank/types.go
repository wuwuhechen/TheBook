package bank

import "TheBook/model/structs"

// 使用别名可以保持银行实现的可读性，同时保留从银行到领域结构的依赖方向。
type Question = structs.Question
type Practice = structs.Practice
type QuestionProgress = structs.QuestionProgress
type RandomSession = structs.RandomSession
type User = structs.User
type RegisterRequest = structs.RegisterRequest
type LoginRequest = structs.LoginRequest
type PracticeResponse = structs.PracticeResponse
type PracticeRecord = structs.PracticeRecord

var NewQuestionProgress = structs.NewQuestionProgress
