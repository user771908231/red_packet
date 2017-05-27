function postJson(url,data,cb)
{
    $.ajax({
        method:"post",url: url
        ,contentType:"application/json; charset=UTF-8"
        ,data:JSON.stringify(data)
        ,success:cb
    });
}

var gameType={
    git:{
        loginImg:"images/avatar.png",
        name:"皮皮麻将会员管理后台",
        agentname:"皮皮麻将会员代理后台",
        color:"background-color:white",
        addDiamondName:"添加钻石",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.3, recommend:0, recharge:0},
            firstLevel:{sca:0.4, recommend:4, recharge:1000},
            secondLevel:{sca:0.5, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.3, recommend:0, recharge:0},
            firstLevel:{sca:0.4, recommend:4, recharge:1000},
            secondLevel:{sca:0.5, recommend:10, recharge:1000}
        },
        dicStatistice:{

        }

    },
    ljmj:{
        loginImg:"images/ljmj.png",
        name:"辽宁麻将会员管理后台",
        agentname:"辽宁麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加辽宁麻将元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:{

        }
    },
    guandan:{
        loginImg:"images/guandan.png",
        name:"掼蛋麻将会员管理后台",
        agentname:"掼蛋麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加掼蛋麻将元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:{

        }
    },
    hbmj:{
        loginImg:"images/hbmj.png",
        name:"河北麻将会员管理后台",
        agentname:"河北麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加河北麻将元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:{

        }
    },
    gsmj:{
        loginImg:"images/gsmj.png",
        name:"甘肃麻将会员管理后台",
        agentname:"甘肃麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加甘肃麻将元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:
            {
                tdh:"推倒胡",
                xa:"西安",
                bj:"宝鸡",
                yl:"榆林",
                hs:"滑水",
                lz:"兰州",
                4:"4局",
                8:"8局",
                c0:"不吃",
                c1:"吃",
                f0:"不带风",
                f1:"带风",
                p0:"不吃胡",
                p1:"吃胡",
                lz0:"无勒子",
                lz1:"有勒子"
            }

    },
    nxmj:{
        loginImg:"images/nxmj.png",
        name:"宁夏麻将会员管理后台",
        agentname:"宁夏麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加宁夏麻将元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:{
            fanghu:"可以放胡",
            jiang258:"258做将",
            nd1:"可胡七对",
            Nod1:"不可胡七对",
            zhong:"红中赖子",
            Nozhong:"不带红中赖子",
            wind:"带风牌",
            Nowind:"不带风牌",
            sixCards:"甩六张",
            NosixCards:"不甩六张",
            NoFish:"不下鱼",
            twoFish:"下2鱼",
            fiveFish:"下5鱼",
            eightFish:"下8鱼"
        }

    },
    shmj:{
        loginImg:"images/shmj.png",
        name:"上海麻将会员管理后台",
        agentname:"上海麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加上海麻将元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:{

        }
    },
    gxphz:{
        loginImg:"images/gxphz.png",
        name:"广西跑胡子会员管理后台",
        agentname:"广西跑胡子会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加广西跑胡子元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice: {
            t0: "桂林",
            t1: "河池",
            r6: "6局",
            r10: "10局",
            r15: "15局",
            r16: "16局",
            r20: "20局",
            r106: "6局",
            r110: "10局",
            r115: "15局",
            r116: "16局",
            r120: "20局",
        }

    },
    henmj:{
        loginImg:"images/henmj.png",
        name:"河南麻将会员管理后台",
        agentname:"河南麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加河南麻将元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:{
            zz:"郑州",
            zz2:"郑州2人",
            zz3:"郑州3人",
            kf:"开封",
            ly:"洛阳",
            jz:"焦作",
            xy:"信阳",


            j01:"1局",
            j04:"4局",
            j06:"6局",
            j08:"8局",
            j10:"10局",
            j15:"15局",
            j16:"16局",
            j20:"20局",

            f0:"不带风",
            f1:"带风",
            dp:"点炮胡",
            zm:"自摸胡",
            hun:"带混",
            zfb:"庄翻倍",
            bc0:"不包次",
            bc1:"包次",
            zui7:"七公嘴",
            zui10:"十公嘴"
        }
    },
    ahmj:{
        loginImg:"images/ahmj.png",
        name:"安徽麻将会员管理后台",
        agentname:"安徽麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加安徽麻将元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:{

        }
    },
    xynmmj:{
        loginImg:"images/xynmmj.png",
        name:"内蒙麻将会员管理后台",
        agentname:"内蒙麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加内蒙麻将元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:{

        }
    },
    ynmj:{
        loginImg:"images/ynmj.png",
        name:"云南麻将会员管理后台",
        agentname:"云南麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加云南麻将元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:{

        }
    },
    hanmj:{
        loginImg:"images/hanmj.png",
        name:"海南麻将会员管理后台",
        agentname:"海南麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加海南麻将元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:
            {
                j4:"4局",
                j8:"8局",
                z:"推到胡",
                s:"海南",
                c0:"不可以吃",
                c1:"可以吃",
                f0:"不带风",
                f1:"带风",
                p0:"不可吃胡",
                p1:"可以吃胡",
                r0:"不带红中",
                r1:"红中赖子",
                q0:"无七对",
                q1:"可胡七对",
                g0:"不叫",
                g1:"叫嘎",
                w2:"2人玩",
                w3:"3人玩",
                w4:"4人玩"
            }
    },
    fjmj:{
        loginImg:"images/fjmj.png",
        name:"福建麻将会员管理后台",
        agentname:"福建麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加福建麻将元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:
            {
                qz:"泉州麻将",
                mq:"闽清麻将",
                zz:"漳州麻将",
                hz:"红中麻将",
                xm:"厦门麻将",
                nd:"宁德麻将",
                pt:"莆田麻将",
                4:"4局",
                8:"8局",
                TG:"托管模式",
                p4:"4人房间",
                p3:"3人房间",
                p2:"2人房间"
            }

    },
    jsmj:{
        loginImg:"images/jsmj.png",
        name:"江苏麻将会员管理后台",
        agentname:"江苏麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加江苏麻将元宝",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:
            {
                jsmj:"江苏麻将",
                nj:"南京",
                yz:"扬州",
                xh:"兴化",
                sz:"苏州",
                yc:"盐城",
                j01:"1局",
                j04:"4局",
                j06:"6局",
                j08:"8局",
                j10:"10局",
                j15:"15局",
                j16:"16局",
                j20:"20局",
                js100:"上限100",
                js200:"上限200",
                js500:"上限500",
                js999:"上限不限",
                qh9:"9番起胡",
                qh13:"13番起胡",
                qh18:"18番起胡",
                m3c4:"3摸4冲",
                m2c3:"2摸3冲",
                qd1:"可胡七对",
                qd2:"不可胡七对"
            }
    },
    phz:{
        loginImg:"images/phz.png",
        name:"跑胡子会员管理后台",
        agentname:"跑胡子会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加跑胡子元宝",
        jump:{
            tip:{url:"tip", name:"选择游戏后台"},
            phz:{url:"http://pdk.coolgamebox.com:88/login.html", name:"跑得快会员管理后台"}
        },
        gameIds: ["phz", "pdk"],//会员管理页面 权限添加
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:
            {
                t1: "邵阳",
                t2: "娄底",
                t3: "剥皮",
                t4: "怀化",
                t5: "常德",
                t6: "衡阳",
                t7: "长沙",
                t8: "湘乡",
                t9: "永州",
                t10: "攸县",
                t11: "岳阳",
                t12: "耒阳",
                t13: "郴州",

                r6: "6局",
                r10: "10局",
                r16: "16局",
                r1000: "",
                r1001: "",
                r1002: "",
                r106: "6局",
                r110: "10局",
                r116: "16局",
                r206: "6局",
                r210: "10局",
                r216: "16局",
            }
    },
    symj:{
        loginImg:"images/avatar.png",
        name:"湖南麻将会员管理后台",
        agentname:"湖南麻将会员代理后台",
        color:"background-color:white",
        addDiamondName:"添加湖南麻将钻石",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:
            {

            }
    },
    scmj:{
        loginImg:"images/scmj.png",
        name:"四川麻将会员管理后台",
        agentname:"四川麻将会员代理后台",
        color:"background-color:#F0F8FF",
        addDiamondName:"添加四川麻将钻石",
        jump:{//http://ddz.coolgamebox.com:88/login.html
            tip:{url:"tip", name:"选择游戏后台"},
            ddz:{url:"http://ddz.coolgamebox.com:88/login.html", name:"斗地主会员管理后台"}
        },
        gameIds: ["scmj", "ddz"],
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.3, recommend:0, recharge:0},
            firstLevel:{sca:0.4, recommend:4, recharge:1000},
            secondLevel:{sca:0.5, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.3, recommend:0, recharge:0},
            firstLevel:{sca:0.4, recommend:4, recharge:1000},
            secondLevel:{sca:0.5, recommend:10, recharge:1000}
        },
        dicStatistice:{
            z:"血战到底",
            n:"内江麻将",
            e:"二人玩法",
            d:"倒倒胡",
            b:"血流成河",
            y:"德阳麻将",
            s:"三人麻将",
            j4:"4局",
            j8:"8局",
            a0:"不换三张",
            a1:"换三张"
        }
    },
    jxmj:{
        loginImg:"images/jxmj.png",
        name:"江西麻将会员管理后台",
        agentname:"江西麻将会员代理后台",
        color:"background-color:#F5F5DC",
        addDiamondName:"添加江西麻将钻石",
        alipayTextTip:"请在浏览器里面操作充值",
        rechargeOpen:true,
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:
            {
                "area0" : "转转麻将",
                "area1" : "南昌麻将",
                "area2" : "上饶麻将",
                "area3" : "贵溪麻将",
                "round4" : "四局",
                "round8" : "八局",
                "round12" : "十二局",
                "playnum2" : "两人",
                "playnum3" : "三人",
                "playnum4" : "四人",
                "nanchangType0" : "翻精自摸",
                "nanchangType1" : "翻精点炮",
                "nanchangType2" : "无下精",
                "nanchangType3" : "埋地雷",
                "nanchangType4" : "回头一笑",
                "nanchangType5" : "同一首歌",
                "dianPaoSanFu0" : "点炮一家付",
                "dianPaoSanFu1" : "点炮多家付",
                "mingGangSanFu0" : "明杠一家付",
                "mingGangSanFu1" : "明杠多家付",
                "canEatHu0" : "自摸",
                "canEatHu1" : "点炮",
                "isHongZhong" : "红中癞子",
                "isJiang258" : "258将",
                "withWind" : "带风",
                "zhuaniaoType0" : "抓2鸟",
                "zhuaniaoType1" : "抓4鸟",
                "zhuaniaoType2" : "抓6鸟",
                "zhuaniaoType3" : "窝窝鸟",
                "zhuaniaoType4" : "上中下鸟"
            }
    },
    sxmj:{
        loginImg:"images/sxmj.png",
        name:"陕西麻将会员管理后台",
        agentname:"陕西麻将会员代理后台",
        color:"background-color:#F0F8FF",
        addDiamondName:"添加陕西麻将钻石",
        alipayTextTip:"请在浏览器里面操作充值",
        rechargeOpen:true,
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:
            {
                tdh:"推倒胡",
                xa:"西安",
                bj:"宝鸡",
                yl:"榆林",
                hs:"滑水",
                ls:"两色",
                wn:"渭南",
                gz:"锅子",

                j04:"4局",
                j08:"8局",

                c0:"不吃",
                c1:"吃",

                f0:"不带风",
                f1:"带风",

                p0:"不带吃胡（只炸不胡）",
                p1:"吃胡",

                lz0:"无赖子",
                lz1:"有赖子"
            }
    },
    gzmj:{
        loginImg:"images/gzmj.png",
        name:"贵州麻将会员管理后台",
        agentname:"贵州麻将会员代理后台",
        color:"background-color:#F0F8FF",
        addDiamondName:"添加贵州麻将钻石",
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.2, recommend:0, recharge:0},
            firstLevel:{sca:0.3, recommend:4, recharge:1000},
            secondLevel:{sca:0.4, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.2, recommend:0, recharge:0},
            firstLevel:{sca:0.3, recommend:4, recharge:1000},
            secondLevel:{sca:0.4, recommend:10, recharge:1000}
        },
        dicStatistice:
            {
                gy:"贵阳",
                zy:"遵义",
                sdg:"三丁拐",
                tr:"铜仁",
                ash:"安顺",
                edg:"二丁拐",
                4:"4局",
                8:"8局",
                12:"12局"
            }

    },
    pdk:{
        loginImg:"images/pdk.png",
        name:"跑得快会员管理后台",
        agentname:"跑得快会员代理后台",
        color:"background-color:#F8F8FF",
        addDiamondName:"添加跑得快红宝石",
        jump:{
            tip:{url:"tip", name:"选择游戏后台"},
            phz:{url:"http://phz.coolgamebox.com:88/login.html", name:"跑胡子会员管理后台"}
        },
        disableAddMember:true,//会员管理页面 添加会员按钮不显示
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:
            {
                16:"经典",
                27:"2副牌",
                15:"15张",
                18:"赖子",
                33:"打筒子",
                j10:"10局",
                j20:"20局",
                j30:"30局",
                j1000:"1000分",
                j600:"600分",
                s0:"不显示",
                s1:"显示",
                f0:"非必出",
                f1:"必出",
                fw:"赢家先",
                f3:"黑3先",
                rw100:"终局100分",
                rw200:"终局200分",
                rw300:"终局300分"
            }
    },
    ddz:{
        loginImg:"images/ddz.png",
        name:"斗地主会员管理后台",
        agentname:"斗地主会员代理后台",
        color:"background-color:#FFFFFF",
        addDiamondName:"添加斗地主红钻",
        jump:{//http://scmj.coolgamebox.com:89/login.html
            tip:{url:"tip", name:"选择游戏后台"},
            scmj:{url:"http://scmj.coolgamebox.com:89/login.html", name:"四川麻将会员管理后台"}
        },
        rechargeOpen:true,
        alipayTextTip:"请在浏览器里面操作充值",
        rebate:{
            basic:{sca:0.3, recommend:0, recharge:0},
            firstLevel:{sca:0.4, recommend:4, recharge:1000},
            secondLevel:{sca:0.5, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.3, recommend:0, recharge:0},
            firstLevel:{sca:0.4, recommend:4, recharge:1000},
            secondLevel:{sca:0.5, recommend:10, recharge:1000}
        },
        dicStatistice :
            {
                d:"斗地主",

                j:"经典",
                s:"四川",
                r:"四人",
                l:"癞子",
                h:"欢乐",
                e:"二人",

                8:"8局",
                10:"10局",
                20:"20局",

                2:"2炸",
                3:"3炸",
                4:"4炸"
            }
    },
    gxmj:{
        loginImg:"images/gxmj.png",
        name:"广西麻将会员管理后台",
        agentname:"广西麻将会员代理后台",
        color:"background-color:#FFFFFF",
        addDiamondName:"添加广西麻将钻石",
        alipayTextTip:"请在浏览器里面操作充值",
        rechargeOpen:true,
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:{
            "gx0":"南宁",
            "gx1":"桂林",
            "gx2":"河池",
            "gx3":"柳州",
            "gx4":"玉林",
            "gx5":"百色",
            "gx6":"贵港",
            "gx7":"钦州",
            "gx8":"转转",
            "j4":"4局",
            "j8":"8局",
            "f0":"不带风",
            "f1":"带风",
            "c0":"不可以吃",
            "c1":"可以吃",
            "p0":"不封胡",
            "p1":"可封胡"
        }
    },
    kwx:{
        loginImg:"images/kwx.png",
        name:"卡五星会员管理后台",
        agentname:"卡五星会员代理后台",
        color:"background-color:#FFFFFF",
        addDiamondName:"添加卡五星钻石",
        alipayTextTip:"请在浏览器里面操作充值",
        rechargeOpen:true,
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:
            {
                SGM0:"无游戏模式",
                SGM1:"襄阳卡五星",
                SGM2:"武汉麻将",
                SGM3:"孝感卡五星",
                SGM4:"十堰卡五星",
                R8:"8局",
                R16:"16局",
                IP1:"带漂",
                IP0:"不带漂",
                MM0:"无买马",
                MM1:"自摸买马",
                MM2:"亮牌自摸买马",
                IKK1:"开口番",
                IKK0:"口口番",
                CE1:"可吃牌",
                CE0:"不可吃牌",
                WW1:"带风",
                WW0:"不带风",
                CEH1:"可吃胡",
                CEH0:"不可吃胡"
            }
    },
    zjmj:{
        loginImg:"images/zjmj.png",
        name:"浙江麻将会员管理后台",
        agentname:"浙江麻将会员代理后台",
        color:"background-color:#FFFFFF",
        addDiamondName:"添加浙江麻将钻石",
        alipayTextTip:"请在浏览器里面操作充值",
        rechargeOpen:true,
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice :{
            hz:"杭州",
            nb:"宁波",
            tz:"台州",
            wz:"温州",
            sx:"绍兴",
            j4:"4局",
            j8:"8局",
            j16:"16局",
            t4:"4台",
            f8:"8番",
            f16:"16番",
            hu300:"300糊",
            hu500:"500糊",
            baida3:"3百搭",
            baida7:"7百搭",
        }
    },
    gdmj:{
        loginImg:"images/gdmj.png",
        name:"广东麻将会员管理后台",
        agentname:"广东麻将会员代理后台",
        color:"background-color:#FFFFFF",
        addDiamondName:"添加广东麻将钻石",
        alipayTextTip:"请在浏览器里面操作充值",
        rechargeOpen:true,
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice :{
            t1:"广州推倒胡",
            t2:"惠州庄麻将",
            t3:"深圳麻将",
            t4:"东莞麻将",
            t5:"鸡平胡",
            t6:"100张",
            r4:"4局",
            r8:"8局",
            g0:"无鬼",
            g1:"红中鬼",
            g2:"翻鬼",
            b0:"无爆炸马",
            b1:"爆炸马",
            s3:"三人推倒胡",
            s4:"四人玩法",
            j0:"无节节高",
            j1:"节节高",
            nj0:"无惠州不可鸡胡",
            nj1:"惠州不可鸡胡",
            m0:"无惠州马跟底",
            m1:"惠州马跟底",
            mq0:"无惠州门清",
            mq1:"惠州门清",
            md0:"无惠州马跟对对胡",
            md1:"惠州马跟对对胡",
            h2:"2马",
            h4:"4马",
            h6:"6马",
            f0:"无风牌",
            f1:"有风牌",
            d0:"不可胡七对",
            d1:"可胡七对",
            c0:"不可吃",
            c1:"可吃",
            p0:"可点炮",
            p1:"不可点炮"
        }
    },
    sdmj:{
        loginImg:"images/sdmj.png",
        name:"山东麻将会员管理后台",
        agentname:"山东麻将会员代理后台",
        color:"background-color:#FFFFFF",
        addDiamondName:"添加山东麻将钻石",
        alipayTextTip:"请在浏览器里面操作充值",
        rechargeOpen:true,
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:{
            "jn":"济南麻将",
            "j4":"4局",
            "j8":"8局",
            "zj1":"258做将",
            "zj2":"不是258做将",
            "df1":"带风",
            "df2":"不带风",
            "dh1":"带花",
            "dh2":"不带花",
            "ml1":"明楼",
            "ml2":"不明楼",
            "sby":"手把一",
            "j6":"6局",
            "j12":"12局",
            "bp1":"包牌",
            "bp2":"不包牌",
            "byw1":"把一是大小王",
            "byw2":"把一不是大小王",
            "wz1":"大小王可以炸",
            "wz2":"大小王不可以炸"
        }
    },
    shxmj:{
        loginImg:"images/shxmj.png",
        name:"山西麻将会员管理后台",
        agentname:"山西麻将会员代理后台",
        color:"background-color:#FFFFFF",
        addDiamondName:"添加山西麻将钻石",
        alipayTextTip:"请在浏览器里面操作充值",
        rechargeOpen:true,
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:{

        }
    },
    dzpk:{
        loginImg:"images/dzpk.png",
        name:"德州扑克会员管理后台",
        agentname:"德州扑克会员代理后台",
        color:"background-color:#FFFFFF",
        addDiamondName:"添加嗨皮德州钻石",
        alipayTextTip:"请在浏览器里面操作充值",
        rechargeOpen:true,
        rebate:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        rebate2:{
            basic:{sca:0.16, recommend:0, recharge:0},
            firstLevel:{sca:0.2, recommend:4, recharge:1000},
            secondLevel:{sca:0.24, recommend:10, recharge:1000}
        },
        dicStatistice:{

        }
    }
};

function getGameType(host){
    var kk = host.split(".");
    var logo = kk[0];
    return gameType[logo];
}
function getHost(host){
    var kk = host.split(".");
    return kk[0];
}

function changeBodyColor() {
    //document.getElementsByTagName("body")[0].setAttribute("style", getGameType(window.location.host).color);
}

$(document).ready(function(){
    changeBodyColor();
});

function getDateCommonFormat(date) {
    var year = date.getFullYear();
    var month = date.getMonth() + 1;
    var day = date.getDate();
    var hour = date.getHours();
    var min = date.getMinutes();
    var sec = date.getSeconds();

    var cur = "";

    cur += year + "-";

    if(month < 10)cur += "0";

    cur += month + "-";

    if(day < 10)cur += "0";

    cur += day + " ";

    if(hour < 10)cur += "0";

    cur += hour + ":";

    if(min < 10)cur += "0";

    cur += min + ":";

    if(sec < 10)cur += "0";

    cur += sec;

    return cur;
}

function getFormatDate(date) {
    //yymmdd
    var year = 0;
    var month = 0;
    var day = 0;
    year = date.getFullYear();
    month = date.getMonth() + 1;
    day = date.getDate();

    var cur = "";
    cur += year;

    if(month >= 10){
        cur += month;
    }
    else {
        cur += "0" + month;
    }

    if(day >= 10){
        cur += day;
    }
    else {
        cur += "0" + day;
    }

    return parseInt(cur);
}

function analyzeStatistics(game, statiStr)//statiStr  统计字符串 "nj_j01_js100";
{
    var data = statiStr.split('_');
    var endStr = "";

    for(var dataI = 0;dataI < data.length;dataI++)
    {
        if(endStr != "") {
            endStr += '_';
        }

        if(!gameType[game].dicStatistice[data[dataI]]) {
            endStr += data[dataI];
        } else {
            endStr += gameType[game].dicStatistice[data[dataI]];
        }
    }

    if(endStr == "") {
        return statiStr;
    } else {
        return endStr;
    }
}
function getStrDate(date) {
    //yymmdd
    var year = 0;
    var month = 0;
    var day = 0;
    year = date.getFullYear();
    month = date.getMonth() + 1;
    day = date.getDate();

    var cur = "";
    cur += year+"-";

    if(month >= 10){
        cur += month+"-";
    }
    else{
        cur += "0" + month+"-";
    }
    if(day >= 10){
        cur += day;
    }
    else {
        cur += "0" + day;
    }
    return cur;
}
function replaceAll(str , replaceKey , replaceVal){
    return str.replace(replaceKey,replaceVal || '');
};
function logout()
{
    postJson("/login/doLogout",{  },function(rtn)
    {
        window.parent.location.href="/login.html";

    });
}
Date.prototype.Format = function (fmt)
{
    var o = {
        "M+": this.getMonth() + 1,
        "d+": this.getDate(),
        "h+": this.getHours(),
        "m+": this.getMinutes(),
        "s+": this.getSeconds(),
        "q+": Math.floor((this.getMonth() + 3) / 3),
        "S": this.getMilliseconds()
    };
    if (/(y+)/.test(fmt))
    {
        fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    }
    for (var k in o)
    {
        if (new RegExp("(" + k + ")").test(fmt))
        {
            fmt =
                fmt.replace(RegExp.$1,
                    (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
        }
    }
    return fmt;
};

