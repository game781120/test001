package utils

type JunShiReq struct {
	Query string `json:"query"`
	Model string `json:"model"`
}

type JunShiRes struct {
	IsPass string `json:"is_pass"`
}

type TCustom struct {
	KnowledgeIds []int64 `json:"knowledge_ids"`
}

type Safe struct {
	ComplianceLabel string `json:"compliance_label"`
	RiskLevel       string `json:"risk_level"`
	SensitiveWords  string `json:"sensitive_words"`
}

type SafeReq struct {
	Text      string `json:"text"`
	RequestID string `json:"request_id"`
	Model     string `json:"model"`
}
type SafeRes struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	RequestID   string `json:"request_id"`
	ResponseID  string `json:"response_id"`
	PassThrough string `json:"pass_through"`
	Model       string `json:"model"`
	Created     int64  `json:"created"`
	Safe
}
type NluReq struct {
	RequestId string `json:"request_id"`
	Question  string `json:"question"`
}
type NluRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Nulagent struct {
			SignIn struct {
				Name string `json:"name"`
			} `json:"signIn"`
			Reserve struct {
				Name string `json:"name"`
			} `json:"reserve"`
			FreeTalk struct {
				How string `json:"How"`
			} `json:"freeTalk"`
		} `json:"nulagent"`
		RequestID string `json:"request_id"`
	} `json:"data"`
}
