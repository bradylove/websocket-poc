var wsHost = "ws://localhost:9292/stream/"

$(document).ready(function() {
  client = new Client();

  client.connect()
})

var Client = function(streamId) {
  return {
    streamId: streamId,
    conn: null,

    connect: function() {
      this.conn = new WebSocket(wsHost + this.streamId)
      this.conn.onclose = this.onclose
      this.conn.onmessage = this.onmessage
    },


    onclose: function(event) {
      console.log("close")
    },

    onmessage: function(event) {
      console.log("message")
      console.log(event.data)
    }
  }
}

