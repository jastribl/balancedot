import React, { useEffect, useState } from 'react'

import { postJSON, get } from '../utils/api'

const SplitwiseLoginCheck = (props) => {
    const [hasSplitwiseSetup, setHasSplitiwseSetup] = useState(false)
    const [authURL, setAuthURL] = useState(null)

    const onAuthClick = () => {
        if (authURL !== null) {
            window.location = authURL
        }
    }

    useEffect(() => {
        get('/api/splitwise_login_check')
            .then(splitwiseLoginCheckResponse => {
                setHasSplitiwseSetup(true)
            })
            .catch(e => {
                setHasSplitiwseSetup(false)
                if (e.status === 401 && e.message === "Authentication Response") {
                    setAuthURL(e['redirect_url'])
                }
            })
    }, [setHasSplitiwseSetup, setAuthURL])

    if (hasSplitwiseSetup) {
        return props.children
    }

    if (authURL !== null) {
        return <input type="button" onClick={onAuthClick} value="Link Splitwise" style={{ marginBottom: 25 + 'px' }} />
    }

    return <div />
}

export default SplitwiseLoginCheck