package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	sdkApiLink = tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL(
			"🔧 SDK / API Docs",
			"https://github.com/Cryptorubic/rubic-sdk",
		),
	)

	integrationGuideLink = tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL(
			"🧠 Integration Guide",
			"https://docs.rubic.finance/integrate-sdk/sdk-overview",
		),
	)

	productOverviewLink = tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL(
			"🚀 Product Overview",
			"https://app.rubic.exchange/?fromChain=ETH&toChain=ETH",
		),
	)
)
