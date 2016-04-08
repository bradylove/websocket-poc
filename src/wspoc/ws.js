var wsHost = "ws://localhost:9292/stream/"

$(document).ready(function() {
  client = new Client();

  client.connect()
})

var Client = function(streamId) {
  if (streamId == undefined) {
    streamId = "";
  }

  return {
    streamId: streamId,
    conn: null,

    connect: function() {
      this.conn = new WebSocket(wsHost + this.streamId)
      this.conn.onmessage = this.onmessage
      this.conn.error = this.onerror
      this.conn.onclose = function() {
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
        this.streamId = data.id
        console.log("StreamID:" + this.streamId)

        return
      }

      $("#messages").append("" + data.count + "\n")
    }
  }
}

