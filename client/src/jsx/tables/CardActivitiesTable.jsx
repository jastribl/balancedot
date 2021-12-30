import React from 'react'
import { Link } from 'react-router-dom'

import { formatAsDate, formatAsMoney } from '../../utils/format'
import { dateComparator } from '../../utils/sorting'
import ExtendableTable from './ExtendableTable'

const CardActivitiesTable = ({ initialSortColumn, ...props }) => {
    if (initialSortColumn === undefined && initialSortColumn !== null) {
        initialSortColumn = 'transaction_date'
    }
    return <ExtendableTable
        columns={[
            'uuid',
            'transaction_date',
            'post_date',
            'description',
            'category',
            'type',
            'amount',
            'splitwise_expense_count',
        ]}
        customRenders={{
            'uuid': (data) =>
                <Link to={'/cards/' + data['card_uuid'] + '/activities/' + data['uuid']}>{data['uuid']}</Link>,
            'transaction_date': (data) => formatAsDate(data['transaction_date']),
            'post_date': (data) => formatAsDate(data['post_date']),
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
        }}
        initialSortColumn={initialSortColumn}
        customSortComparators={{
            'transaction_date': dateComparator,
            'post_date': dateComparator,
        }}
        {...props}
    />
}

export default CardActivitiesTable
