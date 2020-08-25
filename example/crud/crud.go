package main

import (
	"bufio"
	"fmt"

	"gitlab.com/tsuchinaga/gocli"
)

func main() {
	gocli.SetCommandNotExistsMessage("コマンドがありません")
	gocli.SetHelpDescription("今ひらいているヘルプ")
	gocli.SetReturnDescription("ひとつ前に戻る")
	gocli.SetExitDescription("コマンドの終了")

	cli := gocli.NewGocli()
	cli.AddSubCommand(gocli.NewCommand("list", "一覧表示").SetAction(func(*bufio.Scanner) gocli.AfterAction {
		fmt.Println("一覧表示をするよ")
		return gocli.AfterActionReturn
	})).AddSubCommand(gocli.NewCommand("create", "新規追加").SetAction(func(*bufio.Scanner) gocli.AfterAction {
		fmt.Println("新規追加をするよ")
		return gocli.AfterActionReturn
	})).AddSubCommand(gocli.NewCommand("update", "更新").SetAction(func(*bufio.Scanner) gocli.AfterAction {
		fmt.Println("更新をするよ")
		return gocli.AfterActionKeep
	}).AddSubCommand(
		gocli.NewCommand("id", "IDで指定して更新するよ").SetAction(func(*bufio.Scanner) gocli.AfterAction {
			fmt.Println("IDで検索して更新するよ")
			return gocli.AfterActionReturn
		})).
		AddSubCommand(gocli.NewCommand("no", "No.で指定して更新するよ").SetAction(func(*bufio.Scanner) gocli.AfterAction {
			fmt.Println("IDで検索して更新するよ")
			return gocli.AfterActionReturn
		})),
	).AddSubCommand(gocli.NewCommand("delete", "削除").SetAction(func(*bufio.Scanner) gocli.AfterAction {
		fmt.Println("削除をするよ")
		return gocli.AfterActionReturn
	}))

	if err := cli.Run(); err != nil {
		panic(err)
	}
}
