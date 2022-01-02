import React, { useEffect, useState } from 'react'

import { getWithParamsWithHandling } from '../../utils/api'
import ErrorRow from './ErrorRow'
import Spinner from './Spinner'

const LoaderComponent = ({
    path,
    queryParams,
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
            // todo: figure out why this is hitting twice
            getWithParamsWithHandling(
                path,
                queryParams ?? {},
                setData,
                setErrorMessage,
                setLoading
            )
        }
    }, [
        path,
        queryParams,
        parentLoading,
        parentErrorMessage,
    ])

    return <div>
        <Spinner visible={loading || parentLoading} />
        <ErrorRow message={errorMessage ?? parentErrorMessage} />
    </div>
}

export default LoaderComponent
