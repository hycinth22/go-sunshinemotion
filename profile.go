package sunshinemotion

import "time"

type UserInfo struct {
	InClassID      int64  `json:"inClassID"`      // 班级ID
	InClassName    string `json:"inClassName"`    // 班级名称
	InCollegeID    int64  `json:"inCollegeID"`    // 院系ID
	InCollegeName  string `json:"inCollegeName"`  // 院系名称
	InSchoolID     int64  `json:"inSchoolID"`     // 学校ID
	InSchoolName   string `json:"inSchoolName"`   // 学校名称
	InSchoolNumber string `json:"inSchoolNumber"` // 学校编号
	IsTeacher      int    `json:"isTeacher"`      // 是否为教师
	NickName       string `json:"nickName"`       // 昵称
	PhoneNumber    string `json:"phoneNumber"`    // 电话号码
	Sex            string `json:"sex"`            //性别
	StudentName    string `json:"studentName"`    // 学生姓名
	StudentNumber  string `json:"studentNumber"`  // 学生编号
	UserRoleID     int    `json:"userRoleID"`
}

type UserSportResult struct {
	UserID            uint      // 用户ID
	Year              int       // 年度
	Term              string    // 学期
	RunDistance       float64   // 已跑距离
	QualifiedDistance float64   // 达标距离
	LastTime          time.Time // 上次跑步时间
}
