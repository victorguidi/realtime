https://localhost:3000/

let socket = new WebSocket("wss://localhost:3000/wss/login", ["1"]);
socket.onopen = (event) => { socket.send(JSON.stringify({
    sessionToken: "tes",
    id: 2
})) };
socket.onmessage = (event) => console.log(event.data);

socket.send(JSON.stringify({
    from: 2,
    message: "Hellloo"
}))
