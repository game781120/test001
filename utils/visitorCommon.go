package utils

type ChatRequest struct {
	ChatId        string `json:"chatId"`
	MessagesId    string `json:"messagesId"`
	Message       string `json:"message"`
	IgnoreHistory bool   `json:"ignoreHistory"`
	Stream        bool   `json:"stream"`
}
type QaReq struct {
	KeyWords   string   `json:"keyWords"`
	DatasetIds []string `json:"datasetIds"`
	Top        int      `json:"top"`
}
type QaRes struct {
	Content          string  `json:"content"`
	Answer           string  `json:"answer"`
	Qa               int     `json:"qa"`
	ImageUrls        string  `json:"imageUrls"`
	VideoUrls        string  `json:"videoUrls"`
	RelatedQuestions string  `json:"relatedQuestions"`
	Scores           float32 `json:"scores"`
}

type Appointment struct {
	AppointmentId string `json:"appointmentId"`
	//客户信息
	CustomerMobile     string `json:"customerMobile"`
	CustomerName       string `json:"customerName"`
	CustomerCompany    string `json:"customerCompany"`
	CustomerDepartment string `json:"customerDepartment"`
	CustomerPosition   string `json:"customerPosition"`
	//预约信息
	VisitingDataTime    string `json:"visitingDatetime"`
	VisitingAddress     string `json:"visitingAddress"`
	VisitingDescription string `json:"visitingDescription"`
	VisitingTotals      int    `json:"visitingTotals"`
	//接待人信息
	ReceptionMobile     string `json:"receptionMobile"`
	ReceptionName       string `json:"receptionName"`
	UserOpenId          string `json:"userOpenId"`
	ReceptionCompany    string `json:"receptionCompany"`
	ReceptionDepartment string `json:"receptionDepartment"`
	ReceptionPosition   string `json:"receptionPosition"`
	IsSignIn            bool   `json:"isSignIn"`
	//签到状态 0 已预约  1 已签到
	SignInState    int    `json:"signInState"`
	SignInDataTime string `json:"signInDatetime"`
	//更新时间
	UpdateDatetime string `json:"updateDatetime"`
	CreateDatetime string `json:"createDatetime"`
	//预约人与接待人有可能不是同一个人
	AppointmentMobile     string `json:"appointmentMobile"`
	AppointmentName       string `json:"appointmentName"`
	AppointmentCompany    string `json:"appointmentCompany"`
	AppointmentDepartment string `json:"appointmentDepartment"`
	AppointmentPosition   string `json:"appointmentPosition"`
}

type AppointmentTemp struct {
	AppointmentId string `json:"appointmentId"`
	//客户信息
	CustomerMobile     string `json:"customerMobile"`
	CustomerName       string `json:"customerName"`
	CustomerCompany    string `json:"customerCompany"`
	CustomerDepartment string `json:"customerDepartment"`
	CustomerPosition   string `json:"customerPosition"`
	//预约信息
	VisitingDataTime    string `json:"visitingDatetime"`
	VisitingAddress     string `json:"visitingAddress"`
	VisitingDescription string `json:"visitingDescription"`
	VisitingTotals      string `json:"visitingTotals"`
	//接待人信息
	ReceptionMobile     string `json:"receptionMobile"`
	ReceptionName       string `json:"receptionName"`
	ReceptionCompany    string `json:"receptionCompany"`
	ReceptionDepartment string `json:"receptionDepartment"`
	ReceptionPosition   string `json:"receptionPosition"`
	//签到状态 0 已预约  1 已签到
	SignInState    string `json:"signInState"`
	SignInDataTime string `json:"signInDatetime"`
	//更新时间
	UpdateDatetime string `json:"updateDatetime"`
	CreateDatetime string `json:"createDatetime"`
	//预约人与接待人有可能不是同一个人
	AppointmentMobile     string `json:"appointmentMobile"`
	AppointmentName       string `json:"appointmentName"`
	AppointmentCompany    string `json:"appointmentCompany"`
	AppointmentDepartment string `json:"appointmentDepartment"`
	AppointmentPosition   string `json:"appointmentPosition"`
}

type AppointmentReq struct {
	Appointment Appointment `json:"appointment"`
}

type AppointmentRes struct {
	Code        int         `json:"code"`
	Msg         string      `json:"msg"`
	Appointment Appointment `json:"appointment"`
}

type CustomerSignInReq struct {
	AppointmentId string `json:"appointmentId"`
	Mobile        string `json:"mobile"`
	Name          string `json:"name"`
}

type CustomerSignInRes struct {
	Code        int         `json:"code"`
	Msg         string      `json:"msg"`
	Appointment Appointment `json:"appointment"`
}

type QueryRecordReq struct {
	Type            int    `json:"type"`
	KeyWords        string `json:"keyWords"`
	ReceptionName   string `json:"receptionName"`
	ReceptionMobile string `json:"receptionMobile"`
	UserOpenId      string `json:"userOpenId"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	SignInState     int    `json:"signInState"`
	Current         int    `json:"current"`
	PageSize        int    `json:"pageSize"`
}

type ExportRes struct {
	Code            int           `json:"code"`
	Msg             string        `json:"msg"`
	CurrentPage     int           `json:"currentPage"`
	Total           int           `json:"total"`
	AppointmentList []Appointment `json:"appointmentList"`
}

type DelRecordReq struct {
	AppointmentIdList []string `json:"appointmentList"`
}

type DelRecordRes struct {
	Code              int      `json:"code"`
	Msg               string   `json:"msg"`
	AppointmentIdList []string `json:"appointmentList"`
}

type QueryRecordRes struct {
	Code            int           `json:"code"`
	Msg             string        `json:"msg"`
	CurrentPage     int           `json:"currentPage"`
	Total           int           `json:"total"`
	AppointmentList []Appointment `json:"appointmentList"`
}

type ManagerData struct {
	MobileNumber string `json:"mobileNumber"`
	Name         string `json:"name"`
}

type ManagerRs struct {
	Code        int    `json:"code"`
	Msg         string `json:"msg"`
	Data        any    `json:"data"`
	CurrentPage int    `json:"currentPage"`
	Total       int    `json:"total"`
}
