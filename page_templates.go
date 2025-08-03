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
		@StatsModal()
		@Scripts()
	</body>
	</html>
}

// Main chat container
templ ChatContainer() {
	<div class="max-w-3xl mx-auto bg-white rounded-lg shadow-md h-[calc(100vh-2.5rem)] flex flex-col">
		@ChatHeader()
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
		<button onclick="toggleStats()" class="px-4 py-2 bg-purple-600 text-white border-none rounded cursor-pointer text-sm hover:bg-purple-700">Stats</button>
	</div>
}

// New modern stats modal
templ StatsModal() {
	<div id="statsModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 hidden">
		<div class="bg-white rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[80vh] overflow-hidden">
			<div class="flex items-center justify-between p-6 border-b border-gray-200">
				<h2 class="text-xl font-semibold text-gray-800">Chat Statistics</h2>
				<button onclick="closeStats()" class="text-gray-400 hover:text-gray-600 text-2xl font-bold">&times;</button>
			</div>
			<div class="p-6 overflow-y-auto">
				<div id="statsContent" class="space-y-6">
					<div class="text-center text-gray-500">Loading statistics...</div>
				</div>
			</div>
		</div>
	</div>
}

// Scripts section
templ Scripts() {
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