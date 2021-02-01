import React, { useEffect, useState } from 'react'

import { getWithHandling, postJSON } from '../../utils/api'
import ErrorRow from '../common/ErrorRow'
import Spinner from '../common/Spinner'
import SplitwiseLoginCheck from '../SplitwiseLoginCheck'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const SplitwiseExpensesPage = () => {
    const [splitwiseExpenses, setSplitwiseExpenses] = useState(null)
    const [pageLoading, setPageLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)
    const [refreshResponse, setRefreshResponse] = useState(null)

    const refreshSplitwiseExpenses = () => getWithHandling(
        '/api/splitwise_expenses',
        setSplitwiseExpenses,
        setErrorMessage,
        setPageLoading,
    )


    const handleRefreshExpenses = () => {
        setPageLoading(true)
        return postJSON('/api/refresh_splitwise')
            .then(data => {
                setRefreshResponse(data)
                refreshSplitwiseExpenses()
            })
            .catch(e => {
                setPageLoading(false)
                if ('redirect_url' in e) {
                    window.open(e.redirect_url)
                    return
                }
                setErrorMessage(e.message)
            })
    }

    useEffect(() => {
        refreshSplitwiseExpenses()
    }, [
        setPageLoading,
        setSplitwiseExpenses,
        setErrorMessage,
    ])

    // todo: make this nicer looking and more functional
    let refreshResponseRender = null
    if (refreshResponse !== null) {
        refreshResponseRender = (
            <div><pre>{JSON.stringify(refreshResponse, null, 4)}</pre></div>
        )
    }

    return (
        <div>
            <Spinner visible={pageLoading} />
            <h1>Splitwise Expenses</h1>
            <ErrorRow message={errorMessage} />
            <SplitwiseLoginCheck>
                <input type='button' onClick={handleRefreshExpenses} value='Refresh Splitwise' style={{ marginBottom: 25 + 'px' }} />
                {refreshResponseRender}
                <SplitwiseExpenseTable data={splitwiseExpenses} />
            </SplitwiseLoginCheck>
        </div>
    )
}

export default SplitwiseExpensesPage
