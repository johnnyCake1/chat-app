import React, { useState } from 'react';
import { Container, Row, Col, Form, Button } from 'react-bootstrap';
import './ChatWindow.css';
import useLocalStorageState from '../util/userLocalStorage';

const ChatWindow = ({ conversation }) => {
  const [messageInput, setMessageInput] = useState('');

  const handleInputChange = (e) => {
    setMessageInput(e.target.value);
  };

  const handleSendMessage = (e) => {
    e.preventDefault();
    if (messageInput.trim() !== '') {
      // Send the message logic here
      console.log('Sending message:', messageInput);
      // Clear input after sending message
      setMessageInput('');
    }
  };

  if (!conversation) {
    return (
      <Container className="chat-window">
        <div className="empty-chat">Select a conversation to start chatting</div>
      </Container>
    );
  }

  return (
    <Container className="chat-window">
      <Row className="header-bar">
        <Col>
          <div className="conversation-name">{conversation.conversationName}</div>
        </Col>
        <Col className="action-buttons">
          <Button variant="transparent" className="action-button" size='lg'>
            ...
          </Button>
        </Col>
      </Row>
      <Row className="message-list">
        <Col>
          {conversation.messages.map((message) => (
            <div
              key={message.id}
              className={`message ${message.senderID === 1 ? 'sent' : 'received'}`}
            >
              <div className="message-content">{message.text}</div>
              <div className="message-timestamp">
                {new Date(message.timeStamp).toLocaleTimeString()}
              </div>
            </div>
          ))}
        </Col>
      </Row>
      <Row className="input-box">
        <Col>
          <Form onSubmit={handleSendMessage}>
            <Form.Group controlId="messageInput">
              <Form.Control
                type="text"
                autoComplete="off"
                placeholder="Type a message..."
                value={messageInput}
                onChange={handleInputChange}
              />
            </Form.Group>
            <Button variant="primary" type="submit">
              Send
            </Button>
          </Form>
        </Col>
      </Row>
    </Container>
  );
};

export default ChatWindow;
