import React, { useEffect, useRef, useState } from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import ConversationList from './ConversationList';
import ChatWindow from './ChatWindow';
import SearchBar from './SearchBar';
import { API_URL, CHAT_SUB_PROTOCOL } from '../constants';
import useLocalStorageState from '../util/userLocalStorage';

const MainChatPage = () => {
  const [selectedConversation, setSelectedConversation] = useState(null);
  const [conversations, setConversations] = useState([]);
  const [token,] = useLocalStorageState('token');
  const [currentUser, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

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

  const handleReceivedMessage = (messageData) => {
    setConversations(prevConversations => {
      // Check if the conversation exists
      const conversationExists = prevConversations.some(conversation => conversation.id === messageData.sendMessage.chatroomID);
      console.log("Conversation exists:", conversationExists);
      if (!conversationExists) {
        // If the conversation doesn't exist, create a new one
        const newConversation = {
          id: messageData.sendMessage.chatroomID,
          messages: [messageData.sendMessage],
          lastMessage: messageData.sendMessage.text,
          timeStamp: new Date(messageData.sendMessage.timeStamp).toLocaleTimeString(),
          selected: false,
          unreadCount: currentUser.id !== messageData.sendMessage.senderID ? 1 : 0,
        };

        // Add the new conversation to the conversations array
        return [...prevConversations, newConversation];
      } else {
        // If the conversation exists, update it
        const updatedConversations = prevConversations.map(conversation => {
          if (conversation.id === messageData.sendMessage.chatroomID) {
            return {
              ...conversation,
              lastMessage: messageData.sendMessage.text,
              timeStamp: new Date(messageData.sendMessage.timeStamp).toLocaleTimeString(),
              messages: [...conversation.messages, messageData.sendMessage],
              unreadCount: currentUser.id !== messageData.sendMessage.senderID ? conversation.unreadCount + 1 : conversation.unreadCount,
            };
          }
          return conversation;
        });

        // Update selected conversation if it matches the received message's chatroomID
        const selectedConversationIndex = updatedConversations.findIndex(conversation => conversation.id === messageData.sendMessage.chatroomID);
        if (selectedConversationIndex !== -1) {
          setSelectedConversation(updatedConversations[selectedConversationIndex]);
        }

        return updatedConversations;
      }
    });
  };

  const handleCreatePrivateChatroom = (messageData) => {
    setConversations(prevConversations => {
      // Check if the conversation exists
      const conversationExists = prevConversations.some(conversation => conversation.id === messageData.createPrivateChatroom.chatroomID);
      console.log("Conversation exists:", conversationExists);
      if (!conversationExists) {
        // If the conversation doesn't exist, create a new one
        const newConversation = {
          id: messageData.createPrivateChatroom.chatroomID,
          messages: [messageData.createPrivateChatroom.chatMessage],
          lastMessage: messageData.createPrivateChatroom.chatMessage.text,
          profilePictureURL: messageData.createPrivateChatroom.chatroomPictureURL,
          conversationName: messageData.createPrivateChatroom.chatroomName,
          timeStamp: messageData.createPrivateChatroom.chatMessage.timeStamp,
          unreadCount: messageData.createPrivateChatroom.unreadCount
        };

        // Add the new conversation to the conversations array
        return [...prevConversations, newConversation];
      } else {
        // If the conversation exists, update it
        const updatedConversations = prevConversations.map(conversation => {
          if (conversation.id === messageData.createPrivateChatroom.chatroomID) {
            return {
              ...conversation,
              lastMessage: messageData.createPrivateChatroom.chatMessage.text,
              timeStamp: messageData.createPrivateChatroom.chatMessage.timeStamp,
              messages: [...conversation.messages, messageData.createPrivateChatroom.chatMessage],
              conversationName: messageData.createPrivateChatroom.chatroomName,
              unreadCount: messageData.createPrivateChatroom.unreadCount,
            };
          }
          return conversation;
        });

        // Update selected conversation if it matches the received message's chatroomID
        const selectedConversationIndex = updatedConversations.findIndex(conversation => conversation.id === messageData.createPrivateChatroom.chatroomID);
        if (selectedConversationIndex !== -1) {
          setSelectedConversation(updatedConversations[selectedConversationIndex]);
        }

        return updatedConversations;
      }
    });
  };

  // Function to create a new WebSocket instance with all necessary handlers
  const createWebSocket = () => {
    console.log("Creating a new WebSocket instance...");
    const websocket = new WebSocket('ws://localhost/ws', [`${CHAT_SUB_PROTOCOL}`, `${token}`]);

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
      const messageData = JSON.parse(e.data);
      console.log(`Message data received with option ${messageData.messageOption}:`, messageData);
      // decide what to do with the received message data
      switch (messageData.messageOption) {
        case 'SEND_MESSAGE': {
          handleReceivedMessage(messageData);
          break;
        }
        case 'CREATE_PRIVATE_CHATROOM': {
          handleCreatePrivateChatroom(messageData);
          break;
        }
        default: {
          console.log("Unknown message option:", messageData.messageOption);
        }
      }
    };
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

  // resolve additional conversation details
  const resolveConversationAdditionalDetails = ({ id, isGroup, groupName, participants, messages, unreadCount }) => {
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
      groupName,
      participants,
      messages,
      unreadCount,
      // Additional details after resolving:
      lastMessage,
      timeStamp,
      conversationName,
      selected: false,
    };
  }


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
        console.log("Fetched chatrooms:", data);
        setConversations(data.map(resolveConversationAdditionalDetails));
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

  const handleUserSelect = (selectedUser) => {
    console.log("Selected user:", selectedUser);
    console.log("current conversations:", conversations)
    // Find the 1 to 1 conversation with the selected user
    const existingConversation = conversations.find(conversation =>
      !conversation.isGroup &&
      conversation.participants.some(participant => participant.id === selectedUser.id)
    );

    if (existingConversation) {
      // If a conversation with the selected user exists, select it
      console.log("Existing conversation found:", existingConversation);
      setSelectedConversation(existingConversation);
    } else {
      // If no conversation exists, create a new one
      const newConversation = resolveConversationAdditionalDetails({
        participants: [currentUser, selectedUser],
        messages: [],
      })
      console.log("Creating new conversation:", newConversation);

      // Add the new conversation to the conversations array, but don't persist it to the server yet
      setConversations(prevConversations => [...prevConversations, newConversation]);
      // Select the new conversation
      setSelectedConversation(newConversation);
    }
  };
  return (
    <Container fluid>
      <Row>
        <Col md={4}>
          <SearchBar onSelect={handleUserSelect} />
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
