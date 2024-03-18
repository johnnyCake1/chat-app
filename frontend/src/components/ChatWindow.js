import React, { useEffect, useRef, useState } from 'react';
import { Container, Row, Col, Form, Button } from 'react-bootstrap';
import './ChatWindow.css';
import { Navigate } from 'react-router-dom';
import { API_URL, MessageOptions } from '../constants';
import useLocalStorageState from '../util/userLocalStorage';

const ChatWindow = ({ conversation, setConversation /* supposed to be used to update the passed conversation */, ws, currentUser }) => {
  const [messageInput, setMessageInput] = useState('');
  const [page, setPage] = useState(1);
  const [token,] = useLocalStorageState('token');

  const messageListRef = useRef(null);

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

  const loadMessages = async (page) => {
    fetch(`${API_URL}/chatrooms/${conversation.id}/messages?page=${page}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      }
    }).then(response => {
      if (response.ok) {
        return response.json();
      }
      console.log("Failed to load messages:", response.status, response.statusText);
      return [];
    }).then(responseData => {
      console.log("Loaded messages:", responseData);
      // use setConversation to update the conversation with the new messages instead of this:
      conversation.messages = [...responseData, ...conversation.messages];
      setConversation(conversation);
      setPage(prevPage => prevPage + 1);
    }).catch(err => {
      console.error("Error loading messages:", err);
    });
  };

  const handleScroll = (e) => {
    const { scrollTop } = e.currentTarget;
    if (scrollTop === 0) {
      loadMessages(page);
    }
  };

  const handleSendMessage = (e) => {
    e.preventDefault();
    if (ws && messageInput.trim()) {
      setMessageInput('');
      let messageData = {
        messageOption: '',
      };

      if (conversation.id) {
        // If conversation exists, send message to the conversation
        messageData.messageOption = MessageOptions.SEND_MESSAGE;
        messageData.sendMessage = {
          senderID: currentUser.id,
          chatroomID: conversation.id,
          text: messageInput,
        };
      } else {
        // If conversation doesn't exist, create a new private chatroom
        messageData.messageOption = MessageOptions.CREATE_PRIVATE_CHATROOM;
        messageData.createPrivateChatroom = {
          participants: conversation.participants,
          chatMessage: {
            text: messageInput,
            senderID: currentUser.id,
          }
        };
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
    <Container className="chat-window" onScroll={handleScroll}>
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
          {conversation.messages.length > 0 && conversation.messages.map((message, idx) => (
            <div
              key={idx}
              className={`message ${message.senderID === currentUser.id ? 'sent' : 'received'}`}
              ref={idx === conversation.messages.length - 1 ? lastMessageRef : null}
            >
              <div className="message-content">{message.text}</div>
              <div className="message-timestamp">
                {new Date(message.timeStamp).toLocaleTimeString()}
              </div>
            </div>
          )) || <div className="empty-chat">No messages yet. Say hello to {conversation.conversationName}</div>}
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
