package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/textproto"
	"strings"
	"time"
)

var (
	hostport = flag.String("addr", "", "")
	nick     = flag.String("nick", "", "")
	channel  = flag.String("chan", "", "")
)

type Bot struct {
	nick    string
	channel string
	conn    net.Conn
}

func NewBot(nick, channel string) *Bot {
	return &Bot{
		nick:    nick,
		channel: channel,
		conn:    nil,
	}
}

func (bot *Bot) Connect(hostport string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", hostport)
	if err != nil {
		log.Fatal("unable to connect to IRC server ", err)
	}
	bot.conn = conn
	log.Printf("Connected to IRC server %s (%s)\n", hostport, bot.conn.RemoteAddr())
	return bot.conn, nil
}

func (bot *Bot) CmdPrintf(fmtstr string, args ...interface{}) {
	fmt.Fprintf(bot.conn, fmtstr+"\r\n", args...)
}

func (bot *Bot) Say(msg string) {
	bot.CmdPrintf("PRIVMSG %s :%s", bot.channel, msg)
}

func main() {
	flag.Parse()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	says := []string{
		"Ohai, " + *nick + " are delicious",
		"doge",
		"pls",
		"My name is " + *nick + ", pleased to meet you.",
		"PING me if you dare",
		"check me out and contribute to me at http://github.com/shaladdle/ircbot",
	}

	ircbot := NewBot(*nick, *channel)
	conn, err := ircbot.Connect(*hostport)
	if err != nil {
		fmt.Println("error starting connection:", err)
		return
	}
	ircbot.CmdPrintf("USER %s 8 * :%s", ircbot.nick, ircbot.nick)
	ircbot.CmdPrintf("NICK %s", ircbot.nick)
	defer conn.Close()

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	joined := false
	for {
		line, err := tp.ReadLine()
		if err != nil {
			break // break loop on errors
		}
		fmt.Printf("%s\n", line)

		if line[:4] == "PING" {
			reply := "PONG" + line[4:]
			ircbot.CmdPrintf(reply)
			fmt.Println("me:", reply)
		} else if strings.Contains(line, "MODE") && !joined {
			joined = true
			ircbot.CmdPrintf("JOIN %s", ircbot.channel)
		} else if strings.Contains(line, *nick+": PING") || strings.Contains(line, *nick+", PING") {
			ircbot.Say("POOOOOOOOOOOOOOOONG!!!!")
		} else if strings.Contains(line, *nick+":") || strings.Contains(line, *nick+",") {
			ircbot.Say(says[r.Intn(len(says))])
		}
	}
}
