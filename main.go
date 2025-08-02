package main

import (
	"encoding/json"
	"strconv"
	"syscall/js"
	"time"
)

type Message struct {
	Username  string `json:"username"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

var messages []Message
var messageStats = make(map[string]int)
var chart js.Value

func updateStatsDisplay() {
	document := js.Global().Get("document")
	statsText := document.Call("getElementById", "statsText")
	if statsText.IsNull() {
		return
	}

	html := "<h3>Message Statistics</h3>"
	html += "<p>Total Messages: " + strconv.Itoa(len(messages)) + "</p>"

	if len(messageStats) > 0 {
		html += "<ul>"
		for user, count := range messageStats {
			html += "<li>" + user + ": " + strconv.Itoa(count) + " messages</li>"
		}
		html += "</ul>"
	}

	statsText.Set("innerHTML", html)
}

func initChart() {
	// Check if Chart.js is available
	chartConstructor := js.Global().Get("Chart")
	if chartConstructor.IsUndefined() {
		return
	}

	// Get canvas element
	canvas := js.Global().Get("document").Call("getElementById", "messageChart")
	if canvas.IsNull() {
		return
	}

	// Get 2D context
	ctx := canvas.Call("getContext", "2d")
	if ctx.IsNull() {
		return
	}

	// Create chart using JavaScript eval to avoid Go object conversion issues
	js.Global().Call("eval", `
		const canvas = document.getElementById('messageChart');
		const container = canvas.parentElement;
		canvas.width = container.clientWidth;
		canvas.height = container.clientHeight;
		
		window.messageChart = new Chart(canvas.getContext('2d'), {
			type: 'doughnut',
			data: {
				labels: [],
				datasets: [{
					data: [],
					backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF', '#FF9F40'],
					borderWidth: 2,
					borderColor: '#fff'
				}]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				layout: {
					padding: 20
				},
				plugins: {
					legend: {
						position: 'bottom',
						labels: {
							padding: 20,
							usePointStyle: true
						}
					}
				}
			}
		});
	`)

	chart = js.Global().Get("messageChart")
}

func updateChart() {
	if chart.IsUndefined() {
		return
	}

	// Build arrays for labels and data
	labels := make([]interface{}, 0, len(messageStats))
	data := make([]interface{}, 0, len(messageStats))

	for user, count := range messageStats {
		labels = append(labels, user)
		data = append(data, count)
	}

	// Update chart data using JavaScript
	chart.Get("data").Set("labels", labels)
	chart.Get("data").Get("datasets").Index(0).Set("data", data)
	chart.Call("update")
}
func toggleStats() {
	document := js.Global().Get("document")
	statsSection := document.Call("getElementById", "statsSection")

	if statsSection.Get("style").Get("display").String() == "none" {
		statsSection.Get("style").Set("display", "block")
		updateStatsDisplay()
		if chart.IsUndefined() {
			initChart()
		}
		updateChart()
	} else {
		statsSection.Get("style").Set("display", "none")
	}
}
func saveToLocalStorage() {
	data, err := json.Marshal(messages)
	if err != nil {
		return
	}

	js.Global().Get("localStorage").Call("setItem", "chatMessages", string(data))

	statsData, err := json.Marshal(messageStats)
	if err != nil {
		return
	}

	js.Global().Get("localStorage").Call("setItem", "messageStats", string(statsData))
}

func loadFromLocalStorage() {
	// Load messages
	dataValue := js.Global().Get("localStorage").Call("getItem", "chatMessages")
	if !dataValue.IsNull() && !dataValue.IsUndefined() {
		dataStr := dataValue.String()
		if dataStr != "" {
			var loadedMessages []Message
			err := json.Unmarshal([]byte(dataStr), &loadedMessages)
			if err == nil {
				messages = loadedMessages
				updateMessagesDisplay()
			}
		}
	}

	// Load stats
	statsValue := js.Global().Get("localStorage").Call("getItem", "messageStats")
	if !statsValue.IsNull() && !statsValue.IsUndefined() {
		statsStr := statsValue.String()
		if statsStr != "" {
			var loadedStats map[string]int
			err := json.Unmarshal([]byte(statsStr), &loadedStats)
			if err == nil {
				messageStats = loadedStats
			}
		}
	}
}

func addMessage() {
	document := js.Global().Get("document")
	usernameInput := document.Call("getElementById", "username")
	messageInput := document.Call("getElementById", "messageInput")

	if usernameInput.IsNull() || messageInput.IsNull() {
		return
	}

	username := usernameInput.Get("value").String()
	messageText := messageInput.Get("value").String()

	if messageText == "" {
		return
	}

	if username == "" {
		username = "Anonymous"
	}

	message := Message{
		Username:  username,
		Text:      messageText,
		Timestamp: time.Now().Format("15:04:05"),
	}

	messages = append(messages, message)
	messageStats[username]++

	updateMessagesDisplay()
	saveToLocalStorage()
	updateStatsDisplay()
	updateChart()

	messageInput.Set("value", "")
	messageInput.Call("focus")
}

func clearMessages() {
	messages = []Message{}
	messageStats = make(map[string]int)
	updateMessagesDisplay()
	saveToLocalStorage()
	updateStatsDisplay()
	updateChart()
}

func updateMessagesDisplay() {
	document := js.Global().Get("document")
	messagesElement := document.Call("getElementById", "messages")
	if messagesElement.IsNull() {
		return
	}

	messagesElement.Set("innerHTML", "")

	for _, msg := range messages {
		messageDiv := document.Call("createElement", "div")
		messageDiv.Set("className", "message")

		userDiv := document.Call("createElement", "div")
		userDiv.Set("className", "message-user")
		userDiv.Set("textContent", msg.Username)

		textDiv := document.Call("createElement", "div")
		textDiv.Set("className", "message-text")
		textDiv.Set("textContent", msg.Text)

		timeDiv := document.Call("createElement", "div")
		timeDiv.Set("className", "message-time")
		timeDiv.Set("textContent", msg.Timestamp)

		messageDiv.Call("appendChild", userDiv)
		messageDiv.Call("appendChild", textDiv)
		messageDiv.Call("appendChild", timeDiv)

		messagesElement.Call("appendChild", messageDiv)
	}

	messagesElement.Set("scrollTop", messagesElement.Get("scrollHeight"))
}

func addWelcomeMessage() {
	message := Message{
		Username:  "System",
		Text:      "Welcome to Go WASM Chat!",
		Timestamp: time.Now().Format("15:04:05"),
	}

	messages = append(messages, message)
	messageStats["System"]++
	updateMessagesDisplay()
	saveToLocalStorage()
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("addMessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		addMessage()
		return nil
	}))

	js.Global().Set("clearMessages", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		clearMessages()
		return nil
	}))

	js.Global().Set("toggleStats", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		toggleStats()
		return nil
	}))

	loadFromLocalStorage()
	if len(messages) == 0 {
		addWelcomeMessage()
	}

	<-c
}
