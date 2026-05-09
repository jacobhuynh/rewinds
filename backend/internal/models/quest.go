package models

import "time"

type Quest struct {
	ID               string    `json:"id"`
	Type             string    `json:"type"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	PtsRewardFree    int       `json:"pts_reward_free"`
	PtsRewardPremium int       `json:"pts_reward_premium"`
	Action           string    `json:"action"`
	TargetCount      int       `json:"target_count"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
}

type UserQuestProgress struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	QuestID     string    `json:"quest_id"`
	Progress    int       `json:"progress"`
	Completed   bool      `json:"completed"`
	PtsClaimed  bool      `json:"pts_claimed"`
	PeriodStart time.Time `json:"period_start"`
	CreatedAt   time.Time `json:"created_at"`
}
