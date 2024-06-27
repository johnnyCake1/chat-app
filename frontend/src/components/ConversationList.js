import React from 'react';
import { ListGroup } from 'react-bootstrap';
import './ConversationList.css';

const ConversationList = ({ conversations = [], onConverstationSelect }) => {
  return (
    <div className="conversation-list-container">
      <ListGroup>
        {conversations.length > 0 ?
          conversations.map((conversation, idx) => (
            <ListGroup.Item
              key={idx}
              action
              onClick={() => onConverstationSelect(conversation)}
              className={conversation.selected ? 'conversation-item selected' : 'conversation-item'}
            >
              <div className="conversation-info">
                <div className="conversation-title">{conversation.conversationName}</div>
                <div className="conversation-last-message">{conversation.lastMessage}</div>
              </div>
              <div className="conversation-meta">
                <div className="conversation-time">{conversation.timeStamp}</div>
                {conversation.unreadCount > 0 && (
                  <div className="unread-count">{conversation.unreadCount}</div>
                )}
              </div>
            </ListGroup.Item>
          )) : (<div>Start a conversation by searching a user</div>)}
      </ListGroup>
    </div>
  );
};

export default ConversationList;
