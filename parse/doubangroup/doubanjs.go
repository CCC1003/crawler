package doubangroup

import (
	"crawler/collect"
	"time"
)

var DoubangroupJSTask = &collect.TaskModel{
	Property: collect.Property{
		Name:     "js_find_douban_sun_room",
		WaitTime: 1 * time.Second,
		MaxDepth: 5,
		Cookie:   "bid=znKnt-7lWzE; _ga=GA1.1.1024234719.1704808192; _ga_RXNMP372GL=GS1.1.1704808192.1.0.1704808200.52.0.0; viewed=\"1007305\"; ll=\"108303\"; doubangroup-fav-remind=1; ct=y; dbcl2=\"277697100:TT+ck9rXL/M\"; push_noty_num=0; push_doumail_num=0; __utmv=30149280.27769; ck=83M2; __utmc=30149280; ap_v=0,6.0; __utma=223695111.1024234719.1704808192.1706683118.1706683118.1; __utmb=223695111.0.10.1706683118; __utmc=223695111; __utmz=223695111.1706683118.1.1.utmcsr=doubangroup.com|utmccn=(referral)|utmcmd=referral|utmcct=/group/szsh/discussion; _pk_ref.100001.4cf6=%5B%22%22%2C%22%22%2C1706683118%2C%22https%3A%2F%2Fwww.doubangroup.com%2Fgroup%2Fszsh%2Fdiscussion%3Fstart%3D0%22%5D; _pk_id.100001.4cf6=ce808ff896b4422e.1706683118.; _pk_ses.100001.4cf6=1; _vwo_uuid_v2=D5A97586AE54B6E7528B2983C5376EC23|1e344841f547552ff35c1a1717360b9a; __utma=30149280.1024234719.1704808192.1706681767.1706684046.28; __utmz=30149280.1706684046.28.4.utmcsr=bing|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); __utmt=1; __utmb=30149280.2.10.1706684046",
	},
	Root: `
		var arr = new Array();
 		for (var i = 0; i <= 0; i+=25) {
			var obj = {
			   Url: "https://www.doubangroup.com/group/szsh/discussion?start=" + i,
			   Priority: 1,
			   RuleName: "解析网站URL",
			   Method: "GET",
		   };
			arr.push(obj);
		};
		console.log(arr[0].Url);
		AddJsReq(arr);
			`,
	Rules: []collect.RuleModel{
		{
			Name: "解析网站URL",
			ParseFunc: `
			ctx.ParseJSReg("解析阳台房","(https://www.doubangroup.com/group/topic/[0-9a-z]+/)\"[^>]*>([^<]+)</a>");
			`,
		},
		{
			Name: "解析阳台房",
			ParseFunc: `
			//console.log("parse output");
			ctx.OutputJS("<div class=\"topic-content\">[\\s\\S]*?阳台[\\s\\S]*?<div class=\"aside\">");
			`,
		},
	},
}
