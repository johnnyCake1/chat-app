import React from 'react';
import { Card } from 'react-bootstrap';

const ChatWindow = ({ conversation }) => {
  return (
    <Card className="chat-window">
      <Card.Body>
        {conversation ? (
          <div className="conversation">
            {conversation.messages.map(message => (
              <div key={message.id} className="message">
                <div className="sender">{message.sender}</div>
                <div className="text">{message.text}</div>
              </div>
            ))}
          </div>
        ) : (
          <div className="welcome-message">Select a conversation to start chatting</div>
        )}
      </Card.Body>
    </Card>
  );
};

export default ChatWindow;
