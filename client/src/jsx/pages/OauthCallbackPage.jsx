import React, { useEffect, useState } from 'react'
import { useHistory } from 'react-router-dom'

import { postJSON } from '../../utils/api'
import ErrorRow from '../common/ErrorRow'
import Spinner from '../common/Spinner'

const OauthCallbackPage = () => {
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    const history = useHistory()

    useEffect(() => {
        setIsSubmitting(true)
        const urlParams = new URLSearchParams(window.location.search)
        const code = urlParams.get('code')
        const state = urlParams.get('state')
        postJSON('/api/splitwise_oauth_callback', {
            'code': code,
            'state': state,
        })
            .then(data => {
                if (data['message'] === 'success') {
                    history.replace('/splitwise_expenses')
                }
            })
            .catch(e => {
                setErrorMessage(e.message)
            })
            .finally(() => {
                setIsSubmitting(false)
            })
    }, [setIsSubmitting, setErrorMessage])

    return (
        <div>
            <Spinner visible={isSubmitting} />
            <ErrorRow message={errorMessage} />
        </div>
    )
}

export default OauthCallbackPage
