const BACKEND_HOST = 'localhost';

export const API_URL = `http://${BACKEND_HOST}:8080/api/v1`; //docker frontend uses NGNIX to proxy requests to `/` to backend
export const WEB_SOCKET_URL = `ws://${BACKEND_HOST}:8080/ws`; //docker frontend uses NGNIX to proxy requests to `/` to backend

export const CHAT_SUBPROTOCOL = 'chat-protocol';

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