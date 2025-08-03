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
	statsContent := document.Call("getElementById", "statsContent")
	if statsContent.IsNull() {
		return
	}

	// Clear existing content
	statsContent.Set("innerHTML", "")

	// Create overview section
	overviewDiv := createDiv("bg-gray-50 rounded-lg p-4")
	overviewTitle := createElement("h3")
	overviewTitle.Set("className", "text-lg font-semibold text-gray-800 mb-3")
	overviewTitle.Set("textContent", "Overview")
	overviewDiv.Call("appendChild", overviewTitle)

	// Stats grid
	statsGrid := createDiv("grid grid-cols-2 gap-4")

	// Total messages
	totalCard := createStatsCard("Total Messages", strconv.Itoa(len(messages)), "text-blue-600")
	statsGrid.Call("appendChild", totalCard)

	// Active users
	activeUsersCard := createStatsCard("Active Users", strconv.Itoa(len(messageStats)), "text-green-600")
	statsGrid.Call("appendChild", activeUsersCard)

	overviewDiv.Call("appendChild", statsGrid)
	statsContent.Call("appendChild", overviewDiv)

	// User breakdown section
	if len(messageStats) > 0 {
		usersDiv := createDiv("bg-white rounded-lg border border-gray-200")
		usersTitle := createElement("h3")
		usersTitle.Set("className", "text-lg font-semibold text-gray-800 p-4 border-b border-gray-200")
		usersTitle.Set("textContent", "User Activity")
		usersDiv.Call("appendChild", usersTitle)

		usersList := createDiv("p-4 space-y-3")

		// Sort users by message count
		type userStat struct {
			name  string
			count int
		}
		var sortedUsers []userStat
		for user, count := range messageStats {
			sortedUsers = append(sortedUsers, userStat{user, count})
		}

		// Simple bubble sort
		for i := 0; i < len(sortedUsers); i++ {
			for j := 0; j < len(sortedUsers)-1-i; j++ {
				if sortedUsers[j].count < sortedUsers[j+1].count {
					sortedUsers[j], sortedUsers[j+1] = sortedUsers[j+1], sortedUsers[j]
				}
			}
		}

		for _, user := range sortedUsers {
			userCard := createModernUserCard(user.name, user.count)
			usersList.Call("appendChild", userCard)
		}

		usersDiv.Call("appendChild", usersList)
		statsContent.Call("appendChild", usersDiv)
	}
}

// Create modern stats card
func createStatsCard(title, value, colorClass string) js.Value {
	card := createDiv("bg-white rounded-lg border border-gray-200 p-4 text-center")

	valueDiv := createElement("div")
	valueDiv.Set("className", "text-2xl font-bold "+colorClass)
	valueDiv.Set("textContent", value)

	titleDiv := createElement("div")
	titleDiv.Set("className", "text-sm text-gray-600 mt-1")
	titleDiv.Set("textContent", title)

	card.Call("appendChild", valueDiv)
	card.Call("appendChild", titleDiv)

	return card
}

// Create modern user card
func createModernUserCard(username string, messageCount int) js.Value {
	card := createDiv("flex items-center justify-between p-3 bg-gray-50 rounded-lg")

	// Left side with avatar and name
	leftDiv := createDiv("flex items-center space-x-3")

	// Create avatar
	avatar := createDiv("w-10 h-10 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 text-white flex items-center justify-center font-bold text-sm")
	initial := "?"
	if len(username) > 0 {
		initial = string(username[0])
	}
	avatar.Set("textContent", initial)

	// Create name
	nameDiv := createElement("div")
	nameDiv.Set("className", "font-medium text-gray-800")
	nameDiv.Set("textContent", username)

	leftDiv.Call("appendChild", avatar)
	leftDiv.Call("appendChild", nameDiv)

	// Right side with message count
	rightDiv := createDiv("text-right")

	countDiv := createElement("div")
	countDiv.Set("className", "text-lg font-semibold text-gray-800")
	countDiv.Set("textContent", strconv.Itoa(messageCount))

	labelDiv := createElement("div")
	labelDiv.Set("className", "text-xs text-gray-500")
	labelDiv.Set("textContent", "messages")

	rightDiv.Call("appendChild", countDiv)
	rightDiv.Call("appendChild", labelDiv)

	card.Call("appendChild", leftDiv)
	card.Call("appendChild", rightDiv)

	return card
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
	statsModal := document.Call("getElementById", "statsModal")

	if statsModal.Get("classList").Call("contains", "hidden").Bool() {
		statsModal.Get("classList").Call("remove", "hidden")
		updateStatsDisplay()
		showNotification("Statistics opened", "info")
	} else {
		statsModal.Get("classList").Call("add", "hidden")
		showNotification("Statistics closed", "info")
	}
}

func closeStats() {
	document := js.Global().Get("document")
	statsModal := document.Call("getElementById", "statsModal")
	statsModal.Get("classList").Call("add", "hidden")
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

	messageInput.Set("value", "")
	messageInput.Call("focus")
}

func clearMessages() {
	messages = []Message{}
	messageStats = make(map[string]int)
	updateMessagesDisplay()
	saveToLocalStorage()
	updateStatsDisplay()
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

	js.Global().Set("closeStats", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		closeStats()
		return nil
	}))

	loadFromLocalStorage()
	if len(messages) == 0 {
		addWelcomeMessage()
	}

	<-c
}
