package majiang

import "errors"

var CHECK_CASE_STATUS_CHECKING int32 = 0 //表示没有判定过
var CHECK_CASE_STATUS_CHECKED int32 = 1  //表示碰／杠 判定过

var CHECK_CASE_BEAN_STATUS_CHECKING int32 = 0
var CHECK_CASE_BEAN_STATUS_CHECKED int32 = 1 //已经check 了

//checkBean

func (c *CheckBean) IsChecked() bool {
	return c.GetCheckStatus() == CHECK_CASE_BEAN_STATUS_CHECKED;
}

func (c *CheckBean) IsChecking() bool {
	return c.GetCheckStatus() == CHECK_CASE_BEAN_STATUS_CHECKING
}

//checkCase
func (c *CheckCase) GetHuBean(checkStatus int32) *CheckBean {
	for _, b := range c.CheckB {
		if b != nil && b.GetCheckStatus() == checkStatus && b.GetCanHu() {
			return b
		}
	}

	return nil
}

func (c *CheckCase) GetBeanByUserIdAndStatus(userId uint32, status int32) *CheckBean {
	for _, b := range c.CheckB {
		if b != nil && b.GetCheckStatus() == status && b.GetUserId() == userId {
			return b
		}
	}
	return nil
}

//得到下一个需要判断的checkBean

func (c *CheckCase) GetNextBean() *CheckBean {
	if c == nil {
		return nil
	}

	//判断是否已经checked了
	if c.IsChecked() {
		return nil
	}

	//判断的流程，先问 只能乎，不能硼，杠的人，避免逻辑出问题
	var caseBean *CheckBean = nil
	for _, bean := range c.CheckB {
		if bean != nil && !bean.IsChecked() && bean.GetCanHu() && !bean.GetCanPeng() && !bean.GetCanGang() {
			caseBean = bean
			return caseBean
		}
	}

	//判断胡（可以碰或者可以杠）的牌
	for _, bean := range c.CheckB {
		if bean != nil && !bean.IsChecked() && bean.GetCanHu() {
			caseBean = bean
			return caseBean
		}
	}

	//如果这里的caseBean ！=nil 表示还有可以胡牌的人没有进行判定
	for _, bean := range c.CheckB {
		if bean != nil && !bean.IsChecked() && (bean.GetCanGang() || bean.GetCanPeng()) {
			caseBean = bean
			return caseBean

		}
	}

	//判断吃的牌
	for _, bean := range c.CheckB {
		if bean != nil && !bean.IsChecked() && bean.GetCanChi() {
			caseBean = bean
			return caseBean
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
			*bean.CheckStatus = status //已经check过了
			return nil
		}
	}
	return nil
}

//是否已经验证完了...
func (c *CheckCase) IsChecked() bool {
	return c.GetCheckStatus() == CHECK_CASE_STATUS_CHECKED
}
