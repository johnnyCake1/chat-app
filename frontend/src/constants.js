const HOST = 'localhost';

export const API_URL = `http://${HOST}:8080/api/v1`;

export const CHAT_SUB_PROTOCOL = 'chat-protocol';

// enums for message options
export const MessageOptions = {
    SEND_MESSAGE: "SEND_MESSAGE",
    VIEW_MESSAGE: "VIEW_MESSAGE",
    EDIT_MESSAGE: "EDIT_MESSAGE",
    DELETE_MESSAGE: "DELETE_MESSAGE",
    REACT_TO_MESSAGE: "REACT_TO_MESSAGE",
    CREATE_PRIVATE_CHATROOM: "CREATE_PRIVATE_CHATROOM",
    UPDATE_PRIVATE_CHATROOM: "UPDATE_PRIVATE_CHATROOM",
    CREATE_GROUP_CHATROOM: "CREATE_GROUP_CHATROOM",
    UPDATE_GROUP_CHATROOM: "UPDATE_GROUP_CHATROOM",
    DELETE_GROUP_CHATROOM: "DELETE_GROUP_CHATROOM",
};