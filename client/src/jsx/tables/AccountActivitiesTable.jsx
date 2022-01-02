import React from 'react'
import { Link } from 'react-router-dom'

import { formatAsDate, formatAsMoney } from '../../utils/format'
import { dateComparator } from '../../utils/sorting'
import ExtendableTable from './ExtendableTable'

const AccountActivitiesTable = ({ initialSortColumn, ...props }) => {
    if (initialSortColumn === undefined && initialSortColumn !== null) {
        initialSortColumn = 'posting_date'
    }
    return <ExtendableTable
        columns={[
            'uuid',
            'details',
            'posting_date',
            'description',
            'amount',
            'type',
            'splitwise_expense_count',
            'card_activities_count',
        ]}
        customRenders={{
            'uuid': (data) =>
                <Link to={'/accounts/' + data['account_uuid'] + '/activities/' + data['uuid']}>{data['uuid']}</Link>,
            'posting_date': (data) => formatAsDate(data['posting_date']),
            'amount': (data) => formatAsMoney(data['amount']),
            'splitwise_expense_count': (data) => {
                const splitwiseExpenses = data['splitwise_expenses']
                const num = splitwiseExpenses?.length
                if (num > 0) {
                    const sum = splitwiseExpenses
                        .map(d => d.amount_paid)
                        .reduce((a, b) => a + b, 0)
                        .toFixed(2)
                    return <div style={{
                        color: (Math.abs(Math.abs(sum) - Math.abs(data['amount'])) < 0.03 ? 'green' : 'red')
                    }}>{`${num} (${sum})`}</div>
                } else if (num === undefined) {
                    return 'Not loaded...'
                }
                return ''
            },
            'card_activities_count': (data) => {
                const cardActivities = data['card_activities']
                const num = cardActivities?.length
                if (num > 0) {
                    const sum = cardActivities
                        .map(d => d.amount)
                        .reduce((a, b) => a + b, 0)
                        .toFixed(2)
                    return <div style={{
                        color: (Math.abs(Math.abs(sum) - Math.abs(data['amount'])) === 0.00 ? 'green' : 'red')
                    }}>{`${num} (${sum})`}</div>
                } else if (num === undefined) {
                    return 'Not loaded...'
                }
                return ''
            },
        }}
        initialSortColumn={initialSortColumn}
        initialSortInverse
        customSortComparators={{
            'posting_date': dateComparator
        }}
        {...props}
    />
}

export default AccountActivitiesTable
