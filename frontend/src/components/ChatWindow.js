import React, { useEffect, useRef, useState } from 'react';
import { Container, Row, Col, Form, Button } from 'react-bootstrap';
import './ChatWindow.css';
import { Navigate } from 'react-router-dom';

const ChatWindow = ({ conversation, ws, currentUser }) => {
  const [messageInput, setMessageInput] = useState('');
  const messageListRef = useRef(null);

  // enums for message options
  const MessageOptions = {
    MESSAGE_SEND: "MESSAGE_SEND",
    MESSAGE_EDIT: "MESSAGE_EDIT",
    MESSAGE_VIEW: "MESSAGE_VIEW",
    MESSAGE_DELETE: "MESSAGE_DELETE",
    MESSAGE_REACTION: "MESSAGE_REACTION",
  };

  // scroll down whenever user sends/receives message
  const lastMessageRef = useRef(null);
  useEffect(() => {
    if (lastMessageRef.current) {
      lastMessageRef.current.scrollIntoView();
    }
  }, [conversation]);

  const handleInputChange = (e) => {
    setMessageInput(e.target.value);
  };

  


  const handleSendMessage = (e) => {
    e.preventDefault();
    if (ws && messageInput.trim()) {
      const messageData = {
        messageOption: MessageOptions.MESSAGE_SEND,
        chatroomID: conversation.id,
        senderID: currentUser.id,
        text: messageInput,
      }
      console.log("Sending message data:", messageData)
      ws.send(JSON.stringify(messageData));
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

  if (!ws) {
    return <div>No connection with the server</div>
  }

  if (!currentUser) {
    return <Navigate to='/login' />
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
      <Row className="message-list" ref={messageListRef}>
        <Col>
          {conversation.messages.map((message, idx) => (
            <div
              key={message.id}
              className={`message ${message.senderID === currentUser.id ? 'sent' : 'received'}`}
              ref={idx === conversation.messages.length - 1 ? lastMessageRef : null}
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
