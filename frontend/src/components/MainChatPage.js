import React, { useEffect, useState } from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import ConversationList from './ConversationList';
import ChatWindow from './ChatWindow';
import SearchBar from './SearchBar';
import { API_URL, CHAT_SUB_PROTOCOL, MessageOptions } from '../constants';
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

  // Update the selected conversation when the conversations list changes
  useEffect(() => {
    if (selectedConversation) {
      const updatedSelectedConversation = conversations.find(conversation => conversation.id === selectedConversation.id);
      if (updatedSelectedConversation) {
        setSelectedConversation(updatedSelectedConversation);
      }
    }
  }, [conversations, selectedConversation]);

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
        return [newConversation, ...prevConversations];
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
        return updatedConversations;
      }
    });
  };

  const handleCreatePrivateChatroom = (messageData) => {
    setConversations(prevConversations => {
      // Check if the conversation exists
      const conversationExists = prevConversations.some(conversation => conversation.id === messageData.createPrivateChatroom.id);
      console.log("Conversation exists:", conversationExists);
      if (!conversationExists) {
        // If the conversation doesn't exist, create a new one
        const newConversation = {
          id: messageData.createPrivateChatroom.id,
          messages: [messageData.createPrivateChatroom.chatMessage],
          lastMessage: messageData.createPrivateChatroom.chatMessage.text,
          profilePictureURL: messageData.createPrivateChatroom.chatroomPictureURL,
          conversationName: messageData.createPrivateChatroom.chatroomName,
          timeStamp: new Date(messageData.createPrivateChatroom.chatMessage.timeStamp).toLocaleTimeString(),
          unreadCount: messageData.createPrivateChatroom.unreadCount
        };
        // Add the new conversation to the conversations array
        return [newConversation, ...prevConversations];
      }
      // If the conversation exists, update it
      const updatedConversations = prevConversations.map(conversation => {
        if (conversation.id === messageData.createPrivateChatroom.id) {
          return {
            ...conversation,
            lastMessage: messageData.createPrivateChatroom.chatMessage.text,
            timeStamp: new Date(messageData.createPrivateChatroom.chatMessage.timeStamp).toLocaleTimeString(),
            messages: [...conversation.messages, messageData.createPrivateChatroom.chatMessage],
            conversationName: messageData.createPrivateChatroom.chatroomName,
            unreadCount: messageData.createPrivateChatroom.unreadCount,
          };
        }
        return conversation;
      });

      return updatedConversations;
    });
  };

  const handleMarkMessageAsViewed = (messageData) => {
    console.log("Marking message as viewed:", messageData);
    setConversations(prevConversations => {
      const updatedConversations = prevConversations.map(conversation => {
        if (conversation.id === messageData.viewMessage.chatroomID) {
          const updatedConversation = {
            ...conversation,
            messages: conversation.messages.map(message => {
              if (message.id === messageData.viewMessage.messageID) {
                return messageData.viewMessage.chatMessage;
              }
              return message;
            }),
            unreadCount: Math.max(0, conversation.unreadCount - 1),
          };
          console.log("Updating conversation:", updatedConversation, "\nSelected conversation:", selectedConversation);
          return updatedConversation;
        }
        return conversation;
      });
      console.log("updating conversations:", updatedConversations);

      return updatedConversations;
    });
  };
  // Function to update the conversation list with the new conversation and return the updated list
  const setConversation = (newConversation) => {
    setConversations((prevConversations) => {
      let found = false;
      const updatedConversations = prevConversations.map((conversation) => {
        if (conversation.id === newConversation.id) {
          found = true;
          return newConversation;
        }
        return conversation;
      });
      // If the conversation is not found, add it to the beginning of the list
      if (!found) {
        return [newConversation, ...updatedConversations];
      }
      return updatedConversations;
    });
  };

  // Retry timeout for WebSocket connection which increases after each retry
  const initialRetryTimeout = 1000; // start with 1 second
  const maxRetryTimeout = 60000; // max 60 seconds
  let retryTimeout = initialRetryTimeout;
  // Function to create a new WebSocket instance with all necessary handlers
  const createWebSocket = () => {
    console.log("Creating a new WebSocket instance...");
    const websocket = new WebSocket('ws://localhost/ws', [`${CHAT_SUB_PROTOCOL}`, `${token}`]);

    websocket.onopen = () => {
      console.log('WebSocket Connected');
      retryTimeout = initialRetryTimeout; // reset reconnect timeout on successful connection
      // Sending heartbeat check pings every 55 seconds so that websocket connection does not die
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
    };

    websocket.onclose = (event) => {
      console.log(`WebSocket closed with code ${event.code} and reason: ${event.reason}`);
      console.log("Trying to recover the connection...");
      setTimeout(() => {
        setWs(createWebSocket());
        // Increase the timeout after each retry
        retryTimeout *= 2;
        if (retryTimeout > maxRetryTimeout) {
          retryTimeout = maxRetryTimeout;
        }
      }, retryTimeout);
    };

    websocket.onmessage = (e) => {
      const messageData = JSON.parse(e.data);
      console.log(`Message data received with option ${messageData.messageOption}:`, messageData);
      // decide what to do with the received message data
      switch (messageData.messageOption) {
        case MessageOptions.SEND_MESSAGE: {
          handleReceivedMessage(messageData);
          break;
        }
        case MessageOptions.CREATE_PRIVATE_CHATROOM: {
          handleCreatePrivateChatroom(messageData);
          break;
        }
        case MessageOptions.VIEW_MESSAGE: {
          handleMarkMessageAsViewed(messageData);
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

  // resolve additional conversation details // TODO: remove this client side logic and make the api return the additional details
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
      const emptyConversation = resolveConversationAdditionalDetails({
        participants: [currentUser, selectedUser],
        messages: [],
      })
      console.log("Creating new conversation:", emptyConversation);
      // Select the new conversation
      setSelectedConversation(emptyConversation);
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
          <ChatWindow conversation={selectedConversation} setConversation={setConversation} ws={ws} currentUser={currentUser} />
        </Col>
      </Row>
    </Container>
  );
};

export default MainChatPage;
