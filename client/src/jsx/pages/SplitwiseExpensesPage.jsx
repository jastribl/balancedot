import React, { useState } from 'react'

import { postJSON } from '../../utils/api'
import LoaderComponent from '../common/LoaderComponent'
import SplitwiseLoginCheck from '../SplitwiseLoginCheck'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const SplitwiseExpensesPage = () => {
    const [splitwiseExpenses, setSplitwiseExpenses] = useState(null)
    const [pageLoading, setPageLoading] = useState(false)
    const [refreshResponse, setRefreshResponse] = useState(null)

    const handleRefreshExpenses = () => {
        setPageLoading(true)
        return postJSON('/api/refresh_splitwise')
            .then(data => setRefreshResponse(data))
            .catch(e => {
                if ('redirect_url' in e) {
                    window.open(e.redirect_url)
                    return
                }
                setErrorMessage(e.message)
            })
            .finally(() => setPageLoading(false))
    }

    // todo: make this nicer looking and more functional
    let refreshResponseRender = null
    if (refreshResponse !== null) {
        refreshResponseRender = (
            <div><pre>{JSON.stringify(refreshResponse, null, 4)}</pre></div>
        )
    }

    return (
        <div>
            <h1>Splitwise Expenses</h1>
            <LoaderComponent
                path={'/api/splitwise_expenses'}
                parentLoading={pageLoading}
                setData={setSplitwiseExpenses}
            />
            <SplitwiseLoginCheck>
                <input type='button' onClick={handleRefreshExpenses} value='Refresh Splitwise' style={{ marginBottom: 25 + 'px' }} />
                {refreshResponseRender}
                <SplitwiseExpenseTable data={splitwiseExpenses} />
            </SplitwiseLoginCheck>
        </div>
    )
}

export default SplitwiseExpensesPage
