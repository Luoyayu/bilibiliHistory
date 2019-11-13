package main

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
	"time"
)

var DEBUG = false
var AccessToken = ""

var _type string
var _output string
var _format *cli.StringSlice
var _times string
var cnt = 1
var appFs = afero.NewOsFs()

func getHistoryX(num int64) {
	log.Println(_output)
	var f afero.File
	if _output == "file" {
		var err error
		f, err = appFs.Create("history")
		if err != nil {
			log.Fatal("创建文件失败")
		} else {
			defer f.Close()
		}
	}

	max := "0"
	log.Println("开始查询最近", num*20, "条")
	var i int64 = 1
	for i = 1; i <= num; i++ {
		log.Println("获取第", i, "页")
		ret := getHistory(max, AccessToken)
		if ret.Code != 0 {
			log.Println("结束")
		} else {
			max = fmt.Sprint(ret.Cursor.Max)
			//log.Println("游标: ", max)
			if _output == "file" {
				for _, v := range ret.List {
					timeStr := time.Unix(v.ViewAt, 0).Format("2006-01-02 15:04:05")
					_, _ = f.WriteString(fmt.Sprintf("%5d. ", cnt) + timeStr + "\t" + v.Title + "\n")
					cnt++
				}
			} else {
				for _, v := range ret.List {
					timeStr := time.Unix(v.ViewAt, 0).Format("2006-01-02 15:04:05")
					fmt.Println(fmt.Sprintf("%5d. ", cnt) + timeStr + "\t" + v.Title)
					cnt++
				}
			}
			if max == "0" {
				log.Fatal("到底了! 共计", cnt, "条播放记录")

			}
		}
	}

}
func main() {

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "输出到 stdout/file",
			Value:       "stdout",
			Destination: &_output,
		},

		cli.StringSliceFlag{
			Name:  "format, f",
			Usage: "输出格式 title,viewtime,name,mid,uri,duration",
			Value: _format,
		},

		cli.BoolFlag{
			Name:        "debug, dg",
			Usage:       "开启调试",
			Destination: &DEBUG,
		},

		cli.StringFlag{
			Name:        "times, t",
			Usage:       "查看条数, 20倍",
			Value:       "1",
			Destination: &_times,
		},
	}

	app.Name = "bilibili 历史查询工具"
	app.Usage = "可通过用户名密码或者口令登录，进行观看历史查询\n\t\tFlag请置于命令前!"
	app.Version = "1.0.0"
	app.HideHelp = true

	app.Action = func(c *cli.Context) error {
		//log.Println(c.Args())
		if c.NArg() == 0 {
			_ = cli.ShowAppHelp(c)
		}

		//log.Println("times: ", c.String("times"))
		_times = c.String("times")

		//log.Println("output: ", c.String("output"))

		switch c.String("output") {
		case "stdout":
		case "file":
		default:
			log.Fatal("暂不支持该输出")
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "login",
			Aliases: []string{"l"},
			Usage:   "passwd or token",
			Action: func(c *cli.Context) {
				//log.Println("login: ", c.Args().First())
				switch c.Args().First() {
				case "passwd":
					userName := c.Args().Get(1)
					password := c.Args().Get(2)
					if userName == "" || password == "" {
						log.Fatal("用户名或密码不能为空")
					}
					USERNAME = userName
					PASSWORD = password

					AccessToken = doLogin(USERNAME, PASSWORD)
					if AccessToken != "" {
						log.Println("登录成功!")
						log.Println("下次可以使用 AccessToken 登录: ", AccessToken)
					} else {
						log.Fatal("登录失败! 可使用全局标识 --debug开启调试")
					}

				//log.Println(userName, password)
				case "token":
					AccessToken = c.Args().Get(1)
					log.Println(AccessToken)
					if AccessToken == "" {
						log.Fatal("Access Token 不能为空")
					}
				default:
					log.Fatal("输入必须是 passwd 或 token")
				}
			},
			Description: _type,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	if times, err := strconv.ParseInt(_times, 10, 32); err != nil {
	} else {
		if times > 0 && AccessToken != "" {
			getHistoryX(times)
		}
	}

}
