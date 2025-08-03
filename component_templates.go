package main

import "strconv"

// User card component
templ UserCard(username string, messageCount int) {
	<div class="flex items-center bg-white border border-gray-300 rounded p-1 min-w-20 shadow-sm text-xs">
		<div class="w-5 h-5 rounded-full bg-gradient-to-br from-chat-green to-chat-green-hover text-white flex items-center justify-center font-bold mr-1 uppercase text-xs">
			if len(username) > 0 {
				{ string(username[0]) }
			} else {
				?
			}
		</div>
		<div class="flex-1">
			<div class="font-bold text-gray-800 mb-1">{ username }</div>
			<div class="text-sm text-gray-600">{ strconv.Itoa(messageCount) } messages</div>
		</div>
	</div>
}

// Notification component
templ Notification(message, notificationType string) {
	<div class={ "notification", "notification-" + notificationType }>
		<div class="text-lg mr-3 font-bold">
			switch notificationType {
				case "success":
					✓
				case "error":
					✗
				case "info":
					ℹ
				default:
					!
			}
		</div>
		<div class="flex-1">{ message }</div>
		<button class="bg-none border-none text-lg cursor-pointer ml-3 p-0 w-5 h-5 flex items-center justify-center opacity-70 hover:opacity-100" onclick="this.parentElement.remove()">×</button>
	</div>
}

// Message component
templ MessageComponent(username, text, timestamp string) {
	<div class="mb-4 p-2 rounded bg-gray-50">
		<div class="font-bold text-chat-green mb-1">{ username }</div>
		<div class="text-gray-800">{ text }</div>
		<div class="text-xs text-gray-600 mt-1">{ timestamp }</div>
	</div>
}