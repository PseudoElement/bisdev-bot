package pages

import "github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"

type AbstrUserPage struct {
	*Page
}

func NewAbstrUserPage(injector *injector.AppInjector) *AbstrUserPage {
	p := &AbstrUserPage{
		Page: NewPage(injector),
	}

	return p
}

func (this *AbstrUserPage) IsSelectablePage() bool {
	return true
}
