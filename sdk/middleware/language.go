package middleware

import (
	"context"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"golang.org/x/text/language"
)

// Language Set language I18n.
func Language(r *ghttp.Request) {
	ctx := r.GetCtx()
	lang := r.Get("lang").String()
	if lang == "" {
		lang = r.Get("language").String()
	}
	ctx = WithLanguage(ctx, lang, r)
	r.SetCtx(ctx)
	r.Middleware.Next()
}

func WithLanguage(ctx context.Context, lang string, r *ghttp.Request) context.Context {
	var preferred []language.Tag
	var err error
	// 语言
	if lang != "" {
		preferred, _, err = language.ParseAcceptLanguage(lang)
		if err != nil {
			// log err
			lang = "en"
		}
	}
	if preferred == nil && r != nil {
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
	return ctx
}
