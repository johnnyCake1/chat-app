import React, { useEffect, useState } from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import ConversationList from './ConversationList';
import ChatWindow from './ChatWindow';
import SearchBar from './SearchBar';
import { API_URL } from '../constants';
import useLocalStorageState from '../util/userLocalStorage';

const MainChatPage = () => {
  const [selectedConversation, setSelectedConversation] = useState(null);
  const [conversations, setConversations] = useState([]);
  const [token,] = useLocalStorageState('token');
  const [currentUser, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const chatProtocol = 'chat-protocol';

  useEffect(() => {
    if (token) {
      const tokenParts = token.split('.');
      if (tokenParts.length === 3) {
        // Decode the payload (second part)
        const payload = JSON.parse(atob(tokenParts[1]));
        // Extract the user ID from the Issuer payload
        const userId = payload.iss;
        fetch(`${API_URL}/users/${userId}`, { headers: { Authorization: `Bearer ${token}` } })
          .then(response => {
            if (response.ok) {
              response.json().then(userData => {
                console.log("current user:", userData);
                setUser(userData);
                setIsLoading(false);
              })
            }
          })
      } else {
        console.error('Invalid JWT token format');
      }
    }
  }, [token]);

  // Maximum number of connection retry attempts
  const MAX_RETRY_ATTEMPTS = 3;
  let retryCount = 0;
  // Function to create a new WebSocket instance with all necessary handlers
  const createWebSocket = () => {
    if (retryCount >= MAX_RETRY_ATTEMPTS) {
      console.log("Too many reconnect tries!");
      return null;
    }
    const websocket = new WebSocket('ws://localhost/ws', [`${chatProtocol}`, `${token}`]);

    websocket.onopen = () => {
      console.log('WebSocket Connected');
      // Sending pings every 55 seconds so that websocket connection does not die
      const t = setInterval(function () {
        if (websocket.readyState !== 1) {
          clearInterval(t);
          console.log("Websocket connection lost. Trying to recover the connection...");
          // Reconnect if the connection is lost
          setWs(createWebSocket());
        }
        console.log("connection check `PING` sent");
        websocket.send('PING'); // consuming 'ping' string is implemented in the server
      }, 55000);
    };

    websocket.onerror = (error) => {
      console.log('WebSocket Error: ', error, "\n Trying to recover the connection...");
      // Reconnect if the connection is lost (note: it's not really a recursive call as react will re-render this component after calling the setWs so that there's no actual recursive function stack)
      setWs(createWebSocket());
    };

    websocket.onmessage = (e) => {
      const message = JSON.parse(e.data);
      console.log("Message received:", message, "\nCurrent conversations:", conversations);

      setConversations(prevConversations => {
        const updatedConversations = prevConversations.map(conversation => {
          if (conversation.id === message.chatRoomID) {
            return {
              ...conversation,
              lastMessage: message.text,
              timeStamp: new Date(message.timeStamp).toLocaleTimeString(),
              messages: [...conversation.messages, message],
              unreadCount: currentUser.id !== message.senderID ? conversation.unreadCount + 1 : conversation.unreadCount,
            };
          }
          return conversation;
        });

        // Update selected conversation if it matches the received message's chatRoomID
        const selectedConversationIndex = updatedConversations.findIndex(conversation => conversation.id === message.chatRoomID);
        if (selectedConversationIndex !== -1) {
          setSelectedConversation(updatedConversations[selectedConversationIndex]);
        }

        return updatedConversations;
      });
    };
    retryCount++;
    return websocket;
  };

  const [ws, setWs] = useState(null);
  useEffect(() => {
    if (!currentUser) {
      return;
    }

    // Connect to WebSocket
    const websocket = createWebSocket();
    setWs(websocket);

    return () => {
      websocket.close();
    };
  }, [currentUser]);


  const handleConversationSelect = (conversation) => {
    console.log("selected conversation:", conversation);
    setSelectedConversation(conversation);
  };

  const mapClientConversations = (rawConversations) => {
    console.log("raw", rawConversations);
    return rawConversations.map(({ id, isGroup, groupName, participants, messages, unreadCount }) => {
      let conversationName = '';
      if (isGroup) {
        conversationName = groupName;
      } else {
        const otherParticipant = participants.find(participant => participant.id !== currentUser.id);
        if (otherParticipant) {
          conversationName = otherParticipant.nickname || otherParticipant.email;
        }
      }

      let lastMessage = ''
      let timeStamp = ''
      if (messages && messages.length > 0) {
        lastMessage = messages[messages.length - 1].text;
        timeStamp = new Date(messages[messages.length - 1].timeStamp).toLocaleTimeString();
      }
      return {
        id,
        isGroup,
        participants,
        messages,
        lastMessage,
        timeStamp,
        conversationName,
        selected: false,
        unreadCount,
      };
    });
  };


  useEffect(() => {
    if (!currentUser) {
      return
    }
    const fetchChatrooms = async () => {
      try {
        const response = await fetch(`${API_URL}/users/${currentUser.id}/chatrooms`, {
          method: 'GET',
          headers: { Authorization: `Bearer ${token}` },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch chatrooms');
        }

        const data = await response.json();
        const mappedConversations = mapClientConversations(data);
        setConversations(mappedConversations);
      } catch (error) {
        console.error("Failed to fetch chatrooms:", error);
      }
    };

    fetchChatrooms();
  }, [currentUser]);


  if (isLoading) {
    return <div>Loading the user...</div>
  }

  console.log("Conversations:", conversations);

  return (
    <Container fluid>
      <Row>
        <Col md={4}>
          <SearchBar onSelect={handleConversationSelect} />
          <ConversationList conversations={conversations} onConverstationSelect={handleConversationSelect} />
        </Col>
        <Col md={8}>
          <ChatWindow conversation={selectedConversation} ws={ws} currentUser={currentUser} />
        </Col>
      </Row>
    </Container>
  );
};

export default MainChatPage;
