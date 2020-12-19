import React from 'react'

const Spinner = ({ visible }) => {
    return (
        <div style={{
            position: 'absolute',
            display: visible ? 'block' : 'none',
            background: 'whitesmoke',
            opacity: '50%',
            height: '100%',
            width: '100%',
            zIndex: '2'
        }}>
            <div style={{
                position: 'absolute',
                height: '100%',
                left: '50%',
                transform: 'translate(-50%, 0%)',
                display: 'flex',
                alignItems: 'center'
            }}>
                <div style={{
                    height: '64px',
                    width: '64px',
                    animation: 'rotate 2s linear infinite',
                    border: '5px solid firebrick',
                    borderRightColor: 'transparent',
                    borderRadius: '50%',
                }} />
            </div>
        </div >
    )
}

export default Spinner
