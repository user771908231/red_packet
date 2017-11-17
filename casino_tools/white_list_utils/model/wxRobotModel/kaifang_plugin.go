package wxRobotModel

import (
	"casino_common/common/log" // 导入日志包
	"github.com/songtianyi/wechat-go/wxweb"  // 导入协议包
	"fmt"
	"strings"
)

// Register plugin
// 必须有的插件注册函数
func Register(session *wxweb.Session) {
	// 将插件注册到session
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(kaifang), "kaifang")
	if err := session.HandlerRegister.EnableByName("kaifang"); err != nil {
		log.E("微信机器人注册kaifang handler失败。")
	}
}

// 消息处理函数
func kaifang(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	if msg.Content == "测试" {
		session.SendText("测试完成，机器人服务一切正常。", session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
	}

	// 可选:避免此插件对未保存到通讯录的群生效 可以用contact manager来过滤
	contact := session.Cm.GetContactByUserName(wxweb.RealTargetUserName(session, msg))
	if contact == nil {
		log.E("ignore the messages '%s' from %v", msg.Content, wxweb.RealTargetUserName(session, msg))
		return
	}
	if msg.IsGroup {
		log.T("讨论组[%s]:%s", contact.NickName, msg.Content)
		if strings.HasPrefix(msg.Content, "开房") {
			group_info := GetRroupInfoByName(contact.NickName)
			if group_info != nil {
				switch {
				case strings.HasPrefix(group_info.GameType, "红中"),
					strings.HasPrefix(group_info.GameType, "转转"):
					DoZzHzKaifang(group_info, session, msg, contact)
				case strings.HasPrefix(group_info.GameType, "经典跑得快"),
					strings.HasPrefix(group_info.GameType, "十五张跑得快"),
					strings.HasPrefix(group_info.GameType, "十六张跑得快"),
					strings.HasPrefix(group_info.GameType, "跑得快"):
					DoPdkKaifang(group_info, session, msg, contact)
				default:
					log.E("不支持的未知类型游戏。")
					session.SendText(fmt.Sprintf("开房失败，错误：暂不支持'%s'开房", group_info.GameType), session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
				}
			}else {
				session.SendText(fmt.Sprintf("开房失败，错误：当前微信群未注册"), session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
			}
		}
	}else {
		log.T("联系人[%s]:%s", contact.NickName, msg.Content)
	}
}

type CreateConfig [][][]string

//解析关键词
func (c *CreateConfig) GetKeywords(owner_opt,user_opt string) (res []string) {
	for _, paras := range *c{
		default_paras := ""
		has_paras := false
		//匹配用户选项
		for _, alias := range paras{
			for _, key := range alias {
				if strings.Contains(user_opt, key) {
					has_paras = true
					default_paras = alias[0]
					break
				}
			}
			if has_paras {
				break
			}
		}
		//匹配房主选项
		if has_paras == false {
			for paras_i, alias := range paras{
				for alias_i, key := range alias {
					//默认选项
					if paras_i == 0 && alias_i == 0 {
						default_paras = key
					}
					if strings.Contains(owner_opt, key) {
						has_paras = true
						default_paras = alias[0]
						break
					}
				}
				if has_paras {
					break
				}
			}
		}
		//插入
		res = append(res, default_paras)
	}
	return
}
