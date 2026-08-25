package bank

import "TheBook/model/structs"

// Local aliases keep bank implementations readable while preserving the
// dependency direction from banks to domain structures.
type Question = structs.Question
type Practice = structs.Practice
type QuestionProgress = structs.QuestionProgress
type RandomSession = structs.RandomSession
type User = structs.User
type RegisterRequest = structs.RegisterRequest
type LoginRequest = structs.LoginRequest

var NewQuestionProgress = structs.NewQuestionProgress
