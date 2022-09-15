package main

import (
	"fmt"
	"telegram-stats-bot/src/analytic"
	"telegram-stats-bot/src/bot"
)

func main() {
	//file, err := ioutil.ReadFile("chatExportExamples/Eva2.json")
	//if err != nil {
	//	fmt.Printf("Error while read file: %v\n", err)
	//	return
	//}
	//
	//for i := 0; i < 1; i++ {
	//	_, err := worker.AnalyseJson(file)
	//	if err != nil {
	//		fmt.Printf("Error while generation chat statistics: %v\n", err)
	//		return
	//	}
	//}

	worker := analytic.NewAnalyticWorker(1)
	tgBot, err := bot.InitBot(worker)
	if err != nil {
		fmt.Printf("Error while initializing bot: %v", err)
	}
	err = tgBot.RunBot()
	if err != nil {
		fmt.Printf("error while runing bot: %v", err)
	}
}
