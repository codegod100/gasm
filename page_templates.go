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
	<body class="font-sans m-0 p-5 bg-gray-100 h-screen box-border">
		@ChatContainer()
		@Scripts()
	</body>
	</html>
}

// Main chat container
templ ChatContainer() {
	<div class="max-w-3xl mx-auto bg-white rounded-lg shadow-md h-[calc(100vh-2.5rem)] flex flex-col">
		@ChatHeader()
		@StatsSection()
		@MessagesArea()
		@ChatInput()
	</div>
}

// Chat header component
templ ChatHeader() {
	<div class="bg-chat-green text-white p-4 rounded-t-lg text-center">
		<h1>Go WASM Chat with Templ</h1>
	</div>
}

// Statistics section
templ StatsSection() {
	<div class="p-1 border-b border-gray-200 bg-gray-50 h-15 overflow-hidden hidden" id="statsSection">
		<div class="flex items-center gap-1 mb-0">
			<h3 class="m-0 flex-shrink-0 text-xs">Stats</h3>
			<div class="w-8 h-8 relative m-0 flex justify-center items-center bg-white rounded-md p-1 flex-shrink-0">
				<canvas id="messageChart" width="30" height="30" class="max-w-full max-h-full w-auto h-auto"></canvas>
			</div>
			<div class="m-0 flex-1" id="statsText">
				<p>Loading...</p>
			</div>
		</div>
	</div>
}

// Messages area
templ MessagesArea() {
	<div class="flex-1 p-4 overflow-y-auto border-b border-gray-200" id="messages"></div>
}

// Chat input area
templ ChatInput() {
	<div class="flex p-4 gap-2">
		<input
			type="text"
			id="username"
			class="flex-none w-36 px-2 py-2 border border-gray-300 rounded text-base"
			placeholder="Your name"
			value="User"
		/>
		<input
			type="text"
			id="messageInput"
			class="flex-1 px-2 py-2 border border-gray-300 rounded text-base"
			placeholder="Type your message..."
			onkeypress="if(event.key==='Enter') addMessage()"
		/>
		<button onclick="addMessage()" class="px-4 py-2 bg-chat-green text-white border-none rounded cursor-pointer text-base hover:bg-chat-green-hover">Send</button>
		<button onclick="clearMessages()" class="px-4 py-2 bg-chat-red text-white border-none rounded cursor-pointer text-base hover:bg-chat-red-hover">Clear</button>
		<button onclick="toggleStats()" class="px-4 py-2 bg-chat-blue text-white border-none rounded cursor-pointer text-sm hover:bg-chat-blue-hover">Show Stats</button>
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