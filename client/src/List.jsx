import React from 'react';

const List = (props) => {
    const { items } = props;
    if (!items) {
        return <div/>
    }
    return (
        <ul>
            {items.map((item) => {
                return (
                    <li key={item.id}>
                        <span>{item.a} </span>
                        <span>{item.b}</span>
                    </li>
                );
            })}
        </ul>
    );
};

export default List;