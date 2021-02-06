import React, { useState } from 'react'

const Toggle = ({ initialValue, onToggle }) => {
    const [toggleOn, setToggleOn] = useState(initialValue ?? false)

    const handleToggle = () => {
        setToggleOn(!toggleOn)
        if (onToggle) {
            onToggle()
        }
    }

    return <span>
        <label className="switch">
            <input type="checkbox" checked={toggleOn} onChange={handleToggle} />
            <span className="slider round"></span>
        </label>
    </span>
}

export default Toggle
