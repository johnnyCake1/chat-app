import React, { useState } from 'react';
import { Form, FormControl, ListGroup } from 'react-bootstrap';
import { API_URL } from '../constants';
import { useNavigate } from 'react-router-dom';
import useLocalStorageState from '../util/userLocalStorage';

function SearchBar({ onSelect }) {
    const [searchTerm, setSearchTerm] = useState('');
    const [searchResults, setSearchResults] = useState([]);
    const [token,] = useLocalStorageState('token');
    const navigate = useNavigate();
    let timeout;

    const handleSearch = async (event) => {
        const query = event.target.value;
        setSearchTerm(query);
        if (query === '') {
            return
        }
        // Clear previous timeout to debounce input events
        clearTimeout(timeout);
        // Set a new timeout to delay sending the HTTP request
        timeout = setTimeout(async () => {
            try {
                const response = await fetch(`${API_URL}/users/search?searchTerm=${query}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`,
                    },
                });
                if (response.ok) {
                    const searchData = await response.json();
                    console.log("Search result:", searchData);
                    setSearchResults(searchData);
                } else if (response.status === 401) {
                    navigate('/login')
                }
            } catch (error) {
                console.error('Error searching users:', error);
            }

        }, 300); // Wait for 300 milliseconds before sending the request to avoid sending request for every typed letter.
    };

    return (
        <div style={{ position: 'relative' }}>
            <Form>
                <FormControl
                    type="text"
                    placeholder="Search users..."
                    value={searchTerm}
                    onChange={handleSearch}
                />
            </Form>
            {searchTerm && searchResults && (
                <div style={{ position: 'absolute', top: '100%', left: 0, zIndex: 100, width: '100%' }}>
                    <ListGroup>
                        {searchResults.map((user) => (
                            <ListGroup.Item key={user.id} onSelect={onSelect}>
                                {user.nickname} - {user.email}
                            </ListGroup.Item>
                        ))}
                    </ListGroup>
                </div>
            )}
        </div>
    );
}

export default SearchBar;
