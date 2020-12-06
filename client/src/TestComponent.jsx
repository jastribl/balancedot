import React, { useEffect, useState } from 'react';

import List from "./List"

const TestComponent = () => {
    const [state, setState] = useState({
        items: null,
    });

    useEffect(() => {
        const apiUrl = `/api/json`;
        fetch(apiUrl)
            .then((res) => res.json())
            .then((items) => {
                setState({ items: items });
            });
    }, [setState]);
    return (
        <div className='App'>
            <div className='container'>
                <h1>My Items</h1>
            </div>
            <div className='repo-container'>
                <List items={state.items} />
            </div>
        </div>
    );
}


export default TestComponent;