import React from 'react'
import { Link } from 'react-router-dom'

import { formatAsDate, formatAsMoney } from '../../utils/format'
import { dateComparator } from '../../utils/sorting'
import ExtendableTable from './ExtendableTable'

const AccountActivitiesTable = (props) => {
    return <ExtendableTable
        columns={[
            'uuid',
            'details',
            'posting_date',
            'description',
            'amount',
            'type',
            'splitwise_expense_count',
        ]}
        customRenders={{
            'uuid': (data) =>
                <Link to={'/accounts/' + data['account_uuid'] + '/activities/' + data['uuid']}>{data['uuid']}</Link>,
            'posting_date': (data) => formatAsDate(data['posting_date']),
            'amount': (data) => formatAsMoney(data['amount']),
            'splitwise_expense_count': (data) => {
                const splitwiseExpenses = data['splitwise_expenses']
                const num = splitwiseExpenses.length
                if (num > 0) {
                    const sum = splitwiseExpenses
                        .map(d => d.amount_paid)
                        .reduce((a, b) => a + b, 0)
                    return `${num} (${sum})`
                }
                return ''
            },
        }}
        initialSortColumn='posting_date'
        customSortComparators={{
            'posting_date': dateComparator
        }}
        {...props}
    />
}

export default AccountActivitiesTable
