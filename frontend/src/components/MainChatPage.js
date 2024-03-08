import React, { useState } from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import ConversationList from './ConversationList';
import ChatWindow from './ChatWindow';
import SearchBar from './SearchBar';

const MainChatPage = () => {
  const [selectedConversation, setSelectedConversation] = useState(null);

  const handleConversationSelect = (conversationId) => {
    setSelectedConversation(conversationId);
  };

  return (
    <Container fluid>
      <Row>
        <Col md={4}>
          <SearchBar />
          <ConversationList onSelect={handleConversationSelect} />
        </Col>
        <Col md={8}>
          <ChatWindow conversationId={selectedConversation} />
        </Col>
      </Row>
    </Container>
  );
};

export default MainChatPage;
