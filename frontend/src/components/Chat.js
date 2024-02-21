import React, { useEffect, useState } from 'react';

function Chat() {
    const [ws, setWs] = useState(null);
    const [messages, setMessages] = useState([]);
    const [input, setInput] = useState('');

    useEffect(() => {
        // Connect to WebSocket
        const websocket = new WebSocket('ws://localhost/ws'); // Use your actual host here
        setWs(websocket);

        websocket.onopen = () => console.log('WebSocket Connected');
        websocket.onerror = (error) => console.log('WebSocket Error: ', error);
        websocket.onmessage = (e) => {
            const message = JSON.parse(e.data);
            console.log("Message received:", message)
            setMessages(prev => [...prev, message]);
        };

        return () => {
            websocket.close();
        };
    }, []);

    const sendMessage = () => {
        if (ws && input.trim()) {
            ws.send(JSON.stringify({ text: input }));
            console.log("Message sent:", input)
            setInput(''); // Clear input after send
        }
    };

    return (
        <div>
            <ul>
                {messages.map((msg, index) => (
                    <li key={index}>{msg.text}</li>
                ))}
            </ul>
            <input
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' ? sendMessage() : null}
            />
            <button onClick={sendMessage}>Send Message</button>
        </div>
    );
}

export default Chat;
