package middleware

import (
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"golang.org/x/text/language"
)

// Language Set language I18n.
func Language(r *ghttp.Request) {
	var preferred []language.Tag
	var err error
	ctx := r.Context()
	lang := r.Get("lang").String()
	if lang == "" {
		lang = r.Get("language").String()
	}

	if lang != "" {
		preferred, _, err = language.ParseAcceptLanguage(lang)
		if err != nil {
			// log err
			lang = "en"
		}
	}
	if preferred == nil {
		preferred, _, err = language.ParseAcceptLanguage(r.GetHeader("Accept-Language"))
		if err != nil {
			// log err
			lang = "en"
		}
	}
	matcher := language.NewMatcher([]language.Tag{
		language.English,           // 英语
		language.Chinese,           // 中文
		language.SimplifiedChinese, // 简体中文
		language.Malay,             // 马来西亚语 Malay
	})
	code, _, _ := matcher.Match(preferred...)
	base, _ := code.Base()
	lang = base.String()
	ctx = gi18n.WithLanguage(ctx, lang)

	//设置公共信息

	r.SetCtx(ctx)
	r.Middleware.Next()
}
