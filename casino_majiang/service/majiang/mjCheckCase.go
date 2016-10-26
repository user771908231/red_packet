package majiang

import "errors"

var CHECK_CASE_STATUS_CHECKING int32 = 0        //表示没有判定过
var CHECK_CASE_STATUS_CHECKED int32 = 1        //表示碰／杠 判定过
var CHECK_CASE_STATUS_CHECKING_HUED int32 = 2        //已经有人胡了


var CHECK_CASE_BEAN_STATUS_CHECKING int32 = 0
var CHECK_CASE_BEAN_STATUS_CHECKED int32 = 1  //已经check 了
var CHECK_CASE_BEAN_STATUS_PASS int32 = 2     //已经check 了


//checkBean


func (c *CheckBean) IsChecked() bool {
	return c.GetCheckStatus() == CHECK_CASE_BEAN_STATUS_CHECKED;
}

func (c *CheckBean) IsPassed() bool {
	return c.GetCheckStatus() == CHECK_CASE_BEAN_STATUS_PASS;
}

func (c *CheckBean) IsChecking() bool {
	return c.GetCheckStatus() == CHECK_CASE_BEAN_STATUS_CHECKING
}

//checkCase
func (c *CheckCase) GetHuBean(checkStatus int32) *CheckBean {
	for _, b := range c.CheckB {
		if b != nil && b.GetCheckStatus() == checkStatus  && b.GetCanHu() {
			return b
		}
	}

	return nil
}

func (c *CheckCase) GetBeanByUserIdAndStatus(userId uint32, status int32) *CheckBean {
	for _, b := range c.CheckB {
		if b != nil && b.GetCheckStatus() == status  && b.GetUserId() == userId {
			return b
		}
	}
	return nil
}

//得到下一个需要判断的checkBean
func (c *CheckCase) GetNextBean() *CheckBean {
	//判断是否已经checked了
	if c.IsChecked() {
		return nil
	}

	var caseBean *CheckBean = nil
	for _, bean := range c.CheckB {
		if bean != nil && !bean.IsChecked() && !bean.IsPassed() && bean.GetCanHu() {
			caseBean = bean
			break
		}
	}

	//如果之前已经有人胡过了，那么不能再碰或者杠牌了
	if c.GetCheckStatus() == CHECK_CASE_STATUS_CHECKING_HUED {
		return caseBean
	}

	//如果这里的caseBean ！=nil 表示还有可以胡牌的人没有进行判定
	if caseBean == nil {
		for _, bean := range c.CheckB {
			if bean != nil && !bean.IsChecked() && !bean.IsPassed()  && !bean.GetCanHu() {
				caseBean = bean
				break
			}
		}
	}

	return caseBean
}

//修改判断事件的状态
func (c *CheckCase) UpdateChecStatus(status int32) error {
	*c.CheckStatus = status
	return nil
}

func (c *CheckCase) UpdateCheckBeanStatus(userId uint32, status int32) error {
	if c.CheckB == nil {
		return errors.New("不能设置状态，因为checkB 为nil")
	}
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


