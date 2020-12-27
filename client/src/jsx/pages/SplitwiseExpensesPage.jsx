import React, { useEffect, useState } from 'react'
import Moment from 'moment'

import { formatAsMoney } from '../../utils/format'
import { postJSON, get } from '../../utils/api'

import Spinner from "../common/Spinner"
import SplitwiseLoginCheck from '../SplitwiseLoginCheck'
import Table from "../common/Table"

const SplitwiseExpensesPage = () => {
    const [splitiwseExpenses, setSplitwiseExpenses] = useState(null)
    const [pageLoading, setPageLoading] = useState(false)

    const refreshSplitwiseExpenses = () => {
        setPageLoading(true)
        get('/api/splitwise_expenses')
            .then(splitwiseExpensesResponse => {
                setSplitwiseExpenses(splitwiseExpensesResponse)
            })
            .catch(e => { })
            .finally(() => {
                setPageLoading(false)
            })
    }

    const handleRefreshExpenses = () => {
        setPageLoading(true)
        return postJSON('/api/refresh_splitwise', null, 'follow')
            .then((data) => {
                console.log('got back data', data)
                refreshSplitwiseExpenses()
            })
            .catch(e => {
                if ('redirect_url' in e) {
                    window.open(e.redirect_url)
                    return
                }
            })
            .finally(() => {
                setPageLoading(false)
            })
    }

    useEffect(() => {
        refreshSplitwiseExpenses()
    }, [setSplitwiseExpenses])

    return (
        <div>
            <Spinner visible={pageLoading} />
            <h1>Splitwise Expenses</h1>
            <SplitwiseLoginCheck>
                <input type="button" onClick={handleRefreshExpenses} value="Refresh Splitwise" style={{ marginBottom: 25 + 'px' }} />
                <div>
                    <Table rowKey="uuid" columns={{
                        'uuid': 'UUID',
                        'description': 'Description',
                        'details': 'Details',
                        'amount': 'Amount',
                        'amount_paid': 'Amount Paid',
                        'date': 'Date',
                        'category': 'Category',
                    }} rows={splitiwseExpenses} customRenders={{
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
