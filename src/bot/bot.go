package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"net/http"
	"os"
	"telegram-stats-bot/src/analytic"
)

type Bot struct {
	token  string
	tgBot  *tgbotapi.BotAPI
	worker *analytic.Worker
}

func InitBot(worker *analytic.Worker) (*Bot, error) {
	token := os.Getenv("TG_BOT_TOKEN")
	tgBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		token:  token,
		tgBot:  tgBot,
		worker: worker,
	}, nil
}

func (b *Bot) RunBot() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.tgBot.GetUpdatesChan(u)

	for update := range updates {
		go b.handleUpdate(update)
	}
	return nil
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	fmt.Printf("handle new update %v\n", update)
	if update.Message != nil {
		if update.Message.Document != nil {
			text := b.handleDocumentMessage(update.Message.Document)
			msg := tgbotapi.NewMessage(update.FromChat().ID, text)

			_, err := b.tgBot.Send(msg)
			if err != nil {
				fmt.Printf("Error while send telegram message: %v\n", err)
				return
			}
		} else if update.Message.IsCommand() {
			text := b.handleCommandMessage(update.Message.Command())
			msg := tgbotapi.NewMessage(update.FromChat().ID, text)

			_, err := b.tgBot.Send(msg)
			if err != nil {
				fmt.Printf("Error while send telegram message: %v\n", err)
				return
			}
		}
	}
}

func (b *Bot) handleCommandMessage(command string) string {
	if command == "help" || command == "start" {
		return "В приложении телеграмм на компьютере откройте нужный чат. " +
			"Далее нажмите на символ трех точек в верхнем правом углу чата. " +
			"Там выберите пункт \"Экспорт истории чата\".\n\n" +
			"В открывшимся меню уберите все галочки и выберите в пункте формат \"Машиночитаемый JSON\". " +
			"Затем дождитесь окончания загрузки и отправьте полученный файл в этот чат."
	} else {
		return "Неизвестная команда. Чтобы получить инструкцию к боту выполните команду /help"
	}
}

func (b *Bot) handleDocumentMessage(document *tgbotapi.Document) string {
	if document.FileSize >= 20*1024*1024 {
		return fmt.Sprintf("Слишком большой файл. Максимальный допустимый размер файла до 20мб. " +
			"Попробуйте загрузить другой файл")
	}
	data, err := b.downloadFile(document)
	if err != nil {
		fmt.Printf("Error while download file: %v\n", err)
		return fmt.Sprintf("Произошла ошибка во время загрузки файла. Попробуйте позже")
	}

	stats, err := b.worker.AnalyseJson(data)
	if err != nil {
		fmt.Printf("Error while parsing file: %v\n", err)
		return fmt.Sprintf("Произошла ошибка во время обработки файла. " +
			"Проверьте, что вы правильно выгрузили файл. " +
			"Инструкцию можно посмотреть вызвав команду /help")
	}

	text := fmt.Sprintf("Статистика по чату %s\n"+
		"Количество сообщений: %d\n"+
		"Количество слов: %d\n"+
		"Количество символов: %d\n\n",
		stats.Chat.Name, stats.MessageCount, stats.WordsCount, stats.CharCount)

	text += "Список пользователей:\n"

	persons := stats.GetMostActiveProfile(10)
	for _, person := range persons {
		text += fmt.Sprintf("%s - сообщения %d; слова %d; символы %d\n",
			person.Name, person.MessageCount, person.WordsCount, person.CharCount)
	}

	text += fmt.Sprintf("\nСамые активные дни:\n")

	dates := stats.GetMostPopularDate(10)
	for _, date := range dates {
		year, month, day := date.Key.Date()
		text += fmt.Sprintf("%.2d.%.2d.%.4d - %d\n",
			day, month, year, date.Value)
	}

	return text
}

func (b *Bot) downloadFile(document *tgbotapi.Document) ([]byte, error) {
	config := tgbotapi.FileConfig{
		FileID: document.FileID,
	}

	file, err := b.tgBot.GetFile(config)
	if err != nil {
		return nil, err
	}
	fileLink := file.Link(b.token)

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := client.Get(fileLink)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
