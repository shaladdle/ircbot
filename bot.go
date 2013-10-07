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

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

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
    if _, err := fmt.Fprintf(bot.conn, fmtstr+"\r\n", args...); err != nil {
        panic(err)
    }
}

func (bot *Bot) Say(msg string) {
	bot.CmdPrintf("PRIVMSG %s :%s", bot.channel, msg)
}

func main() {
	flag.Parse()

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

        msgPrefix := "PRIVMSG " + *channel

		var msg string

        if strings.Contains(line, msgPrefix) {
            splt := strings.Split(line, msgPrefix)
            if len(splt) != 0 {
                msg = splt[len(splt)-1]
            }
        }

		if line[:4] == "PING" {
			reply := "PONG" + line[4:]
			ircbot.CmdPrintf(reply)
			fmt.Println("me:", reply)
		} else if strings.Contains(line, "MODE") && !joined {
			joined = true
			ircbot.CmdPrintf("JOIN %s", ircbot.channel)
		} else if joined {
            if idx := strings.Index(line, "!"); strings.Contains(msg, *nick+": PING") {
                ircbot.Say("POOOOOOOOOOOOOOOONG!!!!")
            } else if idx != -1 {
                sender := line[1:idx]
                var recipient string
                if strings.Contains(msg, *nick) {
                    recipient = *nick
                }
                handleMessage(Message{sender, recipient, msg}, ircbot)
            }
        }
	}
}

type Message struct {
    sender, recipient, text string
}

type handler func(Message, *Bot)

func randSay(bot *Bot, things []string) {
    thing := r.Intn(len(things))
    bot.Say(things[thing])
}

func berate(msg Message, bot *Bot) {
    badThings := []string{
        msg.sender + ", I don't like you",
        msg.sender + ", That's what your mom said!",
    }

    randSay(bot, badThings)
}

func quote(msg Message, bot *Bot) {
    randSay(bot, []string{
        "If you prick us do we not bleed? If you tickle us do we not laugh? If "+
        "you poison us do we not die? And if you wrong us shall we not revenge?",
        "Love all, trust a few, do wrong to none.",
        "A fool thinks himself to be wise, but a wise man knows himself to be a fool.",
        "If music be the food of love, play on.",
        "It is not in the stars to hold our destiny but in ourselves.",
        "When a father gives to his son, both laugh; when a son gives to his father, both cry.",
        "Better three hours too soon than a minute too late.",
    })
}

func benice(msg Message, bot *Bot) {
    things := []string{
        msg.sender + ", I love you <3",
        msg.sender + ", you're looking classy today",
        msg.sender + ", I completely agree",
        msg.sender + ", amen, brother",
    }
    randSay(bot, things)
}

func question(msg Message, bot *Bot) {
    things := []string{
        msg.sender + ", how can she slap?",
        msg.sender + ", what are you doing later?",
        msg.sender + ", why are you still here?",
        msg.sender + ", what's the average air speed velocity of an unlaiden swallow?",
    }
    randSay(bot, things)
}

func doge(msg Message, bot *Bot) {
	randSay(bot, []string{
		"Ohai, " + *nick + " are delicious",
		"doge",
		"pls",
		"My name is " + *nick + ", pleased to meet you.",
		"PING me if you dare",
		"check me out and contribute to me at http://github.com/shaladdle/ircbot",
	})
}

func waxpoetic(msg Message, bot *Bot) {
}

func ignore(msg Message, bot *Bot) {
}

func handleMessage(msg Message, bot *Bot) {
    var handlers []handler
    if msg.recipient != *nick {
        handlers = []handler{
            berate,
            benice,
            benice,
            benice,
            benice,
            doge,
            doge,
            question,
            waxpoetic,
            quote,
            ignore,
            ignore,
            ignore,
            ignore,
            ignore,
            ignore,
            ignore,
            ignore,
            ignore,
            ignore,
        }
    } else {
        handlers = []handler{
            benice,
            question,
            doge,
            waxpoetic,
            quote,
            ignore,
        }
    }
    hId := r.Intn(len(handlers))

    handlers[hId](msg, bot)
}
