import React, { useState } from 'react'
import { Link } from 'react-router-dom'

import { postJSON } from '../../utils/api'
import LoaderComponent from '../common/LoaderComponent'
import SplitwiseLoginCheck from '../SplitwiseLoginCheck'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const SplitwiseExpensesPage = ({ match }) => {
    const unlinkedOnly = match.path.endsWith('/unlinked')

    const [splitwiseExpenses, setSplitwiseExpenses] = useState(null)
    const [refreshingSplitwise, setRefreshingSplitwise] = useState(false)
    const [rawRefreshResponse, setRawRefreshResponse] = useState(null)

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

                <div>Showing {unlinkedOnly ? 'Unlinked Expenses' : 'All'}</div>
                <div>
                    <Link to={`/splitwise_expenses` + (unlinkedOnly ? '/' : '/unlinked')}>
                        <input
                            type='button'
                            value={unlinkedOnly ? 'View All' : 'View Unlinked'}
                            style={{ marginTop: 25 + 'px', marginBottom: 25 + 'px' }}
                        />
                    </Link>
                </div>

                <SplitwiseExpenseTable
                    data={splitwiseExpenses}
                />
            </SplitwiseLoginCheck>
        </div>
    )
}

export default SplitwiseExpensesPage
