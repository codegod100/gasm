package main

import (
	"syscall/js"
)

func saveMessages() {
	document := js.Global().Get("document")
	messagesElement := document.Call("getElementById", "messages")
	if messagesElement.IsNull() {
		return
	}

	html := messagesElement.Get("innerHTML")
	js.Global().Get("localStorage").Call("setItem", "chatHTML", html)
}

func loadMessages() {
	document := js.Global().Get("document")
	messagesElement := document.Call("getElementById", "messages")
	if messagesElement.IsNull() {
		return
	}

	savedHTML := js.Global().Get("localStorage").Call("getItem", "chatHTML")
	if !savedHTML.IsNull() {
		messagesElement.Set("innerHTML", savedHTML)
		messagesElement.Set("scrollTop", messagesElement.Get("scrollHeight"))
	}
}

func addMessage() {
	document := js.Global().Get("document")
	usernameInput := document.Call("getElementById", "username")
	messageInput := document.Call("getElementById", "messageInput")
	messagesElement := document.Call("getElementById", "messages")

	if usernameInput.IsNull() || messageInput.IsNull() || messagesElement.IsNull() {
		return
	}

	username := usernameInput.Get("value")
	message := messageInput.Get("value")

	if message.IsNull() {
		return
	}

	messageDiv := document.Call("createElement", "div")
	messageDiv.Set("className", "message")

	userDiv := document.Call("createElement", "div")
	userDiv.Set("className", "message-user")
	if username.IsNull() {
		userDiv.Set("textContent", "Anonymous")
	} else {
		userDiv.Set("textContent", username)
	}

	textDiv := document.Call("createElement", "div")
	textDiv.Set("className", "message-text")
	textDiv.Set("textContent", message)

	messageDiv.Call("appendChild", userDiv)
	messageDiv.Call("appendChild", textDiv)
	messagesElement.Call("appendChild", messageDiv)

	messageInput.Set("value", "")
	messageInput.Call("focus")
	messagesElement.Set("scrollTop", messagesElement.Get("scrollHeight"))

	saveMessages()
}

func clearMessages() {
	document := js.Global().Get("document")
	messagesElement := document.Call("getElementById", "messages")
	if !messagesElement.IsNull() {
		messagesElement.Set("innerHTML", "")
		saveMessages()
	}
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

	js.Global().Set("loadMessages", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		loadMessages()
		return nil
	}))

	<-c
}
