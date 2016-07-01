// dummy modbus responser

var 
zmq = require('zmq')
, pub = zmq.socket('pub')
, sub = zmq.socket('sub')
, ipc_pub = "ipc:///tmp/from.modbus"
, ipc_sub = "ipc:///tmp/to.modbus"


pub.bindSync(ipc_pub); // connect to zmq endpoint
sub.bindSync(ipc_sub); // bind to zmq endpoint
sub.subscribe(""); // filter topic

var cmd = {
    "tid": 1,
    "data": 1234,
    "status": "ok"
}

// start listening response
sub.on("message", function(mode, jstr) {
    //console.log(mode.toString());
    console.log(jstr.toString());
    pub.send("tcp", zmq.ZMQ_SNDMORE);
    pub.send(JSON.stringify(cmd));
    console.log("response");
});