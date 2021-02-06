import React, { useEffect, useState } from 'react'

import { getWithHandling } from '../../utils/api'
import ErrorRow from './ErrorRow'
import Spinner from './Spinner'

const LoaderComponent = ({
    path,
    parentLoading,
    parentErrorMessage,
    setData,
}) => {
    parentLoading ??= false
    parentErrorMessage ??= null

    const [loading, setLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    useEffect(() => {
        if (!parentLoading) {
            getWithHandling(
                path,
                setData,
                setErrorMessage,
                setLoading
            )
        }
    }, [
        path,
        parentLoading,
        parentErrorMessage,
    ])

    return <div>
        <Spinner visible={loading || parentLoading} />
        <ErrorRow message={errorMessage ?? parentErrorMessage} />
    </div>


}

export default LoaderComponent
