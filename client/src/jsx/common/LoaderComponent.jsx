import React, { useEffect, useState } from 'react'

import { getWithHandling } from '../../utils/api'
import ErrorRow from './ErrorRow'
import Spinner from './Spinner'

const LoaderComponent = ({ path, parentLoading, setData }) => {
    parentLoading ??= false

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
    ])

    return <div>
        <Spinner visible={loading || parentLoading} />
        <ErrorRow message={errorMessage} />
    </div>


}

export default LoaderComponent
