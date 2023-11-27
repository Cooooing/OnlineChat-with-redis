package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	e "onlineChat/entity"
	"os"
	"strconv"
	"strings"
)

func init() {
	file := "./" + "app" + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[onlineChat] - ")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	return
}

// 设置用户名
var username string
var ctx = context.Background()
var channel = "test"

// 消息队列
var message = make(chan e.Message)
var ReceiveMessageChan = make(chan e.Message, 20)
var lastReceiveUsername = ""

func main() {
	log.Println("application run ...")
	e.InitRedisClient()
	defer func() {
		e.CloseRedisClient()
		log.Println("application close ...")
	}()

	for {
		fmt.Print("please enter your nickname:")
		reader := bufio.NewReader(os.Stdin)
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)
		if username != "" {
			break
		}
	}
	ctx = context.WithValue(ctx, "username", username)
	go e.Subscriber(channel)
	go e.Publisher(message)

	p := tea.NewProgram(initialModel())

	go func() {
		for {
			msg := <-e.ReceiveMessageChan
			ReceiveMessageChan <- msg
			p.Send(receiveMsg{})
		}
	}()

	e.ClearTerminal()
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}

type Model interface {
	Init() tea.Cmd
	Update(tea.Msg) (Model, tea.Cmd)
	View() string
}

type model struct {
	viewport       viewport.Model
	messages       []string
	receiveMessage e.Message
	textarea       textarea.Model
	senderStyle    lipgloss.Style
	err            error
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.senderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(e.GetRandomColor()))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case "enter":
			input := strings.TrimSpace(m.textarea.Value())
			if input == "" {
				m.textarea.Reset()
				break
			}
			sendMessage := e.NewSendMessage(ctx, channel, input)
			message <- sendMessage
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	case receiveMsg:
		receiveMessage := <-ReceiveMessageChan
		showMsg := ""
		if lastReceiveUsername == receiveMessage.Username {
			showMsg = receiveMessage.Msg
		} else {
			showMsg = receiveMessage.Username + " " + receiveMessage.SendTime + "\n" + receiveMessage.Msg
			lastReceiveUsername = receiveMessage.Username
		}
		if receiveMessage.Username == ctx.Value("username") {
			num := <-e.SelfMessageReceiveNumChan
			showMsg = showMsg + "   " + strconv.FormatInt(num-1, 10) + "人收到"
		}
		m.messages = append(m.messages, m.senderStyle.Render(showMsg))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

type (
	errMsg error
)
type receiveMsg struct{}

func (m model) View() string {
	return fmt.Sprintf(
		"%s:\n%s\n\n%s",
		channel,
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(1)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 15)
	//vp.SetContent("Welcome to the chat room!\nType a message and press Enter to send.")

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}
