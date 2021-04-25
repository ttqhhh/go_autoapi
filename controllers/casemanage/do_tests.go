package controllers

type CheckOut struct {
	Url   string                            `json:"url"`
	Uuid  string                            `json:"uuid"`
	Param map[string]interface{}            `json:"param"`
	Check map[string]map[string]interface{} `json:"check_point"`
}

//// 准备测试
//func (c *CaseManageController) performTests() {
//	caseID := c.GetString("case_id")
//	sep := ","
//	caseIdList := strings.Split(caseID, sep)
//	acm := models.TestCaseMongo{}
//	for _, i := range caseIdList{
//		id,_ := strconv.ParseInt(i, 10, 64)
//		tc := acm.GetOneCase(id)
//		libs.DoRequestV2(tc.ApiUrl, tc.Parameter, tc.Checkpoint)
//	}
//}
