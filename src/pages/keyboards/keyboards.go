package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
)

var BackToStartKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"Back to start page",
			consts.BACK_TO_START,
		),
	),
)

var StartPageKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"1️⃣ I'm representing a DEX / Bridge / Chain / Aggregator / Intent protocol",
			consts.COLLABORATE,
		),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"2️⃣ I'm interested in your API / integration docs",
			consts.INTEGRATE,
		),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"3️⃣ I want to report a bug or product issue",
			consts.SUPPORT,
		),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"4️⃣ Something else",
			consts.OTHER,
		),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL(
			"✉️ Chat with support team directly",
			"https://t.me/RubicSupportBot",
		),
	),
)

var SupportPageKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL(
			"✉️ Chat with support team directly",
			"https://t.me/RubicSupportBot",
		),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"Describe issue.",
			consts.DESCRIBE_ISSUE,
		),
	),
	backToStartButton,
)

var IntegrationPageKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	sdkApiLink,
	integrationGuideLink,
	productOverviewLink,
	supportBotLink,
	backToStartButton,
)

// admin pages keyboards
var (
	AdminStartPageKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔍 Show messages",
				consts.SHOW_MESSAGES,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"👤 Show messages of specific user",
				consts.SHOW_MESSAGES_OF_SPECIFIC_USER,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"📉 Show received messages count from users",
				consts.SHOW_MESSAGES_COUNT_BY_TIME,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔗 Check links",
				consts.CHECK_LINKS,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"❌ Delete oldest messages",
				consts.DELETE_MESSAGES,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"❌ Delete messages of specific user",
				consts.DELETE_MESSAGES_OF_USER,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🚷 Block user",
				consts.BLOCK_USER,
			),
		),
	)

	AdminSelectTimeForMsgCountPageKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Last 10 minutes",
				consts.TIME_TO_MIN.Mins_10,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Last 30 minutes",
				consts.TIME_TO_MIN.Mins_30,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Last hour",
				consts.TIME_TO_MIN.Hour_1,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Last 3 hours",
				consts.TIME_TO_MIN.Hours_3,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Last 6 hours",
				consts.TIME_TO_MIN.Hours_6,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Last 12 hours",
				consts.TIME_TO_MIN.Hours_12,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Last day",
				consts.TIME_TO_MIN.Day_1,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Last 3 days",
				consts.TIME_TO_MIN.Days_3,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Last week",
				consts.TIME_TO_MIN.Week_1,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Last 2 weeks",
				consts.TIME_TO_MIN.Weeks_2,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Last month",
				consts.TIME_TO_MIN.Month_1,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Last 3 months",
				consts.TIME_TO_MIN.Months_3,
			),
		),
		backToStartButton,
	)

	AdminLinksPageKeyboard = IntegrationPageKeyboard

	AdminListOfLinksPageKeyboard = BackToStartKeyBoard

	AdminInfoAfterDeletionMsgPageKeyboard = BackToStartKeyBoard

	AdminCountOfReceivedMsgsPageKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔍 Show messages",
				consts.SHOW_MESSAGES,
			),
		),
		backToStartButton,
	)

	AdminOldOrNewMessagesPageKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Show only new messages",
				consts.SHOW_NEW_MESSAGES,
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Show all messages",
				consts.SHOW_ALL_MESSAGES,
			),
		),
		backToStartButton,
	)
)
