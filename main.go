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

// HTML Element Creation Helper Functions
func createElement(tag string) js.Value {
	return js.Global().Get("document").Call("createElement", tag)
}

func createTextNode(text string) js.Value {
	return js.Global().Get("document").Call("createTextNode", text)
}

func createButton(text, className string, onclick func()) js.Value {
	button := createElement("button")
	button.Set("textContent", text)
	button.Set("className", className)

	// Create a wrapper function for the onclick handler
	button.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		onclick()
		return nil
	}))

	return button
}

func createDiv(className string) js.Value {
	div := createElement("div")
	if className != "" {
		div.Set("className", className)
	}
	return div
}

func createInput(inputType, placeholder, className string) js.Value {
	input := createElement("input")
	input.Set("type", inputType)
	input.Set("placeholder", placeholder)
	if className != "" {
		input.Set("className", className)
	}
	return input
}

// Dynamic UI Component Creation
func createUserCard(username string, messageCount int) js.Value {
	card := createDiv("user-card")

	// Create avatar (colored circle with initial)
	avatar := createDiv("user-avatar")
	initial := "?"
	if len(username) > 0 {
		initial = string(username[0])
	}
	avatar.Set("textContent", initial)

	// Create user info
	info := createDiv("user-info")

	nameDiv := createDiv("user-name")
	nameDiv.Set("textContent", username)

	countDiv := createDiv("user-count")
	countDiv.Set("textContent", strconv.Itoa(messageCount)+" messages")

	info.Call("appendChild", nameDiv)
	info.Call("appendChild", countDiv)

	card.Call("appendChild", avatar)
	card.Call("appendChild", info)

	return card
}

func createNotification(message, notificationType string) js.Value {
	notification := createDiv("notification notification-" + notificationType)

	icon := createDiv("notification-icon")
	switch notificationType {
	case "success":
		icon.Set("textContent", "✓")
	case "error":
		icon.Set("textContent", "✗")
	case "info":
		icon.Set("textContent", "ℹ")
	default:
		icon.Set("textContent", "!")
	}

	text := createDiv("notification-text")
	text.Set("textContent", message)

	closeBtn := createButton("×", "notification-close", func() {
		notification.Call("remove")
	})

	notification.Call("appendChild", icon)
	notification.Call("appendChild", text)
	notification.Call("appendChild", closeBtn)

	return notification
}

func renderTemplToString(component interface{}) string {
	// This would work with templ components, but for now let's use a simpler approach
	// var buf bytes.Buffer
	// ctx := context.Background()
	// component.Render(ctx, &buf)
	// return buf.String()

	// Fallback to manual HTML generation for now
	return ""
}
func updateStatsDisplay() {
	document := js.Global().Get("document")
	statsText := document.Call("getElementById", "statsText")
	if statsText.IsNull() {
		return
	}

	// Clear existing content but keep the compact layout
	statsText.Set("innerHTML", "")

	// Create total messages display
	totalDiv := createDiv("total-messages")
	totalDiv.Set("textContent", "Total: "+strconv.Itoa(len(messages)))
	statsText.Call("appendChild", totalDiv)

	// Create user cards container using templ-style approach
	if len(messageStats) > 0 {
		usersContainer := createDiv("users-container")

		for user, count := range messageStats {
			userCard := createTemplUserCard(user, count)
			usersContainer.Call("appendChild", userCard)
		}

		statsText.Call("appendChild", usersContainer)
	}
}

// Create user card using templ-inspired approach
func createTemplUserCard(username string, messageCount int) js.Value {
	card := createDiv("user-card")

	// Create avatar
	avatar := createDiv("user-avatar")
	initial := "?"
	if len(username) > 0 {
		initial = string(username[0])
	}
	avatar.Set("textContent", initial)

	// Create user info container
	info := createDiv("user-info")

	// Create name element
	nameDiv := createDiv("user-name")
	nameDiv.Set("textContent", username)

	// Create count element
	countDiv := createDiv("user-count")
	countDiv.Set("textContent", strconv.Itoa(messageCount)+" messages")

	// Assemble the card
	info.Call("appendChild", nameDiv)
	info.Call("appendChild", countDiv)
	card.Call("appendChild", avatar)
	card.Call("appendChild", info)

	return card
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

	// Create minimal chart using JavaScript eval
	js.Global().Call("eval", `
		const canvas = document.getElementById('messageChart');
		
		// Destroy existing chart if it exists
		if (window.messageChart && typeof window.messageChart.destroy === 'function') {
			window.messageChart.destroy();
		}
		
		// Force canvas size
		canvas.width = 30;
		canvas.height = 30;
		canvas.style.width = '30px';
		canvas.style.height = '30px';
		
		window.messageChart = new Chart(canvas.getContext('2d'), {
			type: 'doughnut',
			data: {
				labels: [],
				datasets: [{
					data: [],
					backgroundColor: ['#4CAF50', '#2196F3', '#FF9800', '#E91E63', '#9C27B0'],
					borderWidth: 0
				}]
			},
			options: {
				responsive: false,
				maintainAspectRatio: true,
				cutout: '60%',
				plugins: {
					legend: {
						display: false
					}
				},
				layout: {
					padding: 5
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
func showNotification(message, notificationType string) {
	document := js.Global().Get("document")
	body := document.Get("body")

	notification := createNotification(message, notificationType)
	body.Call("appendChild", notification)

	// Auto-remove after 3 seconds
	js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !notification.IsNull() {
			notification.Call("remove")
		}
		return nil
	}), 3000)
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
		showNotification("Statistics panel opened", "info")
	} else {
		statsSection.Get("style").Set("display", "none")
		showNotification("Statistics panel closed", "info")
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
