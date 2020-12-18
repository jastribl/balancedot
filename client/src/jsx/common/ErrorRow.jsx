import React from 'react'

const ErrorRow = ({ message }) => {
    if (message === null) {
        return null
    }

    return (
        <div className="row isa_error">
            <span>
                <div style={{ textAlign: 'center', color: 'red' }}>{message}</div>
            </span>
        </div>
    )
}

export default ErrorRow
