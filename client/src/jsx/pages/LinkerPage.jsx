import Moment from 'moment'
import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { get } from '../../utils/api'
import { formatAsMoney } from '../../utils/format'
import ErrorRow from '../common/ErrorRow'
import Spinner from '../common/Spinner'
import Table from '../common/Table'

const LinkerPage = () => {
    const [splitwiseExpenses, setSplitwiseExpenses] = useState(null)
    const [pageLoading, setPageLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    useEffect(() => {
        setPageLoading(false)
        get('/api/splitwise_expenses/unlinked')
            .then(splitwiseExpensesResponse => {
                setSplitwiseExpenses(splitwiseExpensesResponse)
            })
            .catch(e => {
                setErrorMessage(e.message)
            })
            .finally(() => {
                setPageLoading(false)
            })
    }, [setPageLoading, setSplitwiseExpenses, setErrorMessage])

    return (
        <div>
            <Spinner visible={pageLoading} />
            <h1>Splitwise Expense Linking</h1>
            <ErrorRow message={errorMessage} />
            <div>
                <Table
                    rowKey='uuid'
                    rows={splitwiseExpenses}
                    columns={['uuid', 'splitwise_id', 'description', 'details', 'amount', 'amount_paid', 'date', 'category', 'link']}
                    customRenders={{
                        'details': (data) => data['details'].trim(),
                        'date': (data) =>
                            Moment(data['date']).format('YYYY-MM-DD'),
                        'amount': (data) => formatAsMoney(data['amount'], data['currency_code']),
                        'amount_paid': (data) => formatAsMoney(data['amount_paid'], data['currency_code']),
                        // todo: try moving linking links to the main splitwise page
                        // todo: style link the rest of the buttons (along with other links)
                        'link': (data) => <Link to={'/linker/' + data['uuid']}>Link</Link>,
                    }}
                />
            </div>
        </div>
    )
}

export default LinkerPage
