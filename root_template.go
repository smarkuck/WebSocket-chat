package main

import "html/template"

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<script>
	websocket = new WebSocket("ws://{{.}}/socket")
	websocket.onmessage = onMessage
	websocket.onclose = onClose

	function onMessage(m) {
		insertLine(m.data)
	}

	function onClose() {
		insertLine("Connection closed.")
	}

	function send(e) {
		if(event.key === 'Enter') {
			websocket.send(e.value + "\n")
			insertLine(e.value)
			e.value = ""
		}
	}

	function insertLine(text) {
		document.body.appendChild(document.createElement("br"))
		document.body.appendChild(document.createTextNode(text))
	}
</script>
</head>
<body>
<input type="text" onkeydown="send(this)" /><br />
</body>
</html>
`))
