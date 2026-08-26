// model包是各种结构的外部接口，包含结构体定义、管理器接口以及数据库实现和他们的构造函数。
package model

import (
	"TheBook/model/bank"
	"TheBook/model/manager"
	"TheBook/model/structs"
)

// 结构体定义
type Question = structs.Question
type Practice = structs.Practice
type QuestionResponse = structs.QuestionResponse
type PracticeResponse = structs.PracticeResponse
type Request = structs.Request
type PracticeRecord = structs.PracticeRecord
type AnswerRecord = structs.AnswerRecord
type QuestionProgress = structs.QuestionProgress
type RandomSession = structs.RandomSession
type User = structs.User
type RegisterRequest = structs.RegisterRequest
type LoginRequest = structs.LoginRequest
type QuestionPageData = structs.QuestionPageData
type PracticeResultItem = structs.PracticeResultItem
type PracticeResultPageData = structs.PracticeResultPageData

// 管理器接口
type QuestionManager = manager.QuestionManager
type PracticeManager = manager.PracticeManager
type QuestionProgressManager = manager.QuestionProgressManager
type RandomSessionManager = manager.RandomSessionManager
type UserManager = manager.UserManager

// 数据库实现
type QuestionBank = bank.QuestionBank
type SQLiteQuestionBank = bank.SQLiteQuestionBank
type PracticeBank = bank.PracticeBank
type QuestionProgressBank = bank.QuestionProgressBank
type RandomSessionBank = bank.RandomSessionBank
type UserBank = bank.UserBank
type UserBankSQLite = bank.UserBankSQLite

// 构造函数
var NewResponse = structs.NewResponse
var NewQuestionData = structs.NewQuestionData
var NewQuestionProgress = structs.NewQuestionProgress
var NewRandomSession = structs.NewRandomSession
var NewQuestionBank = bank.NewQuestionBank
var NewSQLiteQuestionBank = bank.NewSQLiteQuestionBank
var NewPracticeBank = bank.NewPracticeBank
var NewQuestionProgressBank = bank.NewQuestionProgressBank
var NewRandomSessionBank = bank.NewRandomSessionBank
var NewUserBank = bank.NewUserBank
var NewUserBankSQLite = bank.NewUserBankSQLite
