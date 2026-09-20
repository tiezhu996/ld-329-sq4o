package service

import (
	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
	"cyskillswap/internal/repository"
)

// defaultViewer 总览与旧版匹配接口默认以主页用户视角计算，
// 已套用同校区、信用分与未来一周连续时隙约束。
const defaultViewer = "林澈"

func Overview() model.Overview {
	skills := repository.ListSkills()
	needs := repository.ListNeeds()
	matches := Matches()
	appointments := repository.ListAppointments()
	reviews := repository.ListReviews()
	messages := repository.ListMessages()
	return model.Overview{
		Service:    constants.ServiceName,
		Categories: constants.SkillCategories,
		Metrics: map[string]int{
			"skills": len(skills), "needs": len(needs), "matches": len(matches),
			"appointments": len(appointments), "reviews": len(reviews), "unread": 3,
		},
		Skills: skills, Needs: needs, Matches: matches,
		Appointments: appointments, Reviews: reviews, Messages: messages,
		Profile: repository.GetProfile(),
	}
}

func Skills() []model.Skill             { return repository.ListSkills() }
func Needs() []model.Need               { return repository.ListNeeds() }
func Matches() []model.Match            { return MatchBoard(defaultViewer).Matches }
func Appointments() []model.Appointment { return repository.ListAppointments() }
func Reviews() []model.Review           { return repository.ListReviews() }
func Messages() []model.Conversation    { return repository.ListMessages() }
func Profile() model.Profile            { return repository.GetProfile() }
