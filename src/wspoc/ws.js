var wsHost = "wss://wspoc.cfapps.io:4443/stream/"
// var wsHost = "ws://localhost:9292/stream/"
var streamId = ""

$(document).ready(function() {
  client = new Client();

  client.connect()
})

var Client = function() {
  return {
    conn: null,

    connect: function() {
      this.conn = new WebSocket(wsHost + streamId)
      this.conn.onmessage = this.onmessage
      this.conn.error = this.onerror
      this.conn.onclose = function(event) {
        console.log(event)

        setTimeout(function() {
          console.log("Reconnecting...")
          this.connect()
        }.bind(this), 1000)
      }.bind(this)
    },

    onerror: function(event) {
      console.log(event)
    },

    onclose: function(event) {
      this.connect()
    },

    onmessage: function(event) {
      var data = JSON.parse(event.data)

      if (data.id != null) {
        streamId = data.id
        console.log("StreamID:" + this.streamId)

        return
      }

      $("#messages").append("" + data.count + "\n")
    }
  }
}

