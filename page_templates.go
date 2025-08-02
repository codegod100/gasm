package main

import "strconv"

// Main chat page template
templ ChatPage() {
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
		<title>Go WASM Chat with Templ</title>
		@ChatStyles()
	</head>
	<body>
		@ChatContainer()
		@Scripts()
	</body>
	</html>
}

// Main chat container
templ ChatContainer() {
	<div class="chat-container">
		@ChatHeader()
		@StatsSection()
		@MessagesArea()
		@ChatInput()
	</div>
}

// Chat header component
templ ChatHeader() {
	<div class="chat-header">
		<h1>Go WASM Chat with Templ</h1>
	</div>
}

// Statistics section
templ StatsSection() {
	<div class="stats-section" id="statsSection" style="display: none">
		<div class="stats-text" id="statsText">
			<h3>Message Statistics</h3>
			<p>Loading...</p>
		</div>
		<div class="chart-container">
			<canvas id="messageChart"></canvas>
		</div>
	</div>
}

// Messages area
templ MessagesArea() {
	<div class="chat-messages" id="messages"></div>
}

// Chat input area
templ ChatInput() {
	<div class="chat-input">
		<input
			type="text"
			id="username"
			class="username-input"
			placeholder="Your name"
			value="User"
		/>
		<input
			type="text"
			id="messageInput"
			class="message-input"
			placeholder="Type your message..."
			onkeypress="if(event.key==='Enter') addMessage()"
		/>
		<button onclick="addMessage()" class="send-button">Send</button>
		<button onclick="clearMessages()" class="clear-button">Clear</button>
		<button onclick="toggleStats()" class="toggle-stats">Show Stats</button>
	</div>
}

// Scripts section
templ Scripts() {
	<script src="node_modules/chart.js/dist/chart.umd.js"></script>
	<script src="wasm_exec.js"></script>
	<script>
		const go = new Go();
		WebAssembly.instantiateStreaming(
			fetch("main.wasm"),
			go.importObject,
		).then((result) => {
			go.run(result.instance);
		});
	</script>
}