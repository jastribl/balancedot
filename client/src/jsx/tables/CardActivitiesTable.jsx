import React from 'react'
import { Link } from 'react-router-dom'

import { formatAsDate, formatAsMoney } from '../../utils/format'
import { dateComparator } from '../../utils/sorting'
import ExtendableTable from './ExtendableTable'

const CardActivitiesTable = (props) => {
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
        initialSortColumn='transaction_date'
        customSortComparators={{
            'transaction_date': dateComparator,
            'post_date': dateComparator,
        }}
        {...props}
    />
}

export default CardActivitiesTable
