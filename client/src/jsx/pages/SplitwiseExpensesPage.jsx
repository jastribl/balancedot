import React, { useState } from 'react'

import { postJSON } from '../../utils/api'
import LoaderComponent from '../common/LoaderComponent'
import Toggle from '../common/Toggle'
import SplitwiseLoginCheck from '../SplitwiseLoginCheck'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const SplitwiseExpensesPage = () => {
    const [unlinkedOnly, setUnlinkedOnly] = useState(false)
    const [splitwiseExpenses, setSplitwiseExpenses] = useState(null)
    const [refreshingSplitwise, setRefreshingSplitwise] = useState(false)
    const [rawRefreshResponse, setRawRefreshResponse] = useState(null)

    const handleUnlinkedToggle = () => setUnlinkedOnly(!unlinkedOnly)

    const handleRefreshExpenses = () => {
        setRefreshingSplitwise(true)
        return postJSON('/api/refresh_splitwise')
            .then(data => setRawRefreshResponse(data))
            .catch(e => {
                if ('redirect_url' in e) {
                    window.open(e.redirect_url)
                    return
                }
                setErrorMessage(e.message)
            })
            .finally(() => setRefreshingSplitwise(false))
    }

    // todo: make this nicer looking and more functional
    let refreshResponseRender = null
    if (rawRefreshResponse !== null) {
        refreshResponseRender = (
            <div><pre>{JSON.stringify(rawRefreshResponse, null, 4)}</pre></div>
        )
    }

    return (
        <div>
            <h1>Splitwise Expenses</h1>
            <LoaderComponent
                path={unlinkedOnly ? '/api/splitwise_expenses/unlinked' : '/api/splitwise_expenses'}
                parentLoading={refreshingSplitwise}
                setData={setSplitwiseExpenses}
            />
            <SplitwiseLoginCheck>
                <input
                    type='button'
                    onClick={handleRefreshExpenses}
                    value='Refresh Splitwise'
                    style={{ marginBottom: 25 + 'px' }}
                />
                {refreshResponseRender}
                <div>Unlinked Only: <Toggle onToggle={handleUnlinkedToggle} /></div>
                <SplitwiseExpenseTable
                    data={splitwiseExpenses}
                />
            </SplitwiseLoginCheck>
        </div>
    )
}

export default SplitwiseExpensesPage
