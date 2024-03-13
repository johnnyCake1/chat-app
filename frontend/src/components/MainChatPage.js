import React, { useEffect, useState } from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import ConversationList from './ConversationList';
import ChatWindow from './ChatWindow';
import SearchBar from './SearchBar';
import { API_URL } from '../constants';
import useLocalStorageState from '../util/userLocalStorage';

const MainChatPage = ({ currentUser }) => {
  const currentUserID = 1;
  const [selectedConversation, setSelectedConversation] = useState(null);
  const [conversations, setConversations] = useState([]);
  const [token,] = useLocalStorageState('token');

  const handleConversationSelect = (conversation) => {
    setSelectedConversation(conversation);
  };

  const mapClientConversations = (rawConversations, currentUserID) => {
    return rawConversations.map(({ id, isGroup, groupName, participants, messages }) => {
      let conversationName = '';
      if (isGroup) {
        conversationName = groupName;
      } else {
        const otherParticipant = participants.find(participant => participant.id !== currentUserID);
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
      };
    });
  };


  useEffect(() => {
    const fetchChatrooms = async () => {
      try {
        const response = await fetch(`${API_URL}/users/${currentUserID}/chatrooms`, {
          method: 'GET',
          headers: { Authorization: `Bearer ${token}` },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch chatrooms');
        }

        const data = await response.json();
        const mappedConversations = mapClientConversations(data, currentUserID);
        setConversations(mappedConversations);
      } catch (error) {
        console.error("Failed to fetch chatrooms:", error);
      }
    };

    fetchChatrooms();
  }, []);


  return (
    <Container fluid>
      <Row>
        <Col md={4}>
          <SearchBar onSelect={handleConversationSelect} />
          <ConversationList conversations={conversations} onConverstationSelect={handleConversationSelect} />
        </Col>
        <Col md={8}>
          <ChatWindow conversation={selectedConversation} />
        </Col>
      </Row>
    </Container>
  );
};

export default MainChatPage;
