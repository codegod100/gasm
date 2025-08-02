package main

import "strconv"

// User card component
templ UserCard(username string, messageCount int) {
	<div class="user-card">
		<div class="user-avatar">
			if len(username) > 0 {
				{ string(username[0]) }
			} else {
				?
			}
		</div>
		<div class="user-info">
			<div class="user-name">{ username }</div>
			<div class="user-count">{ strconv.Itoa(messageCount) } messages</div>
		</div>
	</div>
}

// Notification component
templ Notification(message, notificationType string) {
	<div class={ "notification", "notification-" + notificationType }>
		<div class="notification-icon">
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
		<div class="notification-text">{ message }</div>
		<button class="notification-close" onclick="this.parentElement.remove()">×</button>
	</div>
}

// Message component
templ MessageComponent(username, text, timestamp string) {
	<div class="message">
		<div class="message-user">{ username }</div>
		<div class="message-text">{ text }</div>
		<div class="message-time">{ timestamp }</div>
	</div>
}