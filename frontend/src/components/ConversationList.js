import React from 'react';
import { ListGroup } from 'react-bootstrap';

const ConversationList = ({ onSelect }) => {
  const conversations = [];

  return <>Number of chats: {conversations.length}</>
  return (
    <ListGroup className="conversation-list">
      {conversations.map(conversation => (
        <ListGroup.Item key={conversation.id} onClick={() => onSelect(conversation.id)}>
          <div className="avatar">{conversation.avatar}</div>
          <div className="info">
            <div className="name">{conversation.name}</div>
            <div className="message">{conversation.lastMessage}</div>
          </div>
        </ListGroup.Item>
      ))}
    </ListGroup>
  );
};

export default ConversationList;
