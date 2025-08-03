package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"syscall/js"
	"time"

	"github.com/a-h/templ"
)

type Message struct {
	Username  string `json:"username"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

var messages []Message
var messageStats = make(map[string]int)
var chart js.Value

// Helper function to render templ component to HTML string
func renderTemplToString(component templ.Component) string {
	var buf bytes.Buffer
	ctx := context.Background()
	err := component.Render(ctx, &buf)
	if err != nil {
		return ""
	}
	return buf.String()
}

func updateStatsDisplay() {
	document := js.Global().Get("document")
	statsContent := document.Call("getElementById", "statsContent")
	if statsContent.IsNull() {
		return
	}

	// Clear existing content
	statsContent.Set("innerHTML", "")

	// Create overview section HTML
	overviewHTML := `<div class="bg-surface1 rounded-lg p-4 border border-surface2">
		<h3 class="text-lg font-semibold text-text mb-3">Overview</h3>
		<div class="grid grid-cols-2 gap-4">`

	// Add stats cards using templ components
	totalMessagesCard := renderTemplToString(StatsCard("Total Messages", strconv.Itoa(len(messages)), "text-blue"))
	activeUsersCard := renderTemplToString(StatsCard("Active Users", strconv.Itoa(len(messageStats)), "text-green"))

	overviewHTML += totalMessagesCard + activeUsersCard + `</div></div>`

	// Add user breakdown section
	if len(messageStats) > 0 {
		overviewHTML += `<div class="bg-surface1 rounded-lg border border-surface2 mt-4">
			<h3 class="text-lg font-semibold text-text p-4 border-b border-surface2">User Activity</h3>
			<div class="p-4 space-y-3">`

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
			userCard := renderTemplToString(UserCard(user.name, user.count))
			overviewHTML += userCard
		}

		overviewHTML += `</div></div>`
	}

	statsContent.Set("innerHTML", overviewHTML)
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

	// Build HTML string for all messages using templ components
	var messagesHTML string
	for _, msg := range messages {
		messageHTML := renderTemplToString(MessageComponent(msg.Username, msg.Text, msg.Timestamp))
		messagesHTML += messageHTML
	}

	messagesElement.Set("innerHTML", messagesHTML)
	messagesElement.Set("scrollTop", messagesElement.Get("scrollHeight"))
}

func createModernMessage(msg Message) js.Value {
	html := renderTemplToString(MessageComponent(msg.Username, msg.Text, msg.Timestamp))

	// Create a temporary container to parse the HTML
	document := js.Global().Get("document")
	tempDiv := document.Call("createElement", "div")
	tempDiv.Set("innerHTML", html)

	// Return the first child element
	return tempDiv.Get("firstElementChild")
}

func createUserAvatar(username string) js.Value {
	html := renderTemplToString(UserAvatar(username, "w-12 h-12"))

	// Create a temporary container to parse the HTML
	document := js.Global().Get("document")
	tempDiv := document.Call("createElement", "div")
	tempDiv.Set("innerHTML", html)

	// Return the first child element
	return tempDiv.Get("firstElementChild")
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

func addRandomMessages() {
	// Random usernames
	usernames := []string{
		"Alice", "Bob", "Charlie", "Diana", "Eve", "Frank", "Grace", "Henry",
		"Ivy", "Jack", "Kate", "Leo", "Maya", "Noah", "Olivia", "Pete",
		"Quinn", "Ruby", "Sam", "Tara", "Uma", "Victor", "Wendy", "Xander",
		"Yara", "Zoe", "Alex", "Blake", "Casey", "Drew", "Emery", "Finley",
	}

	// Random message templates
	messageTemplates := []string{
		"Hey everyone! 👋",
		"How's everyone doing today?",
		"Just finished a great project!",
		"Anyone else excited about the weekend?",
		"Coffee time! ☕",
		"Working on something cool...",
		"Beautiful day outside! ☀️",
		"Just discovered this amazing chat app!",
		"Quick question - anyone here?",
		"Love the new design! 🎨",
		"Testing this feature out...",
		"Hope everyone is staying safe!",
		"Such a productive day!",
		"Anyone want to collaborate?",
		"This chat is really smooth!",
		"Great to see everyone here!",
		"Just saying hello! 😊",
		"Loving the dark theme!",
		"Quick break from coding...",
		"The stats feature is neat!",
		"Random message time! 🎲",
		"Everything looks so modern!",
		"Nice work on the UI!",
		"Chat bubbles look amazing!",
		"The colors are perfect! 🌈",
	}

	// Generate 5 random messages
	for i := 0; i < 5; i++ {
		// Pick random username and message
		username := usernames[len(messages)%len(usernames)]
		messageText := messageTemplates[len(messages)%len(messageTemplates)]

		message := Message{
			Username:  username,
			Text:      messageText,
			Timestamp: time.Now().Format("15:04:05"),
		}

		messages = append(messages, message)
		messageStats[username]++
	}

	updateMessagesDisplay()
	saveToLocalStorage()
	updateStatsDisplay()
	updateChart()
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

	js.Global().Set("addRandomMessages", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		addRandomMessages()
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
