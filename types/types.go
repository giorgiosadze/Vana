package types

type User struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"user_name"`
	Password string `json:"password"`
}

type Response struct {
	ResponseID    uint64 `json:"response_id"`
	QuestionID    uint64 `json:"question_id"`
	QuizID        uint64 `json:"quiz_id"`
	UserID        uint64 `json:"user_id"`
	CurrentAnswer string `json:"current_answer"`
}

type Quiz struct {
	QuizID    int64      `json:"quiz_id"`
	OwnerID   uint64     `json:"owner_id"`
	Questions []Question `json:"Questions"`
}

type Question struct {
	QuestionID    uint64   `json:"question_id"`
	QuestionText  string   `json:"question_text"`
	Answers       []Answer `json:"answers"`
	CorrectAnswer Answer   `json:"answer"`
}

type Answer struct {
	AnswerID   uint64 `json:"answer_id"`
	AnswerText string `json:"answer_text"`
}
