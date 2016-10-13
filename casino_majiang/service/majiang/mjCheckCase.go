package majiang

var CHECK_CASE_STATUS_CHECKING int32 = 0        //表示没有判定过
var CHECK_CASE_STATUS_CHECKED int32 = 1        //表示碰／杠 判定过


var CHECK_CASE_bean_STATUS_CHECKING int32 = 0
var CHECK_CASE_bean_STATUS_CHECKED int32 = 1      //已经check 了


//checkBean


func (c *CheckBean) IsChecked() bool {
	return false;
}

//checkCase
func (c *CheckCase) GetBuBean(checkStatus int32) *CheckBean {
	for _, b := range c.CheckB {
		if b != nil && b.GetCheckStatus() == checkStatus  && b.GetCanHu() {
			return b
		}
	}

	return nil
}


//修改判断事件的状态
func (c *CheckCase) UpdateChecStatus(status int32) error {
	*c.CheckStatus = status
	return nil
}

func (c *CheckCase) UpdateCheckBeanStatus(userId uint32, status int32) error {
	for _, bean := range c.CheckB {
		if bean != nil && bean.GetUserId() == userId {
			*bean.CheckStatus = status        //已经check过了
			return nil
		}
	}
	return nil
}

//是否已经验证完了...
func (c *CheckCase) IsChecked() bool {
	return c.GetCheckStatus() == CHECK_CASE_STATUS_CHECKED
}


