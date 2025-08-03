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
	overviewDiv := createDiv("bg-surface1 rounded-lg p-4 border border-surface2")
	overviewTitle := createElement("h3")
	overviewTitle.Set("className", "text-lg font-semibold text-text mb-3")
	overviewTitle.Set("textContent", "Overview")
	overviewDiv.Call("appendChild", overviewTitle)

	// Stats grid
	statsGrid := createDiv("grid grid-cols-2 gap-4")

	// Total messages
	totalCard := createStatsCard("Total Messages", strconv.Itoa(len(messages)), "text-blue")
	statsGrid.Call("appendChild", totalCard)

	// Active users
	activeUsersCard := createStatsCard("Active Users", strconv.Itoa(len(messageStats)), "text-green")
	statsGrid.Call("appendChild", activeUsersCard)
	overviewDiv.Call("appendChild", statsGrid)
	statsContent.Call("appendChild", overviewDiv)

	// User breakdown section
	if len(messageStats) > 0 {
		usersDiv := createDiv("bg-surface1 rounded-lg border border-surface2")
		usersTitle := createElement("h3")
		usersTitle.Set("className", "text-lg font-semibold text-text p-4 border-b border-surface2")
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
	card := createDiv("bg-surface1 rounded-lg border border-surface2 p-4 text-center")

	valueDiv := createElement("div")
	valueDiv.Set("className", "text-2xl font-bold "+colorClass)
	valueDiv.Set("textContent", value)

	titleDiv := createElement("div")
	titleDiv.Set("className", "text-sm text-subtext0 mt-1")
	titleDiv.Set("textContent", title)

	card.Call("appendChild", valueDiv)
	card.Call("appendChild", titleDiv)

	return card
}

// Create modern user card
func createModernUserCard(username string, messageCount int) js.Value {
	card := createDiv("flex items-center justify-between p-3 bg-surface1 border border-surface2 rounded-lg hover:bg-surface2 transition-colors")

	// Left side with avatar and name
	leftDiv := createDiv("flex items-center space-x-3")

	// Create avatar using the same function as messages
	avatar := createUserAvatar(username)

	// Create name
	nameDiv := createElement("div")
	nameDiv.Set("className", "font-medium text-text")
	nameDiv.Set("textContent", username)

	leftDiv.Call("appendChild", avatar)
	leftDiv.Call("appendChild", nameDiv)

	// Right side with message count
	rightDiv := createDiv("text-right")

	countDiv := createElement("div")
	countDiv.Set("className", "text-lg font-semibold text-text")
	countDiv.Set("textContent", strconv.Itoa(messageCount))

	labelDiv := createElement("div")
	labelDiv.Set("className", "text-xs text-subtext0")
	labelDiv.Set("textContent", "messages")

	rightDiv.Call("appendChild", countDiv)
	rightDiv.Call("appendChild", labelDiv)

	card.Call("appendChild", leftDiv)
	card.Call("appendChild", rightDiv)

	return card
}

func toggleStats() {
	document := js.Global().Get("document")
	statsModal := document.Call("getElementById", "statsModal")

	if statsModal.Get("classList").Call("contains", "hidden").Bool() {
		statsModal.Get("classList").Call("remove", "hidden")
		updateStatsDisplay()
		if chart.IsUndefined() {
			initChart()
		}
		updateChart()
	} else {
		statsModal.Get("classList").Call("add", "hidden")
	}
}
func closeStats() {
	document := js.Global().Get("document")
	statsModal := document.Call("getElementById", "statsModal")
	statsModal.Get("classList").Call("add", "hidden")
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

	// Create chart using JavaScript eval
	js.Global().Call("eval", `
		const canvas = document.getElementById('messageChart');
		
		// Destroy existing chart if it exists
		if (window.messageChart && typeof window.messageChart.destroy === 'function') {
			window.messageChart.destroy();
		}
		
		window.messageChart = new Chart(canvas.getContext('2d'), {
			type: 'doughnut',
			data: {
				labels: [],
				datasets: [{
					data: [],
					backgroundColor: [
						'#8B5CF6', '#06B6D4', '#10B981', '#F59E0B', 
						'#EF4444', '#EC4899', '#6366F1', '#84CC16'
					],
					borderWidth: 2,
					borderColor: '#ffffff'
				}]
			},
			options: {
				responsive: true,
				maintainAspectRatio: true,
				cutout: '60%',
				plugins: {
					legend: {
						position: 'bottom',
						labels: {
							padding: 20,
							usePointStyle: true,
							font: {
								size: 12
							}
						}
					},
					tooltip: {
						callbacks: {
							label: function(context) {
								const label = context.label || '';
								const value = context.parsed;
								const total = context.dataset.data.reduce((a, b) => a + b, 0);
								const percentage = ((value / total) * 100).toFixed(1);
								return label + ': ' + value + ' messages (' + percentage + '%)';
							}
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
		messageDiv := createModernMessage(msg)
		messagesElement.Call("appendChild", messageDiv)
	}

	messagesElement.Set("scrollTop", messagesElement.Get("scrollHeight"))
}

func createModernMessage(msg Message) js.Value {
	// Main message container with enhanced styling
	messageDiv := createDiv("flex items-start space-x-4 mb-6 group hover:bg-surface1/50 p-4 rounded-xl transition-all duration-300 message-enter hover:scale-[1.01] hover:shadow-lg")

	// User avatar
	avatar := createUserAvatar(msg.Username)
	messageDiv.Call("appendChild", avatar)

	// Message content container
	contentDiv := createDiv("flex-1 min-w-0")

	// Header with username and timestamp
	headerDiv := createDiv("flex items-baseline space-x-2 mb-2")

	usernameDiv := createElement("span")
	usernameDiv.Set("className", "font-semibold text-text text-sm")
	usernameDiv.Set("textContent", msg.Username)

	timestampDiv := createElement("span")
	timestampDiv.Set("className", "text-xs text-subtext0")
	timestampDiv.Set("textContent", msg.Timestamp)

	headerDiv.Call("appendChild", usernameDiv)
	headerDiv.Call("appendChild", timestampDiv)

	// Message bubble with enhanced chat-like styling
	bubbleDiv := createDiv("bg-surface0 border border-surface2 rounded-2xl px-4 py-3 shadow-lg hover:shadow-xl hover:border-overlay0 hover:bg-surface1 transition-all duration-300 relative before:absolute before:left-[-8px] before:top-4 before:w-0 before:h-0 before:border-t-8 before:border-t-transparent before:border-b-8 before:border-b-transparent before:border-r-8 before:border-r-surface0")

	textDiv := createElement("p")
	textDiv.Set("className", "text-text text-base leading-relaxed m-0 font-medium")
	textDiv.Set("textContent", msg.Text)

	bubbleDiv.Call("appendChild", textDiv)

	// Assemble the message
	contentDiv.Call("appendChild", headerDiv)
	contentDiv.Call("appendChild", bubbleDiv)
	messageDiv.Call("appendChild", contentDiv)

	return messageDiv
}

func createUserAvatar(username string) js.Value {
	avatar := createDiv("w-12 h-12 rounded-full flex items-center justify-center text-white font-bold flex-shrink-0 shadow-lg border-2 border-surface2 hover:scale-110 transition-transform duration-200")

	// Generate a consistent Catppuccin color based on username
	colors := []string{
		"bg-mauve", "bg-blue", "bg-green", "bg-yellow",
		"bg-red", "bg-pink", "bg-sapphire", "bg-teal",
		"bg-peach", "bg-lavender", "bg-sky", "bg-maroon",
	}

	// Simple hash function to pick color consistently
	hash := 0
	for i := 0; i < len(username); i++ {
		hash += int(username[i])
	}
	colorClass := colors[hash%len(colors)]

	avatar.Get("classList").Call("add", colorClass)

	// Set initial
	initial := "?"
	if len(username) > 0 {
		initial = string(username[0])
	}
	avatar.Set("textContent", initial)

	return avatar
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
