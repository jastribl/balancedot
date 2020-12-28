import React, { useEffect, useState } from 'react'

import { get } from '../utils/api'
import ErrorBar from './common/ErrorRow'

const SplitwiseLoginCheck = (props) => {
    const [hasSplitwiseSetup, setHasSplitwiseSetup] = useState(false)
    const [authURL, setAuthURL] = useState(null)
    const [errorMessage, setErrorMessage] = useState(null)

    const onAuthClick = () => {
        if (authURL !== null) {
            window.location = authURL
        }
    }

    useEffect(() => {
        get('/api/splitwise_login_check')
            .then(splitwiseLoginCheckResponse => {
                if (splitwiseLoginCheckResponse['message'] === 'success') {
                    setHasSplitwiseSetup(true)
                } else {
                    setErrorMessage('Unknown state white checking splitwise authentication')
                }
            })
            .catch(e => {
                setHasSplitwiseSetup(false)
                if (e.status === 401 && e.message === 'Authentication Required') {
                    setAuthURL(e['redirect_url'])
                } else {
                    setErrorMessage(e.message)
                }
            })
    }, [setHasSplitwiseSetup, setAuthURL])

    if (hasSplitwiseSetup) {
        return props.children
    }

    if (errorMessage !== null) {
        return <ErrorBar message={errorMessage} />
    }

    if (authURL !== null) {
        return <input type='button' onClick={onAuthClick} value='Link Splitwise' style={{ marginBottom: 25 + 'px' }} />
    }

    return <div />
}

export default SplitwiseLoginCheck