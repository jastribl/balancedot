import Moment from 'moment'
import React, { useEffect, useState } from 'react'

import { get, postJSON } from '../../utils/api'
import { formatAsMoney } from '../../utils/format'
import ErrorRow from '../common/ErrorRow'
import Spinner from '../common/Spinner'
import Table from '../common/Table'
import SplitwiseLoginCheck from '../SplitwiseLoginCheck'

const SplitwiseExpensesPage = () => {
    const [splitwiseExpenses, setSplitwiseExpenses] = useState(null)
    const [pageLoading, setPageLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)
    const [refreshResponse, setRefreshResponse] = useState(null)

    const refreshSplitwiseExpenses = () => {
        setPageLoading(true)
        get('/api/splitwise_expenses')
            .then(splitwiseExpensesResponse => {
                setSplitwiseExpenses(splitwiseExpensesResponse)
            })
            .catch(e => {
                setErrorMessage(e.message)
            })
            .finally(() => {
                setPageLoading(false)
            })
    }

    const handleRefreshExpenses = () => {
        setPageLoading(true)
        return postJSON('/api/refresh_splitwise', null, 'follow')
            .then((data) => {
                setRefreshResponse(data)
                refreshSplitwiseExpenses()
            })
            .catch(e => {
                if ('redirect_url' in e) {
                    window.open(e.redirect_url)
                    return
                }
                setErrorMessage(e.message)
            })
            .finally(() => {
                setPageLoading(false)
            })
    }

    useEffect(() => {
        refreshSplitwiseExpenses()
    }, [setSplitwiseExpenses])

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
                <div>
                    <Table rowKey='uuid' columns={{
                        'uuid': 'UUID',
                        'splitwise_id': 'splitwise_id',
                        'description': 'Description',
                        'details': 'Details',
                        'amount': 'Amount',
                        'amount_paid': 'Amount Paid',
                        'date': 'Date',
                        'category': 'Category',
                    }} rows={splitwiseExpenses} customRenders={{
                        'details': (data) => data['details'].trim(),
                        'date': (data) =>
                            Moment(data['date']).format('YYYY-MM-DD'),
                        'amount': (data) => formatAsMoney(data['amount'], data['currency_code']),
                        'amount_paid': (data) => formatAsMoney(data['amount_paid'], data['currency_code'])
                    }} />
                </div>
            </SplitwiseLoginCheck>
        </div>
    )
}

export default SplitwiseExpensesPage
