import React, { useEffect, useRef, useState } from 'react';
import { Container, Row, Col, Form, Button } from 'react-bootstrap';
import './ChatWindow.css';
import { Navigate } from 'react-router-dom';
import { API_URL, MessageOptions } from '../constants';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheckDouble } from '@fortawesome/free-solid-svg-icons';
import { faCheck } from '@fortawesome/free-solid-svg-icons';
import useLocalStorageState from '../util/userLocalStorage';

const ChatWindow = ({ conversation, setConversation, ws, currentUser }) => {
  const [messageInput, setMessageInput] = useState('');
  const [page, setPage] = useState(1);
  const [token,] = useLocalStorageState('token');

  const messageRefs = useRef([]);

  useEffect(() => {
    if (!conversation) {
      return;
    }
    const markMessageAsViewed = (messageId) => {
      console.log("Marking message as viewed:", messageId);
      // convert messageId string to uint
      messageId = parseInt(messageId);
      if (!ws) {
        console.error("No connection with the server");
        return;
      }
      console.log("convo id:", conversation.id, "message id:", messageId, "viewer id:", currentUser.id)
      ws.send(JSON.stringify({
        messageOption: MessageOptions.VIEW_MESSAGE,
        viewMessage: {
          chatroomID: conversation.id,
          messageID: messageId,
          viewerID: currentUser.id
        }
      }));
      console.log("Marking message as viewed:", messageId);
    };
    // Clear the messageRefs array when the conversation changes
    messageRefs.current = messageRefs.current.slice(0, conversation.messages.length);

    const observer = new IntersectionObserver(
      entries => {
        entries.forEach(entry => {
          if (entry.isIntersecting) {
            // Mark the message as viewed when it enters the viewport
            const messageID = entry.target.getAttribute('data-id');
            markMessageAsViewed(messageID);
          }
        });
      },
      { threshold: 1 }
    );

    messageRefs.current.forEach(ref => {
      if (ref) {
        observer.observe(ref);
      }
    });

    return () => {
      messageRefs.current.forEach(ref => {
        if (ref) {
          observer.unobserve(ref);
        }
      });
    };
  }, [conversation, currentUser, ws]);

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
      // Filter out the messages that already exist in the conversation to make sure we don't add duplicates
      const newMessages = responseData.filter(
        (newMessage) => !conversation.messages.some(
          (existingMessage) => existingMessage.id === newMessage.id
        )
      );
      // Append the new messages to the start of the conversation
      conversation.messages = [...newMessages, ...conversation.messages];
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
      const messageData = {
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

  if (!currentUser) {
    return <Navigate to='/login' />
  }

  if (!ws) {
    return <div>No connection with the server</div>
  }

  if (!conversation) {
    return (
      <Container className="chat-window">
        <div className="empty-chat">Select a conversation to start chatting</div>
      </Container>
    );
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
      {/* <Row className="message-list" ref={messageListRef}> */}
      <Row className="message-list">
        <Col>
          {(conversation.messages.length > 0 && conversation.messages.map((message, idx) => (
            <div
              key={idx}
              className={`message ${message.senderID === currentUser.id ? 'sent' : 'received'}`}
              // only set the ref if the message hasn't been viewed yet so that the observer is triggered only once and only for the messages that haven't been viewed
              ref={el => !message.viewed && message.senderID !== currentUser.id ? messageRefs.current[idx] = el : null}
              data-id={message.id}
            >
              <div className="message-content">{message.text}</div>
              <div className="message-timestamp">
                {new Date(message.timeStamp).toLocaleTimeString()}
              </div>
              <div className="message-viewed">
                <FontAwesomeIcon icon={message.viewed ? faCheckDouble : faCheck} />
              </div>
            </div>
          )))
            ||
            <div className="empty-chat">No messages yet. Say hello to {conversation.conversationName}</div>
          }
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
