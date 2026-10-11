package sqlite

import (
	"TheBook/auth"
	"os"
	"path/filepath"
	"testing"
)

// writeFixture 在临时目录中写入 JSON 测试夹具。
func writeFixture(t *testing.T, directory, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(directory, name), []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

// TestReadQuestionsRejectsQuestionWithTooFewChoices 验证题目选项不足会被拒绝。
func TestReadQuestionsRejectsQuestionWithTooFewChoices(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "questions.json")
	writeFixture(t, directory, "questions.json", `[
  {"id":1,"category":"Go","question":"Which?","choices":["A","B"],"answer":0,"explanation":"test"}
]`)

	if _, err := ReadQuestions(path); err == nil {
		t.Fatal("ReadQuestions() error = nil, want an error for fewer than four choices")
	}
}

// TestMigrateDatabaseImportsAndUpsertsAllData 验证迁移导入并更新全部数据类型。
func TestMigrateDatabaseImportsAndUpsertsAllData(t *testing.T) {
	directory := t.TempDir()
	writeFixture(t, directory, "data.json", `[
  {"id":7,"category":"Go","question":"Which keyword declares a variable?","choices":["var","func","package","import"],"answer":0,"explanation":"var declares a variable"}
]`)
	writeFixture(t, directory, "users.json", `[
  {"user_id":0,"username":"admin","password":"plain-password","nickname":"Administrator","role":"admin"}
]`)
	writeFixture(t, directory, "practice_records.json", `[
  {"id":21,"user_id":0,"practice_id":101,"total_questions":1,"correct_count":1,"wrong_count":0,"start_time":"2026-01-02T03:04:05Z","submit_time":"2026-01-02T03:09:05Z","answers":[{"question_id":7,"answer":0,"answered":true,"correct":true}]}
]`)
	writeFixture(t, directory, "question_progress.json", `[
  {"user_id":0,"current_question_id":7,"updated_at":"2026-01-02T03:09:05Z"}
]`)
	writeFixture(t, directory, "wrong_questions.json", `[
  {"user_id":0,"question_id":7,"wrong_count":2,"first_wrong_at":"2026-01-01T03:04:05Z","last_wrong_at":"2026-01-02T03:04:05Z"}
]`)

	db, err := OpenDatabase(filepath.Join(directory, "thebook.db"))
	if err != nil {
		t.Fatalf("OpenDatabase() error = %v", err)
	}
	closeTestDatabase(t, db)
	if err := CreateDatabase(db, []any{
		&Question{}, &User{}, &PracticeRecord{}, &PracticeAnswer{}, &QuestionProgress{}, &WrongQuestion{},
	}); err != nil {
		t.Fatalf("CreateDatabase() error = %v", err)
	}

	if err := MigrateDatabase(db, directory); err != nil {
		t.Fatalf("MigrateDatabase() error = %v", err)
	}

	var question Question
	if err := db.First(&question, 7).Error; err != nil {
		t.Fatalf("read imported question: %v", err)
	}
	if question.OptionA != "var" || question.Answer != 0 {
		t.Errorf("imported question = %#v, want option A %q and answer 0", question, "var")
	}

	var user User
	if err := db.Where("user_id = ?", 0).First(&user).Error; err != nil {
		t.Fatalf("read imported user with ID 0: %v", err)
	}
	if !auth.VerifyPassword(user.Password, "plain-password") {
		t.Error("imported plaintext password was not converted to a valid bcrypt hash")
	}

	var record PracticeRecord
	if err := db.First(&record, 21).Error; err != nil {
		t.Fatalf("read imported practice record: %v", err)
	}
	if record.UserID != 0 || record.PracticeID != 101 || record.CorrectCount != 1 {
		t.Errorf("imported record = %#v, want user 0, practice 101, correct count 1", record)
	}

	var answer PracticeAnswer
	if err := db.Where("practice_record_id = ? AND question_id = ?", 21, 7).First(&answer).Error; err != nil {
		t.Fatalf("read imported practice answer: %v", err)
	}
	if !answer.Answered || !answer.Correct || answer.Answer != 0 {
		t.Errorf("imported answer = %#v, want answered and correct answer 0", answer)
	}

	var progress QuestionProgress
	if err := db.Where("user_id = ?", 0).First(&progress).Error; err != nil {
		t.Fatalf("read imported progress: %v", err)
	}
	if progress.CurrentQuestionID != 7 {
		t.Errorf("progress current question = %d, want 7", progress.CurrentQuestionID)
	}

	var wrong WrongQuestion
	if err := db.Where("user_id = ? AND question_id = ?", 0, 7).First(&wrong).Error; err != nil {
		t.Fatalf("read imported wrong question: %v", err)
	}
	if wrong.WrongCount != 2 {
		t.Errorf("wrong count = %d, want 2", wrong.WrongCount)
	}

	if err := MigrateDatabase(db, directory); err != nil {
		t.Fatalf("second MigrateDatabase() error = %v", err)
	}
	for _, check := range []struct {
		model any
		want  int64
	}{
		{&Question{}, 1},
		{&User{}, 1},
		{&PracticeRecord{}, 1},
		{&PracticeAnswer{}, 1},
		{&QuestionProgress{}, 1},
		{&WrongQuestion{}, 1},
	} {
		var count int64
		if err := db.Model(check.model).Count(&count).Error; err != nil {
			t.Fatalf("count %T: %v", check.model, err)
		}
		if count != check.want {
			t.Errorf("count %T = %d, want %d after repeated migration", check.model, count, check.want)
		}
	}
}

// TestMigrateDatabaseAllowsMissingOptionalJSONFiles 验证可选 JSON 缺失时仍可迁移。
func TestMigrateDatabaseAllowsMissingOptionalJSONFiles(t *testing.T) {
	directory := t.TempDir()
	db, err := OpenDatabase(filepath.Join(directory, "empty.db"))
	if err != nil {
		t.Fatalf("OpenDatabase() error = %v", err)
	}
	closeTestDatabase(t, db)
	if err := CreateDatabase(db, []any{&Question{}, &User{}}); err != nil {
		t.Fatalf("CreateDatabase() error = %v", err)
	}

	if err := MigrateDatabase(db, directory); err != nil {
		t.Fatalf("MigrateDatabase() error = %v, want nil when JSON source files are absent", err)
	}
}
